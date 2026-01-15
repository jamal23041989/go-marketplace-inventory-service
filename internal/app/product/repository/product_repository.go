package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jamal23041989/go-marketplace-inventory-service/internal/core/domain"
	ers "github.com/jamal23041989/go-marketplace-inventory-service/internal/core/errors"
)

type PostgresProductRepository struct {
	db *sql.DB
}

func NewPostgresProductRepository(db *sql.DB) *PostgresProductRepository {
	return &PostgresProductRepository{
		db: db,
	}
}

func (i *PostgresProductRepository) Create(
	ctx context.Context,
	p *domain.Product,
) (domain.Product, error) {
	query := `
		INSERT INTO products (name, description, price, quantity, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	if err := i.db.QueryRowContext(
		ctx,
		query,
		p.Name,
		p.Description,
		p.Price,
		p.Quantity,
		p.CreatedAt,
		p.UpdatedAt,
	).Scan(&p.ID); err != nil {
		return domain.Product{}, fmt.Errorf("error inserting product: %w", err)
	}

	return *p, nil
}

func (i *PostgresProductRepository) GetById(
	ctx context.Context,
	id uuid.UUID,
) (domain.Product, error) {
	var product domain.Product

	query := `
		SELECT id, name, description, price, quantity, created_at, updated_at
		FROM products
		WHERE id = $1
	`

	if err := i.db.QueryRowContext(ctx, query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.Quantity,
		&product.CreatedAt,
		&product.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Product{}, fmt.Errorf("%w: product not found", ers.ErrProductNotFound)
		}
		return domain.Product{}, err
	}

	return product, nil
}

func (i *PostgresProductRepository) GetAll(
	ctx context.Context,
) ([]domain.Product, error) {
	var products []domain.Product

	query := `
		SELECT id, name, description, price, quantity, created_at, updated_at
		FROM products
	`

	rows, err := i.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		if err := rows.Close(); err != nil {
		}
	}(rows)

	for rows.Next() {
		var product domain.Product

		if err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.Quantity,
			&product.CreatedAt,
			&product.UpdatedAt,
		); err != nil {
			return products, err
		}

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return products, nil
}

func (i *PostgresProductRepository) Update(
	ctx context.Context,
	id uuid.UUID,
	p domain.Product,
) (domain.Product, error) {
	query := `
       UPDATE products
       SET name = $1, description = $2, price = $3, quantity = $4, updated_at = $5
       WHERE id = $6
       RETURNING id, name, description, price, quantity, created_at, updated_at
    `

	var updatedProduct domain.Product
	err := i.db.QueryRowContext(
		ctx,
		query,
		p.Name, p.Description, p.Price, p.Quantity, p.UpdatedAt, id,
	).Scan(
		&updatedProduct.ID,
		&updatedProduct.Name,
		&updatedProduct.Description,
		&updatedProduct.Price,
		&updatedProduct.Quantity,
		&updatedProduct.CreatedAt,
		&updatedProduct.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Product{}, fmt.Errorf("%w: not found error", ers.ErrProductNotFound)
		}
		return domain.Product{}, fmt.Errorf("error updating product: %w", err)
	}

	return updatedProduct, nil
}

func (i *PostgresProductRepository) Delete(
	ctx context.Context,
	id uuid.UUID,
) error {
	query := `
		DELETE FROM products
		WHERE id = $1
	`

	result, err := i.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	_, err = result.RowsAffected()
	if err != nil {
		return err
	}

	return nil
}
