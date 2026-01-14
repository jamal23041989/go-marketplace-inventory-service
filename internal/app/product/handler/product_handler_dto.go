package handler

import (
	"github.com/jamal23041989/go-marketplace-inventory-service/internal/core/domain"
)

type CreateProductRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int64  `json:"price"`
	Quantity    int    `json:"quantity"`
}

func (r *CreateProductRequest) ToDomain() domain.Product {
	var product domain.Product

	product.Name = r.Name
	product.Description = r.Description
	product.Price = r.Price
	product.Quantity = r.Quantity

	return product
}

type UpdateProductRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Price       *int64  `json:"price"`
	Quantity    *int    `json:"quantity"`
}

func (r *UpdateProductRequest) ToUpdateDTO() domain.UpdateProductDTO {
	return domain.UpdateProductDTO{
		Name:        r.Name,
		Description: r.Description,
		Price:       r.Price,
		Quantity:    r.Quantity,
	}
}
