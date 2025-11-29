package resolvers

import (
	"context"

	"github.com/mobintmu/go-simple/internal/graph/model"
)

// Query Resolvers
func (r *queryResolver) GetProduct(
	ctx context.Context,
	id string,
) (*model.Product, error) {
	product, err := r.Product.GetProductByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &model.Product{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}, nil
}

func (r *queryResolver) ListProducts(
	ctx context.Context,
	limit *int,
	offset *int,
) (*model.ProductList, error) {
	l := 10
	if limit != nil {
		l = *limit
	}
	o := 0
	if offset != nil {
		o = *offset
	}

	products, err := r.Product.ListProducts(ctx)
	if err != nil {
		return nil, err
	}

	items := make([]*model.Product, len(products))
	for i, p := range products {
		items[i] = &model.Product{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			CreatedAt:   p.CreatedAt,
			UpdatedAt:   p.UpdatedAt,
		}
	}

	return &model.ProductList{
		Items:  items,
		Total:  total,
		Limit:  l,
		Offset: o,
	}, nil
}

// Mutation Resolvers
func (r *mutationResolver) CreateProduct(
	ctx context.Context,
	input model.CreateProductInput,
) (*model.Product, error) {
	product, err := r.Product.Create(ctx, &service.CreateProductInput{
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
	})
	if err != nil {
		return nil, err
	}

	return &model.Product{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}, nil
}

func (r *mutationResolver) UpdateProduct(
	ctx context.Context,
	id string,
	input model.UpdateProductInput,
) (*model.Product, error) {
	product, err := r.Product.Update(ctx, id, &service.UpdateProductInput{
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
	})
	if err != nil {
		return nil, err
	}

	return &model.Product{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}, nil
}

func (r *mutationResolver) DeleteProduct(
	ctx context.Context,
	id string,
) (bool, error) {
	return r.Product.Delete(ctx, id)
}
