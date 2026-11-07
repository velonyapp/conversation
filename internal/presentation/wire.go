package presentation

import (
	"github.com/velonyapp/conversation/internal/presentation/middleware"
	"github.com/velonyapp/conversation/internal/presentation/service"
	"github.com/velonyapp/conversation/internal/presentation/transport"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	service.NewAPI,
	transport.NewGRPCServer,
	transport.NewHTTPServer,
	middleware.NewTracing,
	middleware.NewMetrics,
	middleware.NewErrorMapper,
	middleware.NewAuthentication,
	middleware.NewValidation,
)
