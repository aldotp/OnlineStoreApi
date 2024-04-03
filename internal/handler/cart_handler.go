package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/aldotp/OnlineStore/internal/helper"
	"github.com/aldotp/OnlineStore/internal/model"
	"github.com/aldotp/OnlineStore/internal/repositories"
	"github.com/aldotp/OnlineStore/internal/services"
	"github.com/gorilla/mux"
)

type CartHandler struct {
	cartSvc       services.CartService
	repoCart      *repositories.CartRepository
	repoCartItems *repositories.CartItemsRepository
	repoProduct   *repositories.ProductRepository
}

func NewCartHandler(cartSvc services.CartService, repoCart *repositories.CartRepository, repoCartItems *repositories.CartItemsRepository, repoProduct *repositories.ProductRepository) *CartHandler {
	return &CartHandler{
		cartSvc:       cartSvc,
		repoCart:      repoCart,
		repoCartItems: repoCartItems,
		repoProduct:   repoProduct,
	}
}

func (c *CartHandler) AddToCart(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	userCtx, err := helper.GetUserCtx(ctx)
	if err != nil {
		helper.ErrorJSON(helper.Response{
			Code:    http.StatusUnauthorized,
			Message: "unauthorized",
		}, w, http.StatusUnauthorized)
		return
	}

	var request model.CartItemsRequest

	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		helper.ErrorJSON(helper.Response{
			Code:    http.StatusBadRequest,
			Message: "invalid json body",
		}, w, http.StatusBadRequest)
		return
	}

	_, err = c.cartSvc.AddToCart(ctx, request, userCtx.ID)
	if err != nil {
		helper.ErrorJSON(helper.Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}, w, http.StatusInternalServerError)
		return
	}

	helper.WriteJSON(w, http.StatusCreated, helper.Response{
		Code:    http.StatusCreated,
		Message: "Success",
		Status:  "Product added to cart successfully",
	})

}

func (c *CartHandler) Cart(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()

	userCtx, err := helper.GetUserCtx(ctx)
	if err != nil {
		helper.ErrorJSON(helper.Response{
			Code:    http.StatusUnauthorized,
			Message: "unauthorized",
		}, w, http.StatusUnauthorized)
		return
	}

	response, err := c.cartSvc.ViewCart(ctx, userCtx.ID)
	if err != nil {
		helper.ErrorJSON(helper.Response{
			Code:    http.StatusInternalServerError,
			Message: "cannot get cart",
		}, w, http.StatusInternalServerError)
		return
	}

	helper.WriteJSON(w, http.StatusOK, helper.Response{
		Code:    http.StatusOK,
		Status:  "Success",
		Message: "List of products in the shopping cart",
		Data:    response,
	})

}

func (c *CartHandler) DeleteProductFromCart(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()

	userCtx, err := helper.GetUserCtx(ctx)
	if err != nil {
		helper.ErrorJSON(helper.Response{
			Code:    http.StatusUnauthorized,
			Message: "unauthorized",
		}, w, http.StatusUnauthorized)
		return
	}

	pathParam := mux.Vars(r)["id"]
	productID, err := strconv.Atoi(pathParam)
	if err != nil {
		helper.ErrorJSON(helper.Response{
			Code:    http.StatusBadRequest,
			Message: "invalid product id",
		}, w, http.StatusBadRequest)
		return
	}

	request := model.DeleteProductRequest{
		ProductID: productID,
	}

	err = c.cartSvc.RemoveFromCart(ctx, request, userCtx.ID)
	if err != nil {
		helper.ErrorJSON(helper.Response{
			Code:    http.StatusInternalServerError,
			Message: "cannot delete product from cart",
		}, w, http.StatusInternalServerError)
		return
	}

	helper.WriteJSON(w, http.StatusOK, helper.Response{
		Code:    http.StatusOK,
		Message: "Delete Product from Cart Success",
	})

}

func (c *CartHandler) EmptyCart(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()

	userCtx, err := helper.GetUserCtx(ctx)
	if err != nil {
		helper.ErrorJSON(helper.Response{
			Code:    http.StatusUnauthorized,
			Message: "unauthorized",
		}, w, http.StatusUnauthorized)
		return
	}

	err = c.cartSvc.EmptyCart(ctx, userCtx.ID)
	if err != nil {
		helper.ErrorJSON(helper.Response{
			Code:    http.StatusInternalServerError,
			Message: "cannot empty cart",
		}, w, http.StatusInternalServerError)
		return
	}

	helper.WriteJSON(w, http.StatusOK, helper.Response{
		Code:    http.StatusOK,
		Message: "Empty Cart Success",
	})
}

func (c *CartHandler) ModifyCart(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()

	userCtx, err := helper.GetUserCtx(ctx)
	if err != nil {
		helper.ErrorJSON(helper.Response{
			Code:    http.StatusUnauthorized,
			Message: "unauthorized",
		}, w, http.StatusUnauthorized)
		return
	}

	pathParam := mux.Vars(r)["id"]
	productID, err := strconv.Atoi(pathParam)
	if err != nil {
		helper.ErrorJSON(helper.Response{
			Code:    http.StatusBadRequest,
			Message: "invalid product id",
		}, w, http.StatusBadRequest)
		return
	}

	var request model.ModifyCartRequest

	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		helper.ErrorJSON(helper.Response{
			Code:    http.StatusBadRequest,
			Message: "invalid json body",
		}, w, http.StatusBadRequest)
		return
	}

	request.ProductID = productID

	err = c.cartSvc.ModifyCart(ctx, request, userCtx.ID)
	if err != nil {
		helper.ErrorJSON(helper.Response{
			Code:    http.StatusInternalServerError,
			Message: "cannot modify cart",
		}, w, http.StatusInternalServerError)
		return
	}

	helper.WriteJSON(w, http.StatusOK, helper.Response{
		Code:    http.StatusOK,
		Message: "Modify Cart Success",
	})
}
