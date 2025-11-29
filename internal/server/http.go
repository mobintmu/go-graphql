package server

import (
	"context"
	"fmt"
	"net/http"

	"go-graphql/internal/config"
	"go-graphql/internal/graph/generated"
	"go-graphql/internal/graph/resolvers"

	// gqlgen generated package
	product "go-graphql/internal/product/service"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type HTTPServer struct {
	server *http.Server
	config *config.Config
	log    *zap.Logger
}

func NewHTTPServer(cfg *config.Config, log *zap.Logger) *HTTPServer {
	return &HTTPServer{
		config: cfg,
		log:    log,
	}
}

// NewGraphQLResolver wires your services into the gqlgen resolvers.
func NewGraphQLResolver(productSvc *product.Product) *resolvers.Resolver {
	return &resolvers.Resolver{
		Product: productSvc,
	}
}

func RegisterGraphQLRoutes(
	hs *HTTPServer,
	resolver *resolvers.Resolver, // Assuming your resolver package is named 'resolvers'
) {
	mux := http.NewServeMux()

	// Create the executable schema using the injected resolver
	schema := generated.NewExecutableSchema(generated.Config{Resolvers: resolver})

	// GraphQL handler
	srv := handler.NewDefaultServer(schema)
	mux.Handle("/query", srv)

	// GraphQL Playground
	mux.Handle("/playground", playground.Handler("GraphQL Playground", "/query"))

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"status":"ok"}`)
	})

	hs.server = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", hs.config.HTTPAddress, hs.config.HTTPPort),
		Handler: mux,
	}
}

func StartHTTPServer(lc fx.Lifecycle, hs *HTTPServer) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				hs.log.Info("Server running",
					zap.String("url", fmt.Sprintf("http://%s:%d/playground", hs.config.HTTPAddress, hs.config.HTTPPort)),
				)
				if err := hs.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					hs.log.Error("Server error", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			hs.log.Info("Stopping server...")
			return hs.server.Shutdown(ctx)
		},
	})
}
