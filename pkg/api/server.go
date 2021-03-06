package api

import (
	"context"
	"github.com/ProtocolONE/auth1.protocol.one/pkg/config"
	"github.com/ProtocolONE/auth1.protocol.one/pkg/database"
	"github.com/ProtocolONE/auth1.protocol.one/pkg/helper"
	"github.com/ProtocolONE/auth1.protocol.one/pkg/models"
	"github.com/ProtocolONE/auth1.protocol.one/pkg/service"
	"github.com/ProtocolONE/mfa-service/pkg/proto"
	"github.com/boj/redistore"
	"github.com/go-redis/redis"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/ory/hydra/sdk/go/hydra/client/admin"
	"go.uber.org/zap"
	"gopkg.in/go-playground/validator.v9"
	"html/template"
	"io"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// ServerConfig contains common configuration parameters for start application server
type ServerConfig struct {
	// ApiConfig is common http setting for the application like a port, timeouts & etc.
	ApiConfig *config.Server

	// HydraConfig contains settings for the public and admin url of the Hydra application.
	HydraConfig *config.Hydra

	// HydraAdminApi is client of the Hydra for administration requests.
	HydraAdminApi *admin.Client

	// SessionConfig contains settings for the session.
	SessionConfig *config.Session

	// MfaService describes the interface for working with MFA micro-service.
	MfaService proto.MfaService

	// MgoSession describes the interface for working with Mongo session.
	MgoSession database.MgoSession

	// SessionStore is client for session storage.
	SessionStore *redistore.RediStore

	// RedisClient is Redis client.
	RedisClient *redis.Client

	// Mailer contains settings for the postman service
	Mailer *config.Mailer
}

// Server is the instance of the application
type Server struct {
	// Echo is instance of the Echo framework
	Echo *echo.Echo

	// ApiConfig is common http setting for the application like a port, timeouts & etc.
	ServerConfig *config.Server

	// RedisClient is Redis client.
	RedisHandler *redis.Client

	// HydraConfig contains settings for the public and admin url of the Hydra application.
	HydraConfig *config.Hydra

	// SessionConfig contains settings for the session.
	SessionConfig *config.Session

	// Registry is the registry service
	Registry service.InternalRegistry
}

// Template is used to display HTML pages.
type Template struct {
	templates *template.Template
}

// NewServer creates new instance of the application.
func NewServer(c *ServerConfig) (*Server, error) {
	registryConfig := &service.RegistryConfig{
		MgoSession:    c.MgoSession,
		HydraAdminApi: c.HydraAdminApi,
		MfaService:    c.MfaService,
		RedisClient:   c.RedisClient,
		Mailer:        service.NewMailer(c.Mailer),
	}
	server := &Server{
		Echo:          echo.New(),
		RedisHandler:  c.RedisClient,
		ServerConfig:  c.ApiConfig,
		SessionConfig: c.SessionConfig,
		HydraConfig:   c.HydraConfig,
		Registry:      service.NewRegistryBase(registryConfig),
	}

	t := &Template{
		templates: template.Must(template.ParseGlob("public/templates/*.html")),
	}
	server.Echo.Renderer = t
	server.Echo.HTTPErrorHandler = helper.ErrorHandler
	server.Echo.Use(ZapLogger(zap.L()))
	server.Echo.Use(middleware.Recover())
	// TODO: Validate origins for each application by settings
	server.Echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowHeaders:     []string{echo.HeaderAuthorization, echo.HeaderContentType, echo.HeaderOrigin, echo.HeaderAccept},
		AllowOrigins:     c.ApiConfig.AllowOrigins,
		AllowCredentials: c.ApiConfig.AllowCredentials,
		AllowMethods:     []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete, http.MethodOptions},
	}))
	server.Echo.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup:    "header:X-XSRF-TOKEN",
		CookieName:     "_csrf",
		Skipper:        csrfSkipper,
		CookieSameSite: http.SameSiteNoneMode,
	}))
	server.Echo.Use(session.Middleware(c.SessionStore))
	server.Echo.Use(middleware.RequestID())
	server.Echo.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			db := c.MgoSession.Copy()
			defer db.Close()

			ctx.Set("database", db)

			logger := zap.L().With(
				zap.String(
					echo.HeaderXRequestID,
					ctx.Response().Header().Get(echo.HeaderXRequestID),
				),
			)
			ctx.Set("logger", logger)

			return next(ctx)
		}
	})

	registerCustomValidator(server.Echo)

	if err := server.setupRoutes(); err != nil {
		zap.L().Fatal("Setup routes failed", zap.Error(err))
	}

	return server, nil
}

func registerCustomValidator(e *echo.Echo) {
	v := validator.New()
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})
	e.Validator = &models.CustomValidator{
		Validator: v,
	}
}

func (s *Server) Start() error {
	go func() {
		err := s.Echo.Start(":" + strconv.Itoa(s.ServerConfig.Port))
		if err != nil {
			zap.L().Fatal("Failed to start server", zap.Error(err))
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	select {
	// wait on kill signal
	case <-shutdown:
		zap.L().Fatal("Server is shutting down")
	}

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return s.Echo.Shutdown(ctx)
}

func (s *Server) setupRoutes() error {
	routes := []func(c *Server) error{
		InitLogin,
		InitPasswordLess,
		InitChangePassword,
		InitMFA,
		InitManage,
		InitOauth2,
		InitHealth,
		InitSafariHack,
	}

	for _, r := range routes {
		if err := r(s); err != nil {
			return err
		}
	}

	return nil
}

func (t *Template) Render(w io.Writer, name string, data interface{}, ctx echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func csrfSkipper(ctx echo.Context) bool {
	return true //ctx.Path() != "/oauth2/login" && ctx.Path() != "/oauth2/signup"
}
