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

type ProductService interface {
	GetProductByCategoryID(ctx context.Context, id int) (*model.ProductByCategoryResponse, error)
	StoreProduct(ctx context.Context, request model.ProductRequest) (*model.ProductResponse, error)
	UpdateProduct(ctx context.Context, request model.UpdateProductRequest) error
	DeleteProduct(ctx context.Context, request model.DeleteProductRequest) error
	GetProducts(ctx context.Context) ([]model.ProductResponse, error)
	GetProductByID(ctx context.Context, id int) (*model.ProductResponse, error)
}

type product struct {
	repo         *repositories.ProductRepository
	repoCategory *repositories.CategoryRepository
	redis        *redis.Client
}

func NewProduct(repo *repositories.ProductRepository, repoCategory *repositories.CategoryRepository, redis *redis.Client) ProductService {
	return &product{
		repo:         repo,
		repoCategory: repoCategory,
		redis:        redis,
	}
}

func (p *product) GetProductByCategoryID(ctx context.Context, id int) (*model.ProductByCategoryResponse, error) {

	category, err := p.repoCategory.GetCategoryByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if category == nil {
		return nil, fmt.Errorf("category not found")
	}

	categoryWithProduct, err := p.repo.GetProductsByCategoryID(ctx, category)
	if err != nil {
		return nil, err
	}

	response := model.ProductByCategoryResponse{
		ID:           category.ID,
		CategoryName: category.Name,
		Products:     categoryWithProduct.Products,
	}

	return &response, nil

}

func (p *product) StoreProduct(ctx context.Context, request model.ProductRequest) (*model.ProductResponse, error) {

	categoryID, err := p.repoCategory.GetCategoryByID(ctx, request.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("invalid category id")
	}

	if categoryID == nil {
		return nil, fmt.Errorf("category id not found")
	}

	product := entity.Product{
		Name:        request.Name,
		Description: request.Description,
		Price:       request.Price,
		CategoryID:  request.CategoryID,
	}

	insertedProduct, err := p.repo.StoreProduct(ctx, product)
	if err != nil {
		return nil, err
	}

	response := &model.ProductResponse{
		ID:          insertedProduct.ID,
		Name:        insertedProduct.Name,
		Description: insertedProduct.Description,
		Price:       insertedProduct.Price,
		CategoryID:  insertedProduct.CategoryID,
		CreatedAt:   insertedProduct.CreatedAt.String(),
		UpdatedAt:   insertedProduct.UpdatedAt.String(),
	}

	p.redis.Del(ctx, "products")
	p.redis.Del(ctx, "product/category")

	return response, nil
}

func (p *product) UpdateProduct(ctx context.Context, request model.UpdateProductRequest) error {

	product, err := p.repo.GetProductByID(ctx, request.ProductID)
	if err != nil {
		return err
	}

	if product == nil {
		return fmt.Errorf("product not found")
	}

	err = p.repo.UpdateProduct(ctx, &entity.Product{
		ID:          product.ID,
		Name:        request.Name,
		Description: request.Description,
		Price:       request.Price,
	})
	if err != nil {
		return err
	}

	p.redis.Del(ctx, fmt.Sprintf("product:%d", request.ProductID))
	p.redis.Del(ctx, "products")

	return nil

}

func (p *product) DeleteProduct(ctx context.Context, request model.DeleteProductRequest) error {

	product, err := p.repo.GetProductByID(ctx, request.ProductID)
	if err != nil {
		return err
	}

	if product == nil {
		return fmt.Errorf("product not found")
	}

	err = p.repo.DeleteProduct(ctx, product.ID)
	if err != nil {
		return err
	}

	p.redis.Del(ctx, fmt.Sprintf("product:%d", request.ProductID))
	p.redis.Del(ctx, "products")

	return nil

}

func (p *product) GetProducts(ctx context.Context) ([]model.ProductResponse, error) {
	cachedProducts, err := p.redis.Get(ctx, "products").Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	if err == redis.Nil {
		products, err := p.repo.GetAllProducts(ctx)
		if err != nil {
			return nil, err
		}

		responses := make([]model.ProductResponse, len(products))

		for i, product := range products {
			responses[i] = model.ProductResponse{
				ID:          product.ID,
				Name:        product.Name,
				Description: product.Description,
				Price:       product.Price,
				CategoryID:  product.CategoryID,
				CreatedAt:   product.CreatedAt.String(),
				UpdatedAt:   product.UpdatedAt.String(),
			}
		}

		productJSON, err := json.Marshal(responses)
		if err != nil {
			return nil, err
		}

		err = p.redis.Set(ctx, "products", productJSON, 2*time.Hour).Err()
		if err != nil {
			return nil, err
		}

		return responses, nil
	}

	var responses []model.ProductResponse
	err = json.Unmarshal([]byte(cachedProducts), &responses)
	if err != nil {
		return nil, err
	}

	return responses, nil
}

func (p *product) GetProductByID(ctx context.Context, id int) (*model.ProductResponse, error) {

	cachedProduct, err := p.redis.Get(ctx, fmt.Sprintf("product:%d", id)).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	if err == redis.Nil {
		product, err := p.repo.GetProductByID(ctx, id)
		if err != nil {
			return nil, err
		}

		if product == nil {
			return nil, fmt.Errorf("product not found")
		}

		productResponse := &model.ProductResponse{
			ID:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			CategoryID:  product.CategoryID,
			CreatedAt:   product.CreatedAt.String(),
			UpdatedAt:   product.UpdatedAt.String(),
		}

		productJSON, err := json.Marshal(productResponse)
		if err != nil {
			return nil, err
		}

		err = p.redis.Set(ctx, fmt.Sprintf("product:%d", id), productJSON, 2*time.Hour).Err()
		if err != nil {
			return nil, err
		}

		return productResponse, nil

	}

	var productResponse model.ProductResponse
	err = json.Unmarshal([]byte(cachedProduct), &productResponse)
	if err != nil {
		return nil, err
	}

	return &productResponse, nil

}
