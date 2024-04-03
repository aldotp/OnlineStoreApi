package services

import (
	"context"
	"fmt"

	"github.com/aldotp/OnlineStore/internal/entity"
	"github.com/aldotp/OnlineStore/internal/model"
	"github.com/aldotp/OnlineStore/internal/repositories"
)

type CartService interface {
	AddToCart(ctx context.Context, request model.CartItemsRequest, userID int) (*entity.CartItem, error)
	ViewCart(ctx context.Context, userID int) (*model.ViewCartResponse, error)
	RemoveFromCart(ctx context.Context, request model.DeleteProductRequest, userID int) error
	EmptyCart(ctx context.Context, userID int) error
	ModifyCart(ctx context.Context, request model.ModifyCartRequest, userID int) error
}

type cart struct {
	repo          *repositories.CartRepository
	repoCartItems *repositories.CartItemsRepository
	repoProduct   *repositories.ProductRepository
}

func NewCart(repo *repositories.CartRepository, repoCartItems *repositories.CartItemsRepository, repoProduct *repositories.ProductRepository) CartService {
	return &cart{
		repo:          repo,
		repoCartItems: repoCartItems,
		repoProduct:   repoProduct,
	}
}

func (c *cart) AddToCart(ctx context.Context, request model.CartItemsRequest, userID int) (*entity.CartItem, error) {

	cart, err := c.repo.GetCartByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	product, err := c.repoProduct.GetProductByID(ctx, request.ProductID)
	if err != nil {
		return nil, err
	}

	if product == nil {
		return nil, fmt.Errorf("product not found")
	}

	item, err := c.repo.GetCartItemByUserIDAndProductID(ctx, userID, request.ProductID)
	if err != nil {
		return nil, err
	}

	if item != nil {
		item.Quantity += request.Quantity
		err := c.repo.UpdateCartItem(ctx, item)
		if err != nil {
			return nil, err
		}

		return item, nil
	}

	data := entity.CartItem{
		CartID:    cart.ID,
		ProductID: request.ProductID,
		Quantity:  request.Quantity,
	}

	cartItem, err := c.repoCartItems.StoreCartItems(ctx, &data)
	if err != nil {
		return nil, err
	}

	return cartItem, nil
}

func (c *cart) RemoveFromCart(ctx context.Context, request model.DeleteProductRequest, userID int) error {
	cart, err := c.repo.GetCartByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("cannot get cart")
	}

	err = c.repo.DeleteProductFromCart(ctx, cart.ID, request.ProductID)
	if err != nil {
		return fmt.Errorf("cannot delete product from cart")
	}

	return nil

}

func (c *cart) ViewCart(ctx context.Context, userID int) (*model.ViewCartResponse, error) {

	cart, err := c.repo.GetCartByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("cannot get cart")
	}

	cartItems, err := c.repoCartItems.GetCartItemsByCartID(ctx, cart.ID)
	if err != nil {
		return nil, fmt.Errorf("cannot get cart items")
	}

	var total float64 = 0
	count := 0

	for _, item := range cartItems {
		price := item.Product.Price * float64(item.Quantity)
		total += price
		count += item.Quantity
	}

	return &model.ViewCartResponse{
		CartItems:    cartItems,
		TotalProduct: len(cartItems),
		Total:        count,
		TotalPrice:   total,
	}, nil
}

func (c *cart) EmptyCart(ctx context.Context, userID int) error {

	cart, err := c.repo.GetCartByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("cannot get cart")
	}

	err = c.repo.EmptyCart(ctx, cart.ID)
	if err != nil {
		return fmt.Errorf("cannot empty cart")
	}

	return nil

}

func (c *cart) ModifyCart(ctx context.Context, request model.ModifyCartRequest, userID int) error {

	cart, err := c.repo.GetCartByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("cannot get cart")
	}

	err = c.repo.ModifyCart(ctx, cart.ID, request.ProductID, request.Quantity)
	if err != nil {
		return fmt.Errorf("cannot modify cart")
	}

	return nil
}
