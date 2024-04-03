package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/aldotp/OnlineStore/internal/helper"
	"github.com/aldotp/OnlineStore/internal/services"
)

type CheckoutHandler struct {
	checkoutSvc services.CheckoutService
}

func NewCheckoutHandler(checkoutSvc services.CheckoutService) *CheckoutHandler {
	return &CheckoutHandler{
		checkoutSvc: checkoutSvc,
	}
}

func (h *CheckoutHandler) CheckoutHandler(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	userCtx, err := helper.GetUserCtx(ctx)
	if err != nil {
		helper.ErrorJSON(helper.Response{
			Code:    http.StatusUnauthorized,
			Message: "unauthorized",
		}, w, http.StatusUnauthorized)
		return
	}

	response, err := h.checkoutSvc.Checkout(ctx, userCtx.ID)
	if err != nil {
		helper.ErrorJSON(helper.Response{
			Code:    http.StatusInternalServerError,
			Message: "cannot checkout",
		}, w, http.StatusInternalServerError)
		return
	}

	helper.WriteJSON(w, http.StatusOK, helper.Response{
		Code:    http.StatusOK,
		Status:  "Success",
		Message: "Checkout Successful",
		Data:    response,
	})
}

func (h *CheckoutHandler) CheckoutHistory(w http.ResponseWriter, r *http.Request) {

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

	response, err := h.checkoutSvc.History(ctx, userCtx.ID)
	if err != nil {
		helper.ErrorJSON(helper.Response{
			Code:    http.StatusInternalServerError,
			Message: "cannot get checkout history",
		}, w, http.StatusInternalServerError)
		return
	}

	helper.WriteJSON(w, http.StatusOK, helper.Response{
		Code:    http.StatusOK,
		Status:  "Success",
		Message: "Checkout History",
		Data:    response,
	})
}
