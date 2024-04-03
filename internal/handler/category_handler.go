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
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
)

type CategoryHandler struct {
	repo        *repositories.CategoryRepository
	redis       *redis.Client
	categorySvc services.CategoryService
}

func NewCategoryHandler(categorySvc services.CategoryService, repo *repositories.CategoryRepository, redis *redis.Client) *CategoryHandler {
	return &CategoryHandler{
		repo:        repo,
		redis:       redis,
		categorySvc: categorySvc,
	}
}

func (p *CategoryHandler) StoreCategory(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel()

	_, err := helper.GetUserCtx(ctx)
	if err != nil {
		helper.ErrorJSON(helper.Response{
			Code:    http.StatusUnauthorized,
			Message: "unauthorized",
		}, w, http.StatusUnauthorized)
		return
	}

	var request model.CategoryRequest
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		helper.ErrorJSON(helper.Response{
			Code:    http.StatusBadRequest,
			Message: "invalid json body",
		}, w, http.StatusBadRequest)
		return
	}

	response, err := p.categorySvc.StoreCategory(ctx, request)
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

func (p *CategoryHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	_, err := helper.GetUserCtx(ctx)
	if err != nil {
		helper.ErrorJSON(helper.Response{
			Code:    http.StatusUnauthorized,
			Message: "unauthorized",
		}, w, http.StatusUnauthorized)
		return
	}

	category, err := p.categorySvc.GetCategories(ctx)
	if err != nil {
		helper.ErrorJSON(helper.Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}, w, http.StatusInternalServerError)
		return
	}

	helper.WriteJSON(w, http.StatusOK, helper.Response{
		Code:    http.StatusOK,
		Message: "Success",
		Data:    category,
	})
}

func (p *CategoryHandler) GetCategoryByID(w http.ResponseWriter, r *http.Request) {

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

	category, err := p.categorySvc.GetCategoryByID(ctx, id)
	if err != nil {
		helper.ErrorJSON(helper.Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}, w, http.StatusInternalServerError)
		return
	}

	if category == nil {
		helper.ErrorJSON(helper.Response{
			Code:    http.StatusNotFound,
			Message: "category not found",
		}, w, http.StatusNotFound)
		return
	}

	helper.WriteJSON(w, http.StatusOK, helper.Response{
		Code:    http.StatusOK,
		Message: "Success",
		Data:    category,
	})

}

func (p *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
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

	err = p.categorySvc.DeleteCategoryByID(ctx, id)
	if err != nil {
		helper.ErrorJSON(helper.Response{
			Code:    http.StatusInternalServerError,
			Message: "cannot delete category",
		}, w, http.StatusInternalServerError)
		return
	}

	helper.WriteJSON(w, http.StatusOK, helper.Response{
		Code:    http.StatusOK,
		Message: "Success",
	})

}

func (p *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
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

	var request model.UpdateCategoryRequest

	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		helper.ErrorJSON(helper.Response{
			Code:    http.StatusBadRequest,
			Message: "invalid json body",
		}, w, http.StatusBadRequest)
		return
	}

	request.CategoryID = id

	err = p.categorySvc.UpdateCategory(ctx, request)
	if err != nil {
		helper.ErrorJSON(helper.Response{
			Code:    http.StatusInternalServerError,
			Message: "cannot update category",
		}, w, http.StatusInternalServerError)
		return
	}

	helper.WriteJSON(w, http.StatusOK, helper.Response{
		Code:    http.StatusOK,
		Message: "Success Update Category",
	})

}
