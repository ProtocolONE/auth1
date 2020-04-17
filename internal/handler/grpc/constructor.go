package grpc

import (
	"github.com/ProtocolONE/auth1.protocol.one/internal/domain/service"
	"go.uber.org/fx"
)

type Params struct {
	fx.In

	ProfileService service.ProfileService
	UserService    service.UserService
}

func New(params Params) *Handler {
	return &Handler{
		profile: params.ProfileService,
		user:    params.UserService,
	}
}
