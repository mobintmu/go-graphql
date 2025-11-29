package app

import (
	"go-graphql/internal/config" // gqlgen generated package
	// your resolvers
	"go-graphql/internal/health"
	"go-graphql/internal/pkg/logger"
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
			// health check
			health.New,
			// server
			server.NewHTTPServer,
			// db
			migrate.NewRunner,
			sqlc.New,
			// cache
			cache.NewClient,
			cache.NewCacheStore,
			// services
			productService.New,
			// GraphQL
			server.NewGraphQLResolver,
		),
		fx.Invoke(
			server.RegisterGraphQLRoutes,
			server.StartHTTPServer,
			// migration
			migrate.RunMigrations,
			// life cycle
			logger.RegisterLoggerLifecycle,
		),
	)
}
