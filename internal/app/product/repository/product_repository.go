package repository

import (
	"context"
	"errors"
	"sync"

	"github.com/google/uuid"
	"github.com/jamal23041989/go-marketplace-inventory-service/internal/core/domain"
	ers "github.com/jamal23041989/go-marketplace-inventory-service/internal/core/errors"
)

type InMemoryProductRepository struct {
	mu       sync.RWMutex
	products map[uuid.UUID]domain.Product
}

func NewInMemoryProductRepository() *InMemoryProductRepository {
	return &InMemoryProductRepository{
		products: make(map[uuid.UUID]domain.Product),
	}
}

func (i *InMemoryProductRepository) Create(ctx context.Context, p domain.Product) (domain.Product, error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	if _, ok := i.products[p.ID]; ok {
		return domain.Product{}, errors.New("product already exists")
	}

	i.products[p.ID] = p
	return p, nil
}

func (i *InMemoryProductRepository) GetById(ctx context.Context, id uuid.UUID) (domain.Product, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	if product, ok := i.products[id]; ok {
		return product, nil
	}
	return domain.Product{}, ers.ErrProductNotFound
}

func (i *InMemoryProductRepository) GetAll(ctx context.Context) ([]domain.Product, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()
	products := make([]domain.Product, 0, len(i.products))
	for _, product := range i.products {
		products = append(products, product)
	}
	return products, nil
}

func (i *InMemoryProductRepository) Update(
	ctx context.Context,
	id uuid.UUID,
	p domain.Product,
) (domain.Product, error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	product, ok := i.products[id]
	if !ok {
		return domain.Product{}, ers.ErrProductNotFound
	}

	i.products[id] = p
	return product, nil
}

func (i *InMemoryProductRepository) Delete(ctx context.Context, id uuid.UUID) error {
	i.mu.Lock()
	defer i.mu.Unlock()
	delete(i.products, id)
	return nil
}
