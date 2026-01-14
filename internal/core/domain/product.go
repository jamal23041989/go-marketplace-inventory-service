package domain

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID          uuid.UUID
	Name        string
	Description string
	Price       int64
	Quantity    int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type UpdateProductDTO struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Price       *int64  `json:"price,omitempty"`
	Quantity    *int    `json:"quantity,omitempty"`
}
