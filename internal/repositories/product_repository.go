package repositories

import (
	"context"
	"errors"

	"github.com/habtamubaye-devsec/shoe-shop/internal/models"

	"gorm.io/gorm"
)

var ErrNotFound = errors.New("record not found")

type ProductFilter struct {
	Search     string
	CategoryID uint
	MinPrice   *int64
	MaxPrice   *int64
	Limit      int
	Offset     int
}

type ProductRepository interface {
	Create(ctx context.Context, product *models.Product) error
	GetByID(ctx context.Context, id uint) (*models.Product, error)
	GetBySlug(ctx context.Context, slug string) (*models.Product, error)
	List(ctx context.Context, filter ProductFilter) ([]models.Product, int64, error)
	Update(ctx context.Context, product *models.Product) error
	Delete(ctx context.Context, id uint) error
	CategoryExists(ctx context.Context, categoryID uint) (bool, error)
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{
		db: db,
	}
}

func (r *productRepository) Create(
	ctx context.Context,
	product *models.Product,
) error {
	return r.db.WithContext(ctx).Create(product).Error
}

func (r *productRepository) GetByID(
	ctx context.Context,
	id uint,
) (*models.Product, error) {
	var product models.Product

	err := r.db.WithContext(ctx).
		Preload("Category").
		First(&product, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *productRepository) GetBySlug(
	ctx context.Context,
	slug string,
) (*models.Product, error) {
	var product models.Product

	err := r.db.WithContext(ctx).
		First(&product, "slug = ?", slug).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *productRepository) List(
	ctx context.Context,
	filter ProductFilter,
) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Product{})

	if filter.Search != "" {
		search := "%" + filter.Search + "%"

		query = query.Where(
			"(name ILIKE ? OR brand ILIKE ?)",
			search,
			search,
		)
	}

	if filter.CategoryID != 0 {
		query = query.Where("category_id = ?", filter.CategoryID)
	}

	if filter.MinPrice != nil {
		query = query.Where("price_minor >= ?", *filter.MinPrice)
	}

	if filter.MaxPrice != nil {
		query = query.Where("price_minor <= ?", *filter.MaxPrice)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Preload("Category").
		Order("created_at DESC").
		Limit(filter.Limit).
		Offset(filter.Offset).
		Find(&products).Error

	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *productRepository) Update(
	ctx context.Context,
	product *models.Product,
) error {
	return r.db.WithContext(ctx).Save(product).Error
}

func (r *productRepository) Delete(
	ctx context.Context,
	id uint,
) error {
	result := r.db.WithContext(ctx).Delete(&models.Product{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *productRepository) CategoryExists(
	ctx context.Context,
	categoryID uint,
) (bool, error) {
	var count int64

	err := r.db.WithContext(ctx).
		Model(&models.Category{}).
		Where("id = ?", categoryID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
