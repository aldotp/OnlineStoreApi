package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aldotp/OnlineStore/internal/entity"
	"github.com/aldotp/OnlineStore/internal/model"
	"github.com/aldotp/OnlineStore/internal/repositories"
	"github.com/go-redis/redis/v8"
)

type CategoryService interface {
	GetCategories(ctx context.Context) ([]entity.Category, error)
	StoreCategory(ctx context.Context, request model.CategoryRequest) (*entity.Category, error)
	GetCategoryByID(ctx context.Context, id int) (*entity.Category, error)
	DeleteCategoryByID(ctx context.Context, id int) error
	UpdateCategory(ctx context.Context, request model.UpdateCategoryRequest) error
}

type category struct {
	redis *redis.Client
	repo  *repositories.CategoryRepository
}

func NewCategory(redis *redis.Client, repo *repositories.CategoryRepository) CategoryService {
	return &category{
		redis: redis,
		repo:  repo,
	}
}
func (c *category) GetCategories(ctx context.Context) ([]entity.Category, error) {
	cachedCategories, err := c.redis.Get(ctx, "categories").Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	if cachedCategories == "" {
		categories, err := c.repo.GetAllCategory(ctx)
		if err != nil {
			return nil, err
		}

		categoryJSON, err := json.Marshal(categories)
		if err != nil {
			return nil, err
		}

		err = c.redis.Set(ctx, "categories", categoryJSON, 2*time.Hour).Err()
		if err != nil {
			return nil, err
		}

		return categories, nil
	}

	var categories []entity.Category
	err = json.Unmarshal([]byte(cachedCategories), &categories)
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (c *category) StoreCategory(ctx context.Context, request model.CategoryRequest) (*entity.Category, error) {

	category, err := c.repo.StoreCategory(ctx, &entity.Category{
		Name: request.Name,
	})
	if err != nil {
		return nil, err
	}

	// Store category in Redis cache
	categoryJSON, err := json.Marshal(category)
	if err != nil {
		return nil, err
	}

	err = c.redis.Set(ctx, fmt.Sprintf("category:%d", category.ID), categoryJSON, 2*time.Hour).Err()
	if err != nil {
		return nil, err
	}

	c.redis.Del(ctx, "categories")
	c.redis.Del(ctx, "product/category")
	return category, nil

}

func (c *category) GetCategoryByID(ctx context.Context, id int) (*entity.Category, error) {

	cachedCategory, err := c.redis.Get(ctx, fmt.Sprintf("category:%d", id)).Result()
	if err == redis.Nil || cachedCategory == "" {

		category, err := c.repo.GetCategoryByID(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("category not found")
		}

		categoryJSON, _ := json.Marshal(category)
		err = c.redis.Set(ctx, fmt.Sprintf("category:%d", id), categoryJSON, 2*time.Hour).Err()
		if err != nil {
			return nil, err
		}

		return category, nil
	}

	var category entity.Category
	err = json.Unmarshal([]byte(cachedCategory), &category)
	if err != nil {
		return nil, err
	}

	return &category, nil
}

func (c *category) DeleteCategoryByID(ctx context.Context, id int) error {

	err := c.repo.DeleteCategoryByID(ctx, id)
	if err != nil {
		return err
	}

	c.redis.Del(ctx, fmt.Sprintf("category:%d", id))
	c.redis.Del(ctx, "categories")
	c.redis.Del(ctx, "product/category")

	return nil
}

func (c *category) UpdateCategory(ctx context.Context, request model.UpdateCategoryRequest) error {

	err := c.repo.UpdateCategory(ctx, &entity.Category{
		ID:   request.CategoryID,
		Name: request.Name,
	})
	if err != nil {
		return err
	}

	c.redis.Del(ctx, fmt.Sprintf("category:%d", request.CategoryID))
	c.redis.Del(ctx, "categories")
	c.redis.Del(ctx, "product/category")

	return nil
}
