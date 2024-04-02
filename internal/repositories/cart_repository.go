package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/aldotp/OnlineStore/internal/entity"
)

type CartRepository struct {
	db *sql.DB
}

func NewCartRepository(db *sql.DB) *CartRepository {
	return &CartRepository{
		db: db,
	}
}

func (c *CartRepository) Store(ctx context.Context, cart *entity.Cart) error {
	_, err := c.db.ExecContext(
		ctx,
		"INSERT INTO carts (user_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
		cart.UserID,
		time.Now().UTC(),
		time.Now().UTC(),
	)
	if err != nil {
		return err
	}

	return err
}

func (c *CartRepository) GetCartByUserID(ctx context.Context, userID int) (*entity.Cart, error) {

	row := c.db.QueryRowContext(ctx, "SELECT id, user_id, created_at, updated_at FROM carts WHERE user_id = ?", userID)
	var cart entity.Cart
	if err := row.Scan(&cart.ID, &cart.UserID, &cart.CreatedAt, &cart.UpdatedAt); err != nil {
		return nil, err
	}

	return &cart, nil
}

func (c *CartRepository) DeleteProductFromCart(ctx context.Context, cartID int, productID int) error {

	_, err := c.db.ExecContext(ctx, "DELETE FROM cart_items WHERE cart_id = ? AND product_id = ?", cartID, productID)
	if err != nil {
		return err
	}

	return err
}

func (r *CartRepository) GetCartItemsByUserID(ctx context.Context, userID int) ([]*entity.CartItem, error) {

	query := `
		SELECT id, product_id, quantity, created_at, updated_at FROM cart_items
		WHERE cart_id IN (SELECT id FROM carts WHERE user_id = ?)
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cartItems []*entity.CartItem
	for rows.Next() {
		var cartItem entity.CartItem
		if err := rows.Scan(&cartItem.ID, &cartItem.ProductID, &cartItem.Quantity, &cartItem.CreatedAt, &cartItem.UpdatedAt); err != nil {
			return nil, err
		}

		row := r.db.QueryRowContext(ctx, "SELECT id, name, description, price, category_id, created_at, updated_at FROM products WHERE id = ?", cartItem.ProductID)

		var product entity.Product

		if err := row.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.CategoryID, &product.CreatedAt, &product.UpdatedAt); err != nil {
			return nil, err
		}

		cartItem.Product = product

		cartItems = append(cartItems, &cartItem)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return cartItems, nil
}

// ClearCart clears the user's cart in the database.
func (r *CartRepository) ClearCart(ctx context.Context, userID int) error {
	// Execute a DELETE query to clear the user's cart
	_, err := r.db.ExecContext(ctx, "DELETE FROM cart_items WHERE cart_id IN (SELECT id FROM carts WHERE user_id = ?)", userID)
	if err != nil {
		return err
	}

	return nil
}
func (r *CartRepository) ClearCartWithTransaction(ctx context.Context, tx *sql.Tx, userID int) error {
	// Execute a DELETE query to clear the user's cart
	_, err := tx.ExecContext(ctx, "DELETE FROM cart_items WHERE cart_id IN (SELECT id FROM carts WHERE user_id = ?)", userID)
	return err
}

// GetCartItemByUserIDAndProductID retrieves a cart item by user ID and product ID.
func (r *CartRepository) GetCartItemByUserIDAndProductID(ctx context.Context, userID, productID int) (*entity.CartItem, error) {
	query := "SELECT id, cart_id, product_id, quantity, created_at, updated_at FROM cart_items WHERE cart_id IN (SELECT id FROM carts WHERE user_id = ?) AND product_id = ?"
	row := r.db.QueryRowContext(ctx, query, userID, productID)

	var cartItem entity.CartItem
	if err := row.Scan(&cartItem.ID, &cartItem.CartID, &cartItem.ProductID, &cartItem.Quantity, &cartItem.CreatedAt, &cartItem.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Cart item not found
		}
		return nil, err
	}

	return &cartItem, nil
}

// UpdateCartItem updates a cart item in the repository.
func (r *CartRepository) UpdateCartItem(ctx context.Context, item *entity.CartItem) error {
	query := "UPDATE cart_items SET quantity = ? WHERE id = ?"
	_, err := r.db.ExecContext(ctx, query, item.Quantity, item.ID)
	if err != nil {
		return err
	}
	return nil
}
