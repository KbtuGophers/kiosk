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
		SELECT id, category_id, barcode, name, measure, producer_country, brand_name, description, image, is_weighted
		FROM products
		ORDER BY id`

	dest = make([]product.Entity, 0)
	err = s.db.SelectContext(ctx, &dest, query)

	return
}

func (s *ProductRepository) Create(ctx context.Context, data product.Entity) (id string, err error) {
	query := `
		INSERT INTO products (id,category_id, barcode, name, measure, producer_country, brand_name, description, image, is_weighted)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id`

	args := []any{data.ID, data.CategoryID, data.Barcode, data.Name, data.Measure, data.ProducerCountry,
		data.BrandName, data.Description, data.Image, data.IsWeighted}

	err = s.db.QueryRowContext(ctx, query, args...).Scan(&id)

	return
}

func (s *ProductRepository) Get(ctx context.Context, id string) (dest product.Entity, err error) {
	query := `
		SELECT id, category_id, barcode, name, measure, producer_country, brand_name, description, image, is_weighted
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
	if data.Barcode != nil {
		args = append(args, data.Barcode)
		sets = append(sets, fmt.Sprintf("barcode=$%d", len(args)))
	}

	if data.Name != nil {
		args = append(args, data.Name)
		sets = append(sets, fmt.Sprintf("name=$%d", len(args)))
	}

	if data.Measure != nil {
		args = append(args, data.Measure)
		sets = append(sets, fmt.Sprintf("measure=$%d", len(args)))
	}

	if data.ProducerCountry != nil {
		args = append(args, data.ProducerCountry)
		sets = append(sets, fmt.Sprintf("producer_country=$%d", len(args)))
	}

	if data.BrandName != nil {
		args = append(args, data.BrandName)
		sets = append(sets, fmt.Sprintf("brand_name=$%d", len(args)))
	}

	if data.Description != nil {
		args = append(args, data.Description)
		sets = append(sets, fmt.Sprintf("description=$%d", len(args)))
	}

	if data.Image != nil {
		args = append(args, data.Image)
		sets = append(sets, fmt.Sprintf("image=$%d", len(args)))
	}

	if data.IsWeighted != nil {
		args = append(args, data.IsWeighted)
		sets = append(sets, fmt.Sprintf("is_weighted=$%d", len(args)))
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
