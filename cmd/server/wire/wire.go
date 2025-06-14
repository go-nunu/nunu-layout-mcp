//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/go-nunu/nunu-layout-mcp/internal/handler"
	"github.com/go-nunu/nunu-layout-mcp/internal/repository"
	"github.com/go-nunu/nunu-layout-mcp/internal/server"
	"github.com/go-nunu/nunu-layout-mcp/internal/service"
	"github.com/go-nunu/nunu-layout-mcp/pkg/app"
	"github.com/go-nunu/nunu-layout-mcp/pkg/log"
	"github.com/go-nunu/nunu-layout-mcp/pkg/server/mcp"
	"github.com/go-nunu/nunu-layout-mcp/pkg/sid"
	"github.com/google/wire"
	"github.com/spf13/viper"
)

var repositorySet = wire.NewSet(
	//repository.NewDB,
	//repository.NewTransaction,
	//repository.NewRedis,
	repository.NewRepository,
	repository.NewExampleRepository,
)

var serviceSet = wire.NewSet(
	service.NewService,
	service.NewExampleService,
)

var handlerSet = wire.NewSet(
	handler.NewHandler,
	handler.NewExampleHandler,
)

var serverSet = wire.NewSet(
	server.NewMCPServer,
)

// build App
func newApp(
	mcpServer *mcp.Server,
) *app.App {
	return app.NewApp(
		app.WithServer(
			mcpServer,
		),
		app.WithName("demo-server"),
	)
}

func NewWire(*viper.Viper, *log.Logger) (*app.App, func(), error) {
	panic(wire.Build(
		repositorySet,
		serviceSet,
		serverSet,
		handlerSet,
		sid.NewSid,
		newApp,
	))
}
