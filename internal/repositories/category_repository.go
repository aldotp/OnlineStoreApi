package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/aldotp/OnlineStore/internal/entity"
)

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (c *CategoryRepository) GetCategoryByName(ctx context.Context, name string) (*entity.Category, error) {

	row := c.db.QueryRowContext(ctx, "SELECT id, name, created_at, updated_at FROM categories WHERE name = ?", name)
	var category entity.Category
	if err := row.Scan(&category.ID, &category.Name, &category.CreatedAt, &category.UpdatedAt); err != nil {
		return nil, err
	}

	return &category, nil

}

func (c *CategoryRepository) GetCategoryByID(ctx context.Context, id int) (*entity.Category, error) {

	row := c.db.QueryRowContext(ctx, "SELECT id, name, created_at, updated_at FROM categories WHERE id = ?", id)
	var category entity.Category
	err := row.Scan(&category.ID, &category.Name, &category.CreatedAt, &category.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &category, nil

}

func (c *CategoryRepository) StoreCategory(ctx context.Context, category *entity.Category) (*entity.Category, error) {

	_, err := c.db.ExecContext(ctx, "INSERT INTO categories (name, created_at, updated_at) VALUES (?, ?, ?)", category.Name, time.Now().UTC(), time.Now().UTC())
	if err != nil {
		return nil, err
	}

	row := c.db.QueryRowContext(ctx, "SELECT id, name, created_at, updated_at FROM categories WHERE name = ?", category.Name)

	var insertedCategory entity.Category

	if err := row.Scan(&insertedCategory.ID, &insertedCategory.Name, &insertedCategory.CreatedAt, &insertedCategory.UpdatedAt); err != nil {
		return nil, err
	}

	return &insertedCategory, nil

}

func (c *CategoryRepository) GetAllCategory(ctx context.Context) ([]entity.Category, error) {

	rows, err := c.db.QueryContext(ctx, "SELECT id, name, created_at, updated_at FROM categories")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var categories []entity.Category

	for rows.Next() {
		var category entity.Category
		err := rows.Scan(&category.ID, &category.Name, &category.CreatedAt, &category.UpdatedAt)
		if err != nil {
			return nil, err
		}

		categories = append(categories, category)
	}

	return categories, nil

}

func (c *CategoryRepository) UpdateCategory(ctx context.Context, category *entity.Category) error {

	_, err := c.db.ExecContext(ctx, "UPDATE categories SET name = ?, updated_at = ? WHERE id = ?", category.Name, time.Now().UTC(), category.ID)
	if err != nil {
		return err
	}

	return nil
}

func (c *CategoryRepository) DeleteCategoryByID(ctx context.Context, id int) error {

	_, err := c.db.ExecContext(ctx, "DELETE FROM categories WHERE id = ?", id)
	if err != nil {
		return err
	}

	return nil
}
