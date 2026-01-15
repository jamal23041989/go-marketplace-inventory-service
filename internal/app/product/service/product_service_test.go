package service

import (
	"errors"
	"testing"

	"github.com/jamal23041989/go-marketplace-inventory-service/internal/core/domain"
	ers "github.com/jamal23041989/go-marketplace-inventory-service/internal/core/errors"
)

func TestValidateProduct_Success(t *testing.T) {
	p := &productService{}

	validProduct := domain.Product{
		Name:        "Клавиатура",
		Price:       1500,
		Quantity:    10,
		Description: "Механическая клавиатура с подсветкой",
	}

	err := p.validateProduct(validProduct)
	if err != nil {
		t.Fatalf("product validation failed: %s", err)
	}
}

func TestValidateProduct_NegativePrice(t *testing.T) {
	p := &productService{}

	invalidProduct := domain.Product{
		Name:     "Мышка",
		Price:    -100,
		Quantity: 5,
	}

	err := p.validateProduct(invalidProduct)
	if err == nil {
		t.Fatalf("expected error for negative price, but got nil")
	}

	if !errors.Is(err, ers.ErrInvalidInput) {
		t.Errorf("expected error to be %v, but got %v", ers.ErrInvalidInput, err)
	}
}
