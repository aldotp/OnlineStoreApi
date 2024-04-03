package model

import "github.com/aldotp/OnlineStore/internal/entity"

type CartItemsRequest struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

type ViewCartResponse struct {
	CartItems    []*entity.CartItem `json:"cart_items"`
	Total        int                `json:"total"`
	TotalProduct int                `json:"total_product"`
	TotalPrice   float64            `json:"total_price"`
}

type RemoveCartItemRequest struct {
	ProductID int `json:"product_id"`
}

type CheckoutResponse struct {
	CartItems    []*entity.CartItem `json:"cart_items"`
	Total        int                `json:"total"`
	TotalProduct int                `json:"total_product"`
	TotalPrice   float64            `json:"total_price"`
}

type ModifyCartRequest struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}
