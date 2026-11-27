package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/habtamubaye-devsec/shoe-shop/internal/dto"
	"github.com/habtamubaye-devsec/shoe-shop/internal/models"
	"github.com/habtamubaye-devsec/shoe-shop/internal/repositories"
)

var (
	ErrInvalidInput = errors.New("invalid input")
	ErrNotFound     = errors.New("product not found")
)

type ProductListParams struct {
	Search     string
	CategoryID uint
	MinPrice   *int64
	MaxPrice   *int64
	Page       int
	Limit      int
}

type ProductService interface {
	Create(ctx context.Context, input dto.CreateProductRequest) (*dto.ProductResponse, error)
	GetByID(ctx context.Context, id uint) (*dto.ProductResponse, error)
	List(ctx context.Context, params ProductListParams) (*dto.ProductListResponse, error)
	Update(ctx context.Context, id uint, input dto.UpdateProductRequest) (*dto.ProductResponse, error)
	Delete(ctx context.Context, id uint) error
}

type productService struct {
	repo repositories.ProductRepository
}

func NewProductService(repo repositories.ProductRepository) ProductService {
	return &productService{
		repo: repo,
	}
}

// createSlug generates a simple URL-friendly slug from a product name.
func createSlug(name string) string {
	name = strings.ToLower(strings.TrimSpace(name))
	return strings.Join(strings.Fields(name), "-")
}

func (s *productService) ensureCategoryExists(
	ctx context.Context,
	categoryID uint,
) error {
	exists, err := s.repo.CategoryExists(ctx, categoryID)
	if err != nil {
		return fmt.Errorf("check category: %w", err)
	}

	if !exists {
		return fmt.Errorf("%w: category does not exist", ErrInvalidInput)
	}

	return nil
}

func (s *productService) ensureSlugAvailable(
	ctx context.Context,
	slug string,
	currentID uint,
) error {
	existing, err := s.repo.GetBySlug(ctx, slug)

	if errors.Is(err, repositories.ErrNotFound) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("check product slug: %w", err)
	}

	if existing.ID != currentID {
		return fmt.Errorf("%w: product name already exists", ErrInvalidInput)
	}

	return nil
}

func toProductResponse(product *models.Product) *dto.ProductResponse {
	return &dto.ProductResponse{
		ID:           product.ID,
		Name:         product.Name,
		Slug:         product.Slug,
		Brand:        product.Brand,
		Description:  product.Description,
		PriceMinor:   product.PriceMinor,
		CategoryID:   product.CategoryID,
		CategoryName: product.Category.Name,
		IsActive:     product.IsActive,
	}
}

func (s *productService) Create(
	ctx context.Context,
	input dto.CreateProductRequest,
) (*dto.ProductResponse, error) {
	input.Name = strings.TrimSpace(input.Name)
	input.Brand = strings.TrimSpace(input.Brand)

	if input.Name == "" {
		return nil, fmt.Errorf("%w: product name is required", ErrInvalidInput)
	}

	if input.PriceMinor < 0 {
		return nil, fmt.Errorf("%w: price cannot be negative", ErrInvalidInput)
	}

	if input.CategoryID == 0 {
		return nil, fmt.Errorf("%w: category_id is required", ErrInvalidInput)
	}

	if err := s.ensureCategoryExists(ctx, input.CategoryID); err != nil {
		return nil, err
	}

	slug := createSlug(input.Name)

	if err := s.ensureSlugAvailable(ctx, slug, 0); err != nil {
		return nil, err
	}

	product := &models.Product{
		Name:        input.Name,
		Slug:        slug,
		Brand:       input.Brand,
		Description: strings.TrimSpace(input.Description),
		PriceMinor:  input.PriceMinor,
		CategoryID:  input.CategoryID,
		IsActive:    true,
	}

	if err := s.repo.Create(ctx, product); err != nil {
		return nil, fmt.Errorf("create product: %w", err)
	}

	// Reload so the related category is available in the response.
	created, err := s.repo.GetByID(ctx, product.ID)
	if err != nil {
		return nil, fmt.Errorf("load created product: %w", err)
	}

	return toProductResponse(created), nil
}

func (s *productService) GetByID(
	ctx context.Context,
	id uint,
) (*dto.ProductResponse, error) {
	product, err := s.repo.GetByID(ctx, id)

	if errors.Is(err, repositories.ErrNotFound) {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("get product: %w", err)
	}

	return toProductResponse(product), nil
}

func (s *productService) List(
	ctx context.Context,
	params ProductListParams,
) (*dto.ProductListResponse, error) {
	if params.Page < 1 {
		params.Page = 1
	}

	if params.Limit < 1 {
		params.Limit = 10
	}

	if params.Limit > 100 {
		params.Limit = 100
	}

	if params.MinPrice != nil && *params.MinPrice < 0 {
		return nil, fmt.Errorf("%w: minimum price cannot be negative", ErrInvalidInput)
	}

	if params.MaxPrice != nil && *params.MaxPrice < 0 {
		return nil, fmt.Errorf("%w: maximum price cannot be negative", ErrInvalidInput)
	}

	if params.MinPrice != nil && params.MaxPrice != nil &&
		*params.MinPrice > *params.MaxPrice {
		return nil, fmt.Errorf("%w: minimum price cannot exceed maximum price", ErrInvalidInput)
	}

	maxInt := int(^uint(0) >> 1)

	if params.Page-1 > maxInt/params.Limit {
		return nil, fmt.Errorf("%w: page is too large", ErrInvalidInput)
	}

	offset := (params.Page - 1) * params.Limit

	products, total, err := s.repo.List(
		ctx,
		repositories.ProductFilter{
			Search:     strings.TrimSpace(params.Search),
			CategoryID: params.CategoryID,
			MinPrice:   params.MinPrice,
			MaxPrice:   params.MaxPrice,
			Limit:      params.Limit,
			Offset:     offset,
		},
	)

	if err != nil {
		return nil, fmt.Errorf("list products: %w", err)
	}

	data := make([]dto.ProductResponse, 0, len(products))

	for i := range products {
		data = append(data, *toProductResponse(&products[i]))
	}

	totalPages := int(total / int64(params.Limit))
	if total%int64(params.Limit) != 0 {
		totalPages++
	}

	return &dto.ProductListResponse{
		Data: data,
		Pagination: dto.PaginationResponse{
			Page:       params.Page,
			Limit:      params.Limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *productService) Update(
	ctx context.Context,
	id uint,
	input dto.UpdateProductRequest,
) (*dto.ProductResponse, error) {
	input.Name = strings.TrimSpace(input.Name)
	input.Brand = strings.TrimSpace(input.Brand)

	if input.Name == "" {
		return nil, fmt.Errorf("%w: product name is required", ErrInvalidInput)
	}

	if input.PriceMinor < 0 {
		return nil, fmt.Errorf("%w: price cannot be negative", ErrInvalidInput)
	}

	if input.CategoryID == 0 {
		return nil, fmt.Errorf("%w: category_id is required", ErrInvalidInput)
	}

	product, err := s.repo.GetByID(ctx, id)

	if errors.Is(err, repositories.ErrNotFound) {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("find product for update: %w", err)
	}

	if err := s.ensureCategoryExists(ctx, input.CategoryID); err != nil {
		return nil, err
	}

	newSlug := createSlug(input.Name)

	if err := s.ensureSlugAvailable(ctx, newSlug, id); err != nil {
		return nil, err
	}

	product.Name = input.Name
	product.Slug = newSlug
	product.Brand = input.Brand
	product.Description = strings.TrimSpace(input.Description)
	product.PriceMinor = input.PriceMinor
	product.CategoryID = input.CategoryID
	product.IsActive = input.IsActive

	if err := s.repo.Update(ctx, product); err != nil {
		return nil, fmt.Errorf("update product: %w", err)
	}

	updated, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("reload updated product: %w", err)
	}

	return toProductResponse(updated), nil
}

func (s *productService) Delete(
	ctx context.Context,
	id uint,
) error {
	err := s.repo.Delete(ctx, id)

	if errors.Is(err, repositories.ErrNotFound) {
		return ErrNotFound
	}

	if err != nil {
		return fmt.Errorf("delete product: %w", err)
	}

	return nil
}
