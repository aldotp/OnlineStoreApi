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

func (repo *OrderDetailRepository) GetOrderDetailsByOrderID(ctx context.Context, orderID int) ([]*entity.OrderDetail, error) {

	rows, err := repo.db.QueryContext(ctx, "SELECT id, order_id, product_id, quantity, price, created_at, updated_at FROM order_details WHERE order_id = ?", orderID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var orderDetails []*entity.OrderDetail

	for rows.Next() {
		var orderDetail entity.OrderDetail
		if err := rows.Scan(&orderDetail.ID, &orderDetail.OrderID, &orderDetail.ProductID, &orderDetail.Quantity, &orderDetail.Price, &orderDetail.CreatedAt, &orderDetail.UpdatedAt); err != nil {
			return nil, err
		}

		var product entity.Product
		row := repo.db.QueryRowContext(ctx, "SELECT id, name, description, price, created_at, updated_at FROM products WHERE id = ?", orderDetail.ProductID)
		if err := row.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.CreatedAt, &product.UpdatedAt); err != nil {
			return nil, err
		}

		orderDetail.Product = &product
		orderDetails = append(orderDetails, &orderDetail)
	}

	return orderDetails, nil

}
