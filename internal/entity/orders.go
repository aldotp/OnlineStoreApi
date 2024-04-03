package entity

import "time"

type Order struct {
	ID           int            `json:"id"`
	UserID       int            `json:"user_id"`
	Status       string         `json:"status"`
	TotalAmount  float64        `json:"total_amount"`
	OrderDetails []*OrderDetail `json:"order_details"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}
