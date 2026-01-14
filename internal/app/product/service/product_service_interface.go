package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/jamal23041989/go-marketplace-inventory-service/internal/core/domain"
)

type ProductService interface {
	Create(ctx context.Context, p domain.Product) (domain.Product, error)
	GetById(ctx context.Context, id uuid.UUID) (domain.Product, error)
	GetAll(ctx context.Context) ([]domain.Product, error)
	Update(ctx context.Context, id uuid.UUID, dto domain.UpdateProductDTO) (domain.Product, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
