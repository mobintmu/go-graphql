package server

import (
	"fmt"
	"go-graphql/docs"
	"go-graphql/internal/config"
	"go-graphql/internal/graph/generated"
	"go-graphql/internal/graph/resolvers"
	"go-graphql/internal/health"
	"go-graphql/internal/product/controller"
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RegisterRoutes(
	engine *gin.Engine,
	health *health.Health,
	cfg *config.Config,
	adminProduct *controller.AdminProduct,
	clientProduct *controller.ClientProduct,
	resolver *resolvers.Resolver,
) {
	log.Println("🚀 Registering routes...")

	// Health check
	engine.GET("/health", health.Handle)

	// Admin Product routes
	adminGroup := engine.Group("/api/v1/admin/products")
	adminProduct.RegisterRoutes(adminGroup, cfg)

	// Client Product routes
	clientGroup := engine.Group("/api/v1/products")
	clientProduct.RegisterRoutes(clientGroup)

	// GraphQL schema + handler
	schema := generated.NewExecutableSchema(generated.Config{Resolvers: resolver})
	graphqlHandler := handler.NewDefaultServer(schema)

	// GraphQL endpoints
	engine.POST("/query", gin.WrapH(graphqlHandler))
	engine.GET("/playground", gin.WrapH(playground.Handler("GraphQL Playground", "/query")))

	// Swagger docs
	docs.SwaggerInfo.Title = "My API"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Description = "This is a sample API with Gin and Swagger."
	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%d", cfg.HTTPAddress, cfg.HTTPPort)
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http"} // or {"https"} in production
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Handle 404 for unknown routes
	engine.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "The requested endpoint does not exist",
			"path":    c.Request.URL.Path,
		})
	})
}
