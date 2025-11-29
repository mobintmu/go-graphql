package app

import (
	"go-graphql/internal/config"
	"go-graphql/internal/health"
	"go-graphql/internal/pkg/logger"
	productController "go-graphql/internal/product/controller"
	productService "go-graphql/internal/product/service"
	"go-graphql/internal/server"
	"go-graphql/internal/storage/cache"
	"go-graphql/internal/storage/sql"
	"go-graphql/internal/storage/sql/migrate"
	"go-graphql/internal/storage/sql/sqlc"

	"go.uber.org/fx"
)

func NewApp() *fx.App {
	return fx.New(
		fx.Provide(
			logger.NewLogger,
			config.NewConfig,
			sql.InitialDB,
			//server
			health.New,
			server.NewGinEngine,
			server.CreateHTTPServer,
			//db
			migrate.NewRunner, // migration runner
			sqlc.New,
			//cache
			cache.NewClient,
			cache.NewCacheStore,
			//controller
			productController.NewAdmin,
			productController.NewClient,
			//service
			productService.New,
		),
		fx.Invoke(
			server.RegisterRoutes,
			server.StartHTTPServer,
			//migration
			migrate.RunMigrations,
			//life cycle
			logger.RegisterLoggerLifecycle,
		),
	)
}
