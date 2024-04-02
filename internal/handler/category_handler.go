package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/aldotp/OnlineStore/internal/entity"
	"github.com/aldotp/OnlineStore/internal/helper"
	"github.com/aldotp/OnlineStore/internal/model"
	"github.com/aldotp/OnlineStore/internal/repositories"
	"github.com/go-redis/redis/v8"
)

type CategoryHandler struct {
	repo  *repositories.CategoryRepository
	redis *redis.Client
}

func NewCategoryHandler(repo *repositories.CategoryRepository, redis *redis.Client) *CategoryHandler {
	return &CategoryHandler{
		repo:  repo,
		redis: redis,
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

	request := new(model.CategoryRequest)
	err = json.NewDecoder(r.Body).Decode(request)
	if err != nil {
		helper.ErrorJSON(helper.Response{
			Code:    http.StatusBadRequest,
			Message: "invalid json body",
		}, w, http.StatusBadRequest)
		return
	}

	category, err := p.repo.StoreCategory(ctx, &entity.Category{
		Name: request.Name,
	})
	if err != nil {
		helper.ErrorJSON(helper.Response{
			Code:    http.StatusInternalServerError,
			Message: "cannot store category",
		}, w, http.StatusInternalServerError)
		return
	}

	// Store category in Redis cache
	categoryJSON, _ := json.Marshal(category)

	key := fmt.Sprintf("category:%d", category.ID)
	err = p.redis.Set(ctx, key, categoryJSON, 24*time.Hour).Err()
	if err != nil {
		helper.ErrorJSON(helper.Response{
			Code:    http.StatusInternalServerError,
			Message: "cannot store category in redis",
		}, w, http.StatusInternalServerError)
		return
	}

	p.redis.Del(ctx, "categories")

	helper.WriteJSON(w, http.StatusCreated, helper.Response{
		Code:    http.StatusCreated,
		Message: "Success",
		Data:    category,
	})

}

func (p *CategoryHandler) GetCategory(w http.ResponseWriter, r *http.Request) {
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

	cachedCategories, err := p.redis.Get(ctx, "categories").Result()
	if err == redis.Nil {
		category, err := p.repo.GetAllCategory(ctx)
		if err != nil {
			helper.ErrorJSON(helper.Response{
				Code:    http.StatusInternalServerError,
				Message: "cannot get all category",
			}, w, http.StatusInternalServerError)
			return
		}

		categoryJSON, _ := json.Marshal(category)
		err = p.redis.Set(ctx, "categories", categoryJSON, 24*time.Hour).Err()
		if err != nil {
			helper.ErrorJSON(helper.Response{
				Code:    http.StatusInternalServerError,
				Message: "cannot store categories in redis",
			}, w, http.StatusInternalServerError)
			return
		}

		helper.WriteJSON(w, http.StatusOK, helper.Response{
			Code:    http.StatusOK,
			Message: "Success",
			Data:    category,
		})
		return

	} else if err != nil {
		helper.ErrorJSON(helper.Response{
			Code:    http.StatusInternalServerError,
			Message: "cannot get categories from redis",
		}, w, http.StatusInternalServerError)
		return
	}

	var category []entity.Category
	err = json.Unmarshal([]byte(cachedCategories), &category)
	if err != nil {
		helper.ErrorJSON(helper.Response{
			Code:    http.StatusInternalServerError,
			Message: "cannot unmarshal categories",
		}, w, http.StatusInternalServerError)
		return
	}

	helper.WriteJSON(w, http.StatusOK, helper.Response{
		Code:    http.StatusOK,
		Message: "Success",
		Data:    category,
	})
}
