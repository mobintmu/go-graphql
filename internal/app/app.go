package app

import (
	"go-graphql/internal/config"
	"go-graphql/internal/graph"
	"go-graphql/internal/graph/resolvers"
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
			NewGraphQLResolver,
			NewGraphQLSchema,
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

func NewGraphQLResolver(productService *productService.Product) *graph.Resolver {
	return &graph.Resolver{
		Product: productService,
	}
}

func NewGraphQLSchema(resolver *graph.Resolver) *graph.ExecutableSchema {
	return graph.NewExecutableSchema(graph.Config{
		Resolvers: &resolvers.Resolver{
			Resolver: resolver,
		},
	})
}
