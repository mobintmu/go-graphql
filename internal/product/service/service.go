package service

import (
	"context"
	"database/sql"
	"go-graphql/internal/config"
	"go-graphql/internal/graph/model"
	"go-graphql/internal/pkg/logger/utils"
	"go-graphql/internal/product/dto"
	"go-graphql/internal/storage/cache"
	"go-graphql/internal/storage/sql/sqlc"

	"go.uber.org/zap"
)

type Product struct {
	query  *sqlc.Queries
	log    *zap.Logger
	memory *cache.Store
	cfg    *config.Config
}

func New(q *sqlc.Queries,
	log *zap.Logger,
	memory *cache.Store,
	cfg *config.Config) *Product {
	return &Product{
		query:  q,
		log:    log,
		memory: memory,
		cfg:    cfg,
	}
}

func (s *Product) Create(ctx context.Context, req dto.AdminCreateProductRequest) (dto.ProductResponse, error) {
	arg := sqlc.CreateProductParams{
		ProductName:        req.Name,
		ProductDescription: req.Description,
		Price:              req.Price,
		IsActive:           true,
	}
	product, err := s.query.CreateProduct(ctx, arg)
	if err != nil {
		return dto.ProductResponse{}, err
	}
	s.log.Info("Product created", zap.Int32("id", product.ID))
	s.memory.Set(ctx, s.memory.KeyProduct(product.ID), product, s.cfg.Redis.DefaultTTL)
	s.memory.Delete(ctx, s.memory.KeyAllProducts())
	return dto.ProductResponse{
		ID:          product.ID,
		Name:        product.ProductName,
		Description: product.ProductDescription,
		Price:       product.Price,
	}, nil
}

func (s *Product) Update(ctx context.Context, req dto.AdminUpdateProductRequest) (dto.ProductResponse, error) {
	arg := sqlc.UpdateProductParams{
		ID:                 int32(req.ID),
		ProductName:        req.Name,
		ProductDescription: req.Description,
		Price:              req.Price,
		IsActive:           req.IsActive,
	}
	product, err := s.query.UpdateProduct(ctx, arg)
	if err != nil {
		return dto.ProductResponse{}, err
	}
	s.memory.Set(ctx, s.memory.KeyProduct(product.ID), product, s.cfg.Redis.DefaultTTL)
	s.memory.Delete(ctx, s.memory.KeyAllProducts())
	return dto.ProductResponse{
		ID:          product.ID,
		Name:        product.ProductName,
		Description: product.ProductDescription,
		Price:       product.Price,
	}, nil
}

func (s *Product) Delete(ctx context.Context, id int32) error {
	s.memory.Delete(ctx, s.memory.KeyProduct(id))
	s.memory.Delete(ctx, s.memory.KeyAllProducts())
	return s.query.DeleteProduct(ctx, id)
}

func (s *Product) GetProductByID(ctx context.Context, id int32) (dto.ProductResponse, error) {
	var product sqlc.Product
	err := s.memory.Get(ctx, s.memory.KeyProduct(id), &product)
	if err != nil {
		product, err = s.query.GetProduct(ctx, id)
		if err != nil {
			return dto.ProductResponse{}, err
		}
		s.memory.Set(ctx, s.memory.KeyProduct(product.ID), product, s.cfg.Redis.DefaultTTL)
	}
	result := dto.ProductResponse{
		ID:          product.ID,
		Name:        product.ProductName,
		Description: product.ProductDescription,
		Price:       product.Price,
	}
	return result, nil
}

func (s *Product) ListProductsWithoutFilter(ctx context.Context) (dto.ClientListProductsResponse, error) {
	var resp []dto.ProductResponse
	if err := s.memory.Get(ctx, s.memory.KeyAllProducts(), &resp); err == nil {
		return resp, nil
	}
	products, err := s.query.ListProducts(ctx)
	if err != nil {
		return nil, err
	}
	resp = make([]dto.ProductResponse, 0, len(products))
	for _, product := range products {
		resp = append(resp, dto.ProductResponse{
			ID:          product.ID,
			Name:        product.ProductName,
			Description: product.ProductDescription,
			Price:       product.Price,
		})
	}
	s.memory.Set(ctx, s.memory.KeyAllProducts(), resp, s.cfg.Redis.DefaultTTL)
	return resp, nil
}

func (s *Product) ListProducts(ctx context.Context, filter *model.ProductFilter, pagination *model.PaginationInput) (*model.ProductConnection, error) {
	params := s.graphqlFilterToSQLCParams(filter, pagination)
	products, err := s.query.ListProductsWithFilters(ctx, params)
	if err != nil {
		s.log.Error("Failed to list products", zap.Error(err))
		return nil, err
	}

	paramsCount := s.graphqlCountProductsFilterToSQLCParams(filter)
	total, err := s.query.CountProductsWithFilters(ctx, paramsCount)
	if err != nil {
		s.log.Error("Failed to count products", zap.Error(err))
		return nil, err
	}

	var result []*model.Product
	for _, p := range products {
		result = append(result, &model.Product{
			ID:          int(p.ID),
			Name:        p.ProductName,
			Description: p.ProductDescription,
			Price:       p.Price,
			IsActive:    p.IsActive,
		})
	}
	return &model.ProductConnection{
		Products: result,
		Total:    int(total),
	}, nil
}

// graphqlFilterToSQLCParams converts GraphQL ProductFilter to SQLC params
func (s *Product) graphqlFilterToSQLCParams(
	filter *model.ProductFilter,
	pagination *model.PaginationInput,
) sqlc.ListProductsWithFiltersParams {
	// Initialize with defaults
	params := sqlc.ListProductsWithFiltersParams{
		Limit:  sql.NullInt64{Int64: 10, Valid: true},
		Offset: sql.NullInt64{Int64: 0, Valid: true},
	}

	// Apply filter if provided
	if filter != nil {
		params.ID = utils.ToNullInt32(filter.ID)
		params.ProductName = utils.ToNullString(filter.Name)
		params.ProductDescription = utils.ToNullString(filter.Description)
		params.MinPrice = utils.ToNullInt64(filter.MinPrice)
		params.MaxPrice = utils.ToNullInt64(filter.MaxPrice)
		params.IsActive = utils.ToNullBool(filter.IsActive)
	}

	if pagination != nil {
		if pagination.Limit != nil {
			params.Limit = sql.NullInt64{Int64: int64(*pagination.Limit), Valid: true}
		}
		if pagination.Offset != nil {
			params.Offset = sql.NullInt64{Int64: int64(*pagination.Offset), Valid: true}
		}
	}

	return params
}

func (s *Product) graphqlCountProductsFilterToSQLCParams(
	filter *model.ProductFilter,
) sqlc.CountProductsWithFiltersParams {
	// Initialize with defaults
	params := sqlc.CountProductsWithFiltersParams{}

	// Apply filter if provided
	if filter != nil {
		params.ID = utils.ToNullInt32(filter.ID)
		params.ProductName = utils.ToNullString(filter.Name)
		params.ProductDescription = utils.ToNullString(filter.Description)
		params.MinPrice = utils.ToNullInt64(filter.MinPrice)
		params.MaxPrice = utils.ToNullInt64(filter.MaxPrice)
		params.IsActive = utils.ToNullBool(filter.IsActive)
	}

	return params
}
