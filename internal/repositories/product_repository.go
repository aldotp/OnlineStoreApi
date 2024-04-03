package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/aldotp/OnlineStore/internal/entity"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

func (u *ProductRepository) GetProductsByCategoryID(ctx context.Context, ctg *entity.Category) (*entity.Category, error) {

	var categories entity.Category

	rows, err := u.db.QueryContext(ctx, "SELECT id, name, description, price, category_id, created_at, updated_at FROM products WHERE category_id = ?", ctg.ID)
	if err != nil {
		return nil, err
	}

	var products []entity.Product

	for rows.Next() {
		var product entity.Product
		err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.CategoryID, &product.CreatedAt, &product.UpdatedAt)
		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	categories.Products = products

	return &categories, nil
}

func (u *ProductRepository) GetAllProducts(ctx context.Context) ([]entity.Product, error) {

	rows, err := u.db.QueryContext(ctx, "SELECT id, name, description, price, category_id, created_at, updated_at FROM products")
	if err != nil {
		return nil, err
	}

	var products []entity.Product

	for rows.Next() {
		var product entity.Product
		err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.CategoryID, &product.CreatedAt, &product.UpdatedAt)
		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	return products, nil
}

func (u *ProductRepository) GetProductByID(ctx context.Context, id int) (*entity.Product, error) {

	row := u.db.QueryRowContext(ctx, "SELECT id, name, description, price, category_id, created_at, updated_at FROM products WHERE id = ?", id)
	var product entity.Product
	if err := row.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.CategoryID, &product.CreatedAt, &product.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &product, nil

}

func (u *ProductRepository) StoreProduct(ctx context.Context, product entity.Product) (*entity.Product, error) {

	tNow := time.Now().UTC()
	result, err := u.db.ExecContext(ctx, "INSERT INTO products (name, description, price, category_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)", product.Name, product.Description, product.Price, product.CategoryID, tNow, tNow)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	var insertedProduct entity.Product
	err = u.db.QueryRowContext(ctx, "SELECT id, name, description, price, category_id, created_at, updated_at FROM products WHERE id = ?", id).
		Scan(&insertedProduct.ID, &insertedProduct.Name, &insertedProduct.Description, &insertedProduct.Price, &insertedProduct.CategoryID, &insertedProduct.CreatedAt, &insertedProduct.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &insertedProduct, nil

}

func (u *ProductRepository) UpdateProduct(ctx context.Context, product *entity.Product) error {

	_, err := u.db.ExecContext(ctx, "UPDATE products SET name = ?, description = ?, price = ?, updated_at = ? WHERE id = ?", product.Name, product.Description, product.Price, time.Now().UTC(), product.ID)
	if err != nil {
		return nil
	}

	return nil

}

func (u *ProductRepository) DeleteProduct(ctx context.Context, id int) error {

	_, err := u.db.ExecContext(ctx, "DELETE FROM products WHERE id = ?", id)
	if err != nil {
		return err
	}

	return nil

}
