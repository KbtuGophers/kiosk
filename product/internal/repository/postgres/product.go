package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
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

func (s *ProductRepository) Select(ctx context.Context, r *http.Request) (dest []product.Entity, err error) {
	filters, args := s.prepareFilters(r)
	query := fmt.Sprintf(`SELECT id, category_id, barcode, name, measure, cost, producer_country, brand_name, description, image, is_weighted `+
		`FROM products WHERE %s`, strings.Join(filters, " "))
	query += " 1=1"

	fmt.Println(query, args)

	dest = make([]product.Entity, 0)
	err = s.db.SelectContext(ctx, &dest, query, args...)

	return
}

func (s *ProductRepository) prepareFilters(r *http.Request) (filters []string, args []any) {
	priceGTEFilter, err := strconv.Atoi(r.URL.Query().Get("cost_gte"))
	if err == nil {
		args = append(args, priceGTEFilter)
		filters = append(filters, fmt.Sprintf("cost >= $%d AND", len(args)))
	}

	priceLTEFilter, err := strconv.Atoi(r.URL.Query().Get("cost_lte"))
	if err == nil {
		args = append(args, priceLTEFilter)
		filters = append(filters, fmt.Sprintf("cost <= $%d AND", len(args)))
	}

	searchFilter := strings.ToLower(r.URL.Query().Get("search"))
	if searchFilter != "" {
		args = append(args, "%"+searchFilter+"%")
		filters = append(filters, fmt.Sprintf("name LIKE $%d AND", len(args)))
	}
	return
}

func (s *ProductRepository) Create(ctx context.Context, data product.Entity) (id string, err error) {
	query := `
		INSERT INTO products (id,category_id, barcode, name, measure, cost, producer_country, brand_name, description, image, is_weighted)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id`

	args := []any{data.ID, data.CategoryID, data.Barcode, data.Name, data.Measure, data.Cost, data.ProducerCountry,
		data.BrandName, data.Description, data.Image, data.IsWeighted}

	err = s.db.QueryRowContext(ctx, query, args...).Scan(&id)

	return
}

func (s *ProductRepository) Get(ctx context.Context, id string) (dest product.Entity, err error) {
	query := `
		SELECT id, category_id, barcode, name, measure, cost, producer_country, brand_name, description, image, is_weighted
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

	if data.Cost != nil {
		args = append(args, data.Cost)
		sets = append(sets, fmt.Sprintf("cost=$%d", len(args)))
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
