package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"

	"product/internal/domain/product"
	"product/pkg/store"
)

type ProductRepository struct {
	db *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) *ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

func (s *ProductRepository) Select(ctx context.Context) (dest []product.Entity, err error) {
	query := `
		SELECT id, name, description, cost, category
		FROM products
		ORDER BY id`

	dest = make([]product.Entity, 0)
	err = s.db.SelectContext(ctx, &dest, query)

	return
}

func (s *ProductRepository) Create(ctx context.Context, data product.Entity) (id string, err error) {
	query := `
		INSERT INTO products (name, description, cost, category)
		VALUES ($1, $2, $3, $4)
		RETURNING id`

	args := []any{data.Name, data.Description, data.Cost, data.Category}

	err = s.db.QueryRowContext(ctx, query, args...).Scan(&id)

	return
}

func (s *ProductRepository) Get(ctx context.Context, id string) (dest product.Entity, err error) {
	query := `
		SELECT id, name, description, cost, category
		FROM products
		WHERE id=$1`

	args := []any{id}

	if err = s.db.GetContext(ctx, &dest, query, args...); err != nil && err != sql.ErrNoRows {
		return
	}

	if err == sql.ErrNoRows {
		err = store.ErrorNotFound
	}

	return
}

func (s *ProductRepository) Update(ctx context.Context, id string, data product.Entity) (err error) {
	sets, args := s.prepareArgs(data)
	if len(args) > 0 {

		args = append(args, id)
		sets = append(sets, "updated_at=CURRENT_TIMESTAMP")

		query := fmt.Sprintf("UPDATE products SET %s WHERE id=$%d", strings.Join(sets, ", "), len(args))
		_, err = s.db.ExecContext(ctx, query, args...)
	}

	return
}

func (s *ProductRepository) prepareArgs(data product.Entity) (sets []string, args []any) {
	if data.Name != nil {
		args = append(args, data.Name)
		sets = append(sets, fmt.Sprintf("name=$%d", len(args)))
	}

	if data.Description != nil {
		args = append(args, data.Description)
		sets = append(sets, fmt.Sprintf("genre=$%d", len(args)))
	}

	if *data.Cost > 0 {
		args = append(args, data.Cost)
		sets = append(sets, fmt.Sprintf("isbn=$%d", len(args)))
	}

	if data.Category != nil {
		args = append(args, data.Category)
		sets = append(sets, fmt.Sprintf("authors=$%d", len(args)))
	}

	return
}

func (s *ProductRepository) Delete(ctx context.Context, id string) (err error) {
	query := `
		DELETE 
		FROM products
		WHERE id=$1`

	args := []any{id}

	_, err = s.db.ExecContext(ctx, query, args...)

	return
}
