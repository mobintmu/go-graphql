package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go-graphql/internal/config"
	"go-graphql/internal/graph/resolvers"

	// gqlgen generated package
	product "go-graphql/internal/product/service"

	"github.com/gin-contrib/timeout"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type HTTPServer struct {
	server *http.Server
	config *config.Config
	log    *zap.Logger
}

// NewGraphQLResolver wires your services into the gqlgen resolvers.
func NewGraphQLResolver(productSvc *product.Product) *resolvers.Resolver {
	return &resolvers.Resolver{
		ProductService: productSvc,
	}
}

func NewGinEngine() *gin.Engine {
	if gin.Mode() != gin.ReleaseMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Logger(),
		gin.Recovery(),
		timeout.New(timeout.WithTimeout(60*time.Second)))

	return r
}

func NewHTTPServer(engine *gin.Engine, cfg *config.Config, logger *zap.Logger) *HTTPServer {
	httpServer := http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.HTTPAddress, cfg.HTTPPort),
		Handler: engine,
	}
	return &HTTPServer{
		server: &httpServer,
		config: cfg,
		log:    logger,
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
