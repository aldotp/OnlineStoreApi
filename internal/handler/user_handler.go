package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/aldotp/OnlineStore/internal/helper"
	"github.com/aldotp/OnlineStore/internal/model"
	"github.com/aldotp/OnlineStore/internal/services"
)

type UserHandler struct {
	UserService services.UserService
}

func NewUserHandler(userSvc services.UserService) *UserHandler {
	return &UserHandler{
		UserService: userSvc,
	}
}

func (u *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	request := model.LoginRequest{}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		helper.ErrorJSON(helper.Response{
			Code:    http.StatusBadRequest,
			Message: "invalid json body",
		}, w, http.StatusBadRequest)
		return
	}

	loginResponse, err := u.UserService.LoginUser(ctx, request)
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
		Data:    loginResponse,
	})
}

func (u *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {

	var request model.RegisterRequest

	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel()

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		helper.ErrorJSON(helper.Response{
			Code:    http.StatusBadRequest,
			Message: "invalid json body",
		}, w, http.StatusBadRequest)
		return
	}

	res, err := u.UserService.CreateUser(ctx, request)
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
		Data:    res,
	}, nil)

}
