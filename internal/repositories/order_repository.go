package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/aldotp/OnlineStore/internal/entity"
)

// orderRepository implements the OrderRepository interface.
type OrderRepository struct {
	db *sql.DB
}

// NewOrderRepository creates a new instance of orderRepository.
func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{
		db: db,
	}
}

func (r *OrderRepository) CreateOrder(ctx context.Context, order *entity.Order) (*entity.Order, error) {

	tNow := time.Now().UTC()
	result, err := r.db.ExecContext(ctx, "INSERT INTO orders (user_id, total_amount, created_at, updated_at) VALUES (?, ?, ?, ?)", order.UserID, order.TotalAmount, tNow, tNow)
	if err != nil {
		return nil, err
	}

	orderID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	order.ID = int(orderID)

	return order, nil
}

// CreateOrderWithTransaction membuat order dalam transaksi yang diberikan.
func (repo *OrderRepository) CreateOrderWithTransaction(ctx context.Context, tx *sql.Tx, order *entity.Order) (*entity.Order, error) {

	createdAt := time.Now().UTC()
	updatedAt := time.Now().UTC()
	query := "INSERT INTO orders (user_id, total_amount, created_at, updated_at) VALUES (?, ?, ?, ?)"
	result, err := tx.ExecContext(ctx, query, order.UserID, order.TotalAmount, createdAt, updatedAt)
	if err != nil {
		return nil, err
	}

	orderID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	order.ID = int(orderID)
	order.CreatedAt = createdAt
	order.UpdatedAt = updatedAt

	return order, nil
}
func (repo *OrderRepository) BeginTransaction(ctx context.Context) (*sql.Tx, error) {
	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (r *OrderRepository) UpdateOrderStatusWithTransaction(ctx context.Context, tx *sql.Tx, orderID int, status string) error {
	_, err := tx.ExecContext(ctx, "UPDATE orders SET status = ? WHERE id = ?", status, orderID)
	return err

}
