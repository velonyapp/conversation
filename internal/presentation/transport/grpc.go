package transport

import (
	v1 "github.com/velonyapp/conversation/gen/api/v1"
	"github.com/velonyapp/conversation/internal/conf"
	"github.com/velonyapp/conversation/internal/presentation/middleware"
	"github.com/velonyapp/conversation/internal/presentation/service"

	"github.com/go-kratos/kratos/v3/middleware/recovery"
	"github.com/go-kratos/kratos/v3/transport/grpc"
)

func NewGRPCServer(
	c *conf.Transport,
	service *service.API,
	tracing middleware.Tracing,
	metrics middleware.Metrics,
	errorMapper middleware.ErrorMapper,
	authentication middleware.Authentication,
	validation middleware.Validation,
) *grpc.Server {
	opts := []grpc.ServerOption{
		grpc.Address(c.Grpc.Address),
		grpc.Middleware(
			recovery.Recovery(),
			tracing.Middleware(),
			metrics.Middleware(),
			errorMapper.Middleware(),
			authentication.Middleware(),
			validation.Middleware(),
		),
	}

	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}

	srv := grpc.NewServer(opts...)

	v1.RegisterConversationServiceServer(srv, service)

	return srv
}
