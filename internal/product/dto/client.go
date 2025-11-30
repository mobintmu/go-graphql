package dto

type ProductResponse struct {
	ID          int32  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int64  `json:"price"`
	IsActive    bool   `json:"isActive,omitempty"`
}

type ClientListProductsResponse []ProductResponse
