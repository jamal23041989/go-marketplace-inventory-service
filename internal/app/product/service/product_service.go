package service

import (
	"context"
	"errors"
	"fmt"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/jamal23041989/go-marketplace-inventory-service/internal/app/product/repository"
	"github.com/jamal23041989/go-marketplace-inventory-service/internal/core/domain"
	ers "github.com/jamal23041989/go-marketplace-inventory-service/internal/core/errors"
)

type productService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{
		repo: repo,
	}
}

func (p *productService) Create(ctx context.Context, product domain.Product) (domain.Product, error) {
	if err := p.validateProduct(product); err != nil {
		return domain.Product{}, err
	}

	product.ID = uuid.New()
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()

	return p.repo.Create(ctx, product)
}

func (p *productService) GetById(ctx context.Context, id uuid.UUID) (domain.Product, error) {
	if id == uuid.Nil {
		return domain.Product{}, errors.New("invalid product id")
	}
	return p.repo.GetById(ctx, id)
}

func (p *productService) GetAll(ctx context.Context) ([]domain.Product, error) {
	return p.repo.GetAll(ctx)
}

func (p *productService) Update(
	ctx context.Context,
	id uuid.UUID,
	dto domain.UpdateProductDTO,
) (domain.Product, error) {
	if id == uuid.Nil {
		return domain.Product{}, errors.New("invalid product id")
	}

	currentProduct, err := p.repo.GetById(ctx, id)
	if err != nil {
		return domain.Product{}, err
	}

	if dto.Name != nil {
		currentProduct.Name = *dto.Name
	}
	if dto.Description != nil {
		currentProduct.Description = *dto.Description
	}
	if dto.Price != nil {
		currentProduct.Price = *dto.Price
	}
	if dto.Quantity != nil {
		currentProduct.Quantity = *dto.Quantity
	}
	currentProduct.UpdatedAt = time.Now()

	if err := p.validateProduct(currentProduct); err != nil {
		return domain.Product{}, err
	}

	return p.repo.Update(ctx, id, currentProduct)
}

func (p *productService) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("invalid product id")
	}
	return p.repo.Delete(ctx, id)
}

func (p *productService) validateProduct(product domain.Product) error {
	if product.Name == "" {
		return fmt.Errorf("%w: product name is required", ers.ErrInvalidInput)
	}
	if product.Price < 0 {
		return fmt.Errorf("%w: product price cannot be negative", ers.ErrInvalidInput)
	}
	if product.Quantity < 0 {
		return fmt.Errorf("%w: product quantity cannot be negative", ers.ErrInvalidInput)
	}
	if utf8.RuneCountInString(product.Description) == 0 || utf8.RuneCountInString(product.Description) > 500 {
		return fmt.Errorf("%w: product description is too large", ers.ErrInvalidInput)
	}
	return nil
}
