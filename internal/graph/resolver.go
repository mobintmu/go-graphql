package graph

import (
	"go-graphql/internal/product/service"
)

// This file will be shared across your resolvers

type Resolver struct {
	ProductService *service.Product
}
