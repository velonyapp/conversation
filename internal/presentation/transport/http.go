package transport

import (
	v1 "github.com/velonyapp/conversation/gen/api/v1"
	"github.com/velonyapp/conversation/internal/conf"
	"github.com/velonyapp/conversation/internal/presentation/middleware"
	"github.com/velonyapp/conversation/internal/presentation/service"

	"github.com/go-kratos/kratos/v3/middleware/recovery"
	"github.com/go-kratos/kratos/v3/transport/http"
	"github.com/rs/cors"
)

func NewHTTPServer(
	c *conf.Transport,
	service *service.API,
	tracing middleware.Tracing,
	metrics middleware.Metrics,
	errorMapper middleware.ErrorMapper,
	authentication middleware.Authentication,
	validation middleware.Validation,
) *http.Server {
	opts := []http.ServerOption{
		http.Address(c.Http.Address),
		http.Filter(cors.New(cors.Options{
			AllowedOrigins: []string{
				"https://api-docs.velony.app",
			},
			AllowedMethods: []string{
				"GET",
				"POST",
				"PUT",
				"PATCH",
				"DELETE",
			},
			AllowedHeaders: []string{
				"Content-Type",
				"Authorization",
			},
		}).Handler),
		http.Middleware(
			recovery.Recovery(),
			tracing.Middleware(),
			metrics.Middleware(),
			errorMapper.Middleware(),
			authentication.Middleware(),
			validation.Middleware(),
		),
	}

	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}

	srv := http.NewServer(opts...)

	v1.RegisterConversationServiceHTTPServer(srv, service)

	return srv
}
