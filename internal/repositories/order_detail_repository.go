package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/aldotp/OnlineStore/internal/entity"
)

// orderDetailRepository implements the OrderDetailRepository interface.
type OrderDetailRepository struct {
	db *sql.DB
}

// NewOrderDetailRepository creates a new instance of orderDetailRepository.
func NewOrderDetailRepository(db *sql.DB) *OrderDetailRepository {
	return &OrderDetailRepository{
		db: db,
	}
}

func (repo *OrderDetailRepository) CreateOrderDetailWithTransaction(ctx context.Context, tx *sql.Tx, orderDetail *entity.OrderDetail) error {
	query := "INSERT INTO order_details (order_id, product_id, quantity, price, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)"
	_, err := tx.ExecContext(ctx, query, orderDetail.OrderID, orderDetail.ProductID, orderDetail.Quantity, orderDetail.Price, time.Now(), time.Now())
	return err
}
