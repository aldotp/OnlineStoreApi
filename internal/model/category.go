package model

import "github.com/aldotp/OnlineStore/internal/entity"

type CategoryRequest struct {
	Name string `json:"name"`
}

type ProductByCategoryResponse struct {
	ID           int              `json:"category_id,omitempty"`
	CategoryName string           `json:"category_name,omitempty"`
	Products     []entity.Product `json:"products"`
}

type UpdateCategoryRequest struct {
	CategoryID int    `json:"category_id"`
	Name       string `json:"name"`
}
