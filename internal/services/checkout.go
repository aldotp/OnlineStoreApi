package services

import (
	"context"
	"fmt"

	"github.com/aldotp/OnlineStore/internal/entity"
	"github.com/aldotp/OnlineStore/internal/model"
	"github.com/aldotp/OnlineStore/internal/repositories"
)

type CheckoutService interface {
	Checkout(ctx context.Context, userID int) (*model.CheckoutResponse, error)
}

type checkout struct {
	orderRepo       *repositories.OrderRepository
	cartRepo        *repositories.CartRepository
	orderDetailRepo *repositories.OrderDetailRepository
	paymentSvc      PaymentService
}

func NewCheckout(orderRepo *repositories.OrderRepository, cartRepo *repositories.CartRepository, orderDetailRepo *repositories.OrderDetailRepository, paymentSvc PaymentService) CheckoutService {
	return &checkout{
		orderRepo:       orderRepo,
		cartRepo:        cartRepo,
		orderDetailRepo: orderDetailRepo,
		paymentSvc:      paymentSvc,
	}
}

func (c *checkout) Checkout(ctx context.Context, userID int) (*model.CheckoutResponse, error) {

	// start transaction
	tx, err := c.orderRepo.BeginTransaction(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	cartItems, err := c.cartRepo.GetCartItemsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("cannot get cart items")
	}

	var totalAmount float64
	var count int = 0
	for _, item := range cartItems {
		totalAmount += float64(item.Quantity) * item.Product.Price
		count += item.Quantity
	}

	order := &entity.Order{
		UserID:      userID,
		TotalAmount: totalAmount,
	}

	createdOrder, err := c.orderRepo.CreateOrderWithTransaction(ctx, tx, order)
	if err != nil {
		return nil, fmt.Errorf("cannot create order")
	}

	for _, item := range cartItems {
		orderDetail := &entity.OrderDetail{
			OrderID:   createdOrder.ID,
			ProductID: item.Product.ID,
			Quantity:  item.Quantity,
			Price:     item.Product.Price,
		}

		err := c.orderDetailRepo.CreateOrderDetailWithTransaction(ctx, tx, orderDetail)
		if err != nil {
			return nil, fmt.Errorf("cannot create order detail")
		}
	}

	paymentStatus := c.paymentSvc.ProcessPayment(ctx, totalAmount)
	if paymentStatus != "success" {
		return nil, fmt.Errorf("payment failed")
	}

	err = c.orderRepo.UpdateOrderStatusWithTransaction(ctx, tx, createdOrder.ID, "paid")
	if err != nil {
		return nil, fmt.Errorf("cannot update order status")
	}

	err = c.cartRepo.ClearCartWithTransaction(ctx, tx, userID)
	if err != nil {
		return nil, fmt.Errorf("cannot clear cart")
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("transaction commit failed")
	}

	return &model.CheckoutResponse{
		CartItems:    cartItems,
		Total:        count,
		TotalProduct: count,
		TotalPrice:   totalAmount,
	}, nil

}
