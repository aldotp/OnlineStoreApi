package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/aldotp/OnlineStore/internal/entity"
)

type CartItemsRepository struct {
	db *sql.DB
}

func NewCartItemsRepository(db *sql.DB) *CartItemsRepository {
	return &CartItemsRepository{db: db}
}

func (c *CartItemsRepository) StoreCartItems(ctx context.Context, cartItem *entity.CartItem) (*entity.CartItem, error) {

	_, err := c.db.ExecContext(
		ctx,
		"INSERT INTO cart_items (cart_id, product_id, quantity, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
		cartItem.CartID,
		cartItem.ProductID,
		cartItem.Quantity,
		time.Now().UTC(),
		time.Now().UTC(),
	)
	if err != nil {
		return nil, err
	}

	var insertedCartItem entity.CartItem
	row := c.db.QueryRowContext(ctx, "SELECT id, cart_id, product_id, quantity, created_at, updated_at FROM cart_items WHERE cart_id = ? AND product_id = ?", cartItem.CartID, cartItem.ProductID)
	if err := row.Scan(&insertedCartItem.ID, &insertedCartItem.CartID, &insertedCartItem.ProductID, &insertedCartItem.Quantity, &insertedCartItem.CreatedAt, &insertedCartItem.UpdatedAt); err != nil {
		return nil, err
	}

	row = c.db.QueryRowContext(ctx, "SELECT id, name, description, price, category_id, created_at, updated_at FROM products WHERE id = ?", insertedCartItem.ProductID)
	var product entity.Product
	if err := row.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.CategoryID, &product.CreatedAt, &product.UpdatedAt); err != nil {
		return nil, err
	}

	insertedCartItem.Product = product

	return &insertedCartItem, err
}

func (c *CartItemsRepository) GetCartItemsByCartID(ctx context.Context, cartID int) ([]*entity.CartItem, error) {

	rows, err := c.db.QueryContext(ctx, "SELECT id, cart_id, product_id, quantity, created_at, updated_at FROM cart_items WHERE cart_id = ?", cartID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var cartItems []*entity.CartItem

	for rows.Next() {
		var cartItem entity.CartItem
		err := rows.Scan(&cartItem.ID, &cartItem.CartID, &cartItem.ProductID, &cartItem.Quantity, &cartItem.CreatedAt, &cartItem.UpdatedAt)
		if err != nil {
			return nil, err
		}

		row := c.db.QueryRowContext(ctx, "SELECT id, name, description, price, category_id, created_at, updated_at FROM products WHERE id = ?", cartItem.ProductID)
		var product entity.Product

		err = row.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.CategoryID, &product.CreatedAt, &product.UpdatedAt)
		if err != nil {
			return nil,  err
		}

		cartItem.Product = product
		cartItems = append(cartItems, &cartItem)
	}

	return cartItems, nil
}
