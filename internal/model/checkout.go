package model

type CheckoutHistoryResponse struct {
	ID           int            `json:"id"`
	UserID       int            `json:"user_id"`
	Status       string         `json:"status"`
	TotalPrice   float64        `json:"total_price"`
	TotalProduct int            `json:"total_product"`
	CreatedAt    string         `json:"created_at"`
	UpdatedAt    string         `json:"updated_at"`
	OrderDetails []*OrderDetail `json:"order_details"`
}

type OrderDetail struct {
	ID          int     `json:"id"`
	ProductID   int     `json:"product_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
}
