package dto

// CreateProductRequest represents the data needed to create a product.
type CreateProductRequest struct {
	Name        string `json:"name"`
	Brand       string `json:"brand"`
	Description string `json:"description"`
	PriceMinor  int64  `json:"price_minor"`
	CategoryID  uint   `json:"category_id"`
}

// UpdateProductRequest represents the data needed to update a product.
type UpdateProductRequest struct {
	Name        string `json:"name"`
	Brand       string `json:"brand"`
	Description string `json:"description"`
	PriceMinor  int64  `json:"price_minor"`
	CategoryID  uint   `json:"category_id"`
	IsActive    bool   `json:"is_active"`
}

// ProductResponse is the product information returned to the client.
type ProductResponse struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	Slug         string `json:"slug"`
	Brand        string `json:"brand"`
	Description  string `json:"description"`
	PriceMinor   int64  `json:"price_minor"`
	CategoryID   uint   `json:"category_id"`
	CategoryName string `json:"category_name"`
	IsActive     bool   `json:"is_active"`
}

// PaginationResponse describes the current page.
type PaginationResponse struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// ProductListResponse contains products and pagination information.
type ProductListResponse struct {
	Data       []ProductResponse  `json:"data"`
	Pagination PaginationResponse `json:"pagination"`
}
