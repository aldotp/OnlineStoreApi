package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/aldotp/OnlineStore/internal/helper"
	"github.com/aldotp/OnlineStore/internal/model"
	"github.com/aldotp/OnlineStore/internal/services"
	"github.com/gorilla/mux"
)

type ProductHandler struct {
	productSvc services.ProductService
}

func NewProductHandler(productSvc services.ProductService) *ProductHandler {
	return &ProductHandler{
		productSvc: productSvc,
	}
}

func (p *ProductHandler) GetProductsByCategory(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	pathParam := mux.Vars(r)["id"]
	categoryID, err := strconv.Atoi(pathParam)
	if err != nil {
		helper.ErrorJSON(helper.Response{
			Code:    http.StatusBadRequest,
			Message: "invalid category id",
		}, w, http.StatusBadRequest)
		return
	}

	response, err := p.productSvc.GetProductByCategoryID(ctx, categoryID)
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

	response, err := p.productSvc.StoreProduct(ctx, request)
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

func (p *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {

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

	paramID := mux.Vars(r)["id"]

	id, err := strconv.Atoi(paramID)
	if err != nil {
		helper.ErrorJSON(helper.Response{
			Code:    http.StatusBadRequest,
			Message: "invalid id",
		}, w, http.StatusBadRequest)
		return
	}

	var request model.UpdateProductRequest

	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		helper.ErrorJSON(helper.Response{
			Code:    http.StatusBadRequest,
			Message: "invalid json body",
		}, w, http.StatusBadRequest)
		return
	}

	request.ProductID = id

	err = p.productSvc.UpdateProduct(ctx, request)
	if err != nil {
		helper.ErrorJSON(helper.Response{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}, w, http.StatusBadRequest)
		return
	}

	helper.WriteJSON(w, http.StatusOK, helper.Response{
		Code:    http.StatusOK,
		Message: "Success Update Product",
	})
}

func (p *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {

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

	paramID := mux.Vars(r)["id"]

	id, err := strconv.Atoi(paramID)
	if err != nil {
		helper.ErrorJSON(helper.Response{
			Code:    http.StatusBadRequest,
			Message: "invalid id",
		}, w, http.StatusBadRequest)
		return
	}

	var request model.DeleteProductRequest
	request.ProductID = id

	err = p.productSvc.DeleteProduct(ctx, request)
	if err != nil {
		helper.ErrorJSON(helper.Response{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}, w, http.StatusBadRequest)
		return
	}

	helper.WriteJSON(w, http.StatusOK, helper.Response{
		Code:    http.StatusOK,
		Message: "Success Delete Product",
	})
}

func (p *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {

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

	response, err := p.productSvc.GetProducts(ctx)
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

func (p *ProductHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()

	_, err := helper.GetUserCtx(ctx)
	if err != nil {
		helper.ErrorJSON(helper.Response{
			Code:    http.StatusUnauthorized,
			Message: "unauthorized",
		}, w, http.StatusUnauthorized)
		return
	}

	paramID := mux.Vars(r)["id"]

	id, err := strconv.Atoi(paramID)
	if err != nil {
		helper.ErrorJSON(helper.Response{
			Code:    http.StatusBadRequest,
			Message: "invalid id",
		}, w, http.StatusBadRequest)
		return
	}

	response, err := p.productSvc.GetProductByID(ctx, id)
	if err != nil {
		helper.ErrorJSON(helper.Response{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}, w, http.StatusBadRequest)
		return
	}

	helper.WriteJSON(w, http.StatusOK, helper.Response{
		Code:    http.StatusOK,
		Message: "Success Get Product",
		Data:    response,
	})
}
