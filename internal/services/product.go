package services

import (
	"context"
	"fmt"
	"strconv"

	"github.com/aldotp/OnlineStore/internal/entity"
	"github.com/aldotp/OnlineStore/internal/model"
	"github.com/aldotp/OnlineStore/internal/repositories"
)

type ProductService interface {
	GetProductByCategoryID(ctx context.Context, queryParam string) (*model.ProductByCategoryResponse, error)
	StoreProduct(ctx context.Context, request model.ProductRequest) (*model.ProductResponse, error)
}

type product struct {
	repo         *repositories.ProductRepository
	repoCategory *repositories.CategoryRepository
}

func NewProduct(repo *repositories.ProductRepository, repoCategory *repositories.CategoryRepository) ProductService {
	return &product{
		repo:         repo,
		repoCategory: repoCategory,
	}
}

func (p *product) GetProductByCategoryID(ctx context.Context, queryParam string) (*model.ProductByCategoryResponse, error) {

	if queryParam == "" {
		products, err := p.repo.GetAllProducts(ctx)
		if err != nil {
			return nil, err
		}

		return &model.ProductByCategoryResponse{
			Products: products,
		}, nil
	}

	categoryID, err := strconv.Atoi(queryParam)
	if err != nil {
		return nil, err
	}

	category, err := p.repoCategory.GetCategoryByID(ctx, categoryID)
	if err != nil {
		return nil, err
	}

	if category == nil {
		return nil, fmt.Errorf("category not found")
	}

	categoryWithProduct, err := p.repo.GetProductsByCategoryID(category)
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

	return response, nil
}
