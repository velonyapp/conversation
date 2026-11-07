//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"log/slog"

	"github.com/velonyapp/conversation/internal/application"
	"github.com/velonyapp/conversation/internal/conf"
	"github.com/velonyapp/conversation/internal/domain"
	"github.com/velonyapp/conversation/internal/infrastructure"
	"github.com/velonyapp/conversation/internal/presentation"

	"github.com/go-kratos/kratos/v3"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(
	*conf.Auth,
	*conf.Transport,
	*conf.Data,
	*slog.Logger,
) (*kratos.App, func(), error) {
	panic(wire.Build(
		domain.ProviderSet,
		application.ProviderSet,
		infrastructure.ProviderSet,
		presentation.ProviderSet,
		newApp,
	))
}
