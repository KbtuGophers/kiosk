package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	category "product/internal/domain/category"
	"product/pkg/store"
	"strings"
)

type CategoryRepository struct {
	db *sqlx.DB
}

func NewCategoryRepository(db *sqlx.DB) *CategoryRepository {
	return &CategoryRepository{
		db: db,
	}
}

func (s *CategoryRepository) Select(ctx context.Context) (dest []category.Entity, err error) {
	query := `
		SELECT id, name
		FROM category
		ORDER BY id`

	err = s.db.SelectContext(ctx, &dest, query)

	return
}

func (s *CategoryRepository) Create(ctx context.Context, data category.Entity) (id string, err error) {
	query := `
		INSERT INTO category (name)
		VALUES ($1)
		RETURNING id`

	args := []any{data.Name}

	err = s.db.QueryRowContext(ctx, query, args...).Scan(&id)

	return
}

func (s *CategoryRepository) Get(ctx context.Context, id string) (dest category.Entity, err error) {
	query := `
		SELECT id, name
		FROM category
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

func (s *CategoryRepository) Update(ctx context.Context, id string, data category.Entity) (err error) {
	sets, args := s.prepareArgs(data)
	if len(args) > 0 {

		args = append(args, id)
		sets = append(sets, "updated_at=CURRENT_TIMESTAMP")

		query := fmt.Sprintf("UPDATE category SET %s WHERE id=$%d", strings.Join(sets, ", "), len(args))
		_, err = s.db.ExecContext(ctx, query, args...)
		if err != nil && err != sql.ErrNoRows {
			return
		}

		if err == sql.ErrNoRows {
			err = store.ErrorNotFound
		}
	}

	return
}

func (s *CategoryRepository) prepareArgs(data category.Entity) (sets []string, args []any) {
	if data.Name != nil {
		args = append(args, data.Name)
		sets = append(sets, fmt.Sprintf("full_name=$%d", len(args)))
	}

	return
}

func (s *CategoryRepository) Delete(ctx context.Context, id string) (err error) {
	query := `
		DELETE 
		FROM category
		WHERE id=$1`

	args := []any{id}

	_, err = s.db.ExecContext(ctx, query, args...)
	if err != nil && err != sql.ErrNoRows {
		return
	}

	if err == sql.ErrNoRows {
		err = store.ErrorNotFound
	}

	return
}
