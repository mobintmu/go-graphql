package server

import (
	"context"
	"fmt"
	"net/http"

	"go-graphql/internal/config"
	"go-graphql/internal/graph"

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

func RegisterGraphQLRoutes(
	hs *HTTPServer,
	schema *graph.ExecutableSchema,
) {
	mux := http.NewServeMux()

	// GraphQL handler
	srv := handler.NewDefaultServer(schema)
	mux.Handle("/query", srv)

	// GraphQL Playground
	mux.Handle("/playground", playground.Handler("GraphQL Playground", "/query"))

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"ok"}`)
	})

	hs.server = &http.Server{
		Addr:    hs.config.HTTPAddress,
		Handler: mux,
	}
}

func StartHTTPServer(lc fx.Lifecycle, hs *HTTPServer) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				hs.log.Info(fmt.Sprintf("Server running at http://%s:%s/playground", hs.config.HTTPAddress, hs.config.HTTPPort))
				if err := hs.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					hs.log.Error(fmt.Sprintf("Server error: %v", err))
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
