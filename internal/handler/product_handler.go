package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/aldotp/OnlineStore/internal/helper"
	"github.com/aldotp/OnlineStore/internal/model"
	"github.com/aldotp/OnlineStore/internal/services"
	"github.com/gorilla/mux"
)

type ProductHandler struct {
	ProductSvc services.ProductService
}

func NewProductHandler(ProductSvc services.ProductService) *ProductHandler {
	return &ProductHandler{
		ProductSvc: ProductSvc,
	}
}

func (p *ProductHandler) GetProductsByCategory(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	pathParam := mux.Vars(r)["category"]
	response, err := p.ProductSvc.GetProductByCategoryID(ctx, pathParam)
	if err != nil {
		helper.ErrorJSON(helper.Response{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}, w, http.StatusBadRequest)
		return
	}

	helper.WriteJSON(w, http.StatusOK, helper.Response{
		Code:    http.StatusOK,
		Message: "Success",
		Data:    response,
	})

}

func (p *ProductHandler) StoreProducts(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	_, err := helper.GetUserCtx(ctx)
	if err != nil {
		helper.ErrorJSON(helper.Response{
			Code:    http.StatusUnauthorized,
			Message: "unauthorized",
		}, w, http.StatusUnauthorized)
		return
	}

	var request model.ProductRequest
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		helper.ErrorJSON(helper.Response{
			Code:    http.StatusBadRequest,
			Message: "invalid json body",
		}, w, http.StatusBadRequest)
		return
	}

	response, err := p.ProductSvc.StoreProduct(ctx, request)
	if err != nil {
		helper.ErrorJSON(helper.Response{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}, w, http.StatusBadRequest)
		return
	}

	helper.WriteJSON(w, http.StatusCreated, helper.Response{
		Code:    http.StatusCreated,
		Message: "Success",
		Data:    response,
	})

}
