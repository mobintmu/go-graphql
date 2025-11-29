package resolvers

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require
// here.

import product "go-graphql/internal/product/service"

type Resolver struct {
	Product *product.Product
}
