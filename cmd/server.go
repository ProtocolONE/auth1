package cmd

import (
	"github.com/ProtocolONE/auth1.protocol.one/pkg/api"
	"github.com/ProtocolONE/auth1.protocol.one/pkg/config"
	"github.com/ProtocolONE/auth1.protocol.one/pkg/database"
	"github.com/ProtocolONE/mfa-service/pkg"
	"github.com/ProtocolONE/mfa-service/pkg/proto"
	"github.com/boj/redistore"
	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/go-redis/redis"
	"github.com/micro/go-micro"
	"github.com/micro/go-plugins/client/selector/static"
	"github.com/ory/hydra/sdk/go/hydra/client"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"log"
	"net/http"
	"net/url"
	"os"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run AuthOne api server with given configuration",
	Run:   runServer,
}

func init() {
	command.AddCommand(serverCmd)
}

func runServer(cmd *cobra.Command, args []string) {
	db := createDatabase(&cfg.Database)
	defer db.Close()

	go func() {
		log.Println(http.ListenAndServe(":6060", nil))
	}()

	store, err := redistore.NewRediStore(
		cfg.Session.Size,
		cfg.Session.Network,
		cfg.Session.Address,
		cfg.Session.Password,
		[]byte(cfg.Session.Secret),
	)
	if err != nil {
		zap.L().Fatal("Unable to start redis session store", zap.Error(err))
	}
	defer store.Close()

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
	})
	defer redisClient.Close()

	var options []micro.Option

	if os.Getenv("MICRO_SELECTOR") == "static" {
		zap.L().Info("Use micro selector `static`")
		options = append(options, micro.Selector(static.NewSelector()))
	}

	zap.L().Info("Initialize micro service")

	service := micro.NewService(options...)
	service.Init()

	ms := proto.NewMfaService(mfa.ServiceName, service.Client())

	u, err := url.Parse(cfg.Hydra.AdminURL)
	if err != nil {
		zap.L().Fatal("Invalid of the Hydra admin url", zap.Error(err))
	}

	transport := httptransport.New(u.Host, "", []string{u.Scheme})
	transport.DefaultAuthentication = runtime.ClientAuthInfoWriterFunc(func(req runtime.ClientRequest, _ strfmt.Registry) error {
		req.SetHeaderParam("X-Forwarded-Proto", "https")
		return nil
	})
	hydraSDK := client.New(transport, nil)
	if err != nil {
		zap.L().Fatal("Hydra SDK creation failed", zap.Error(err))
	}

	serverConfig := api.ServerConfig{
		ApiConfig:     &cfg.Server,
		HydraConfig:   &cfg.Hydra,
		SessionConfig: &cfg.Session,
		MfaService:    ms,
		MgoSession:    db,
		SessionStore:  store,
		RedisClient:   redisClient,
		HydraAdminApi: hydraSDK.Admin,
		Mailer:        &cfg.Mailer,
	}

	server, err := api.NewServer(&serverConfig)
	if err != nil {
		zap.L().Fatal("Failed to create server", zap.Error(err))
	}

	zap.L().Info("Starting up server")
	if err = server.Start(); err != nil {
		zap.L().Fatal("Error running server", zap.Error(err))
	}
}

func createDatabase(cfg *config.Database) database.MgoSession {
	db, err := database.NewConnection(cfg)
	if err != nil {
		zap.L().Fatal("Name connection failed with error", zap.Error(err))
	}

	return db
}
