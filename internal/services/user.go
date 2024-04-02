package services

import (
	"context"
	"fmt"

	"github.com/aldotp/OnlineStore/internal/config"
	"github.com/aldotp/OnlineStore/internal/entity"
	"github.com/aldotp/OnlineStore/internal/helper"
	"github.com/aldotp/OnlineStore/internal/middleware"
	"github.com/aldotp/OnlineStore/internal/model"
	"github.com/aldotp/OnlineStore/internal/repositories"
)

type UserService interface {
	CreateUser(ctx context.Context, request model.RegisterRequest) (*model.RegisterResponse, error)
	LoginUser(ctx context.Context, user model.LoginRequest) (*model.LoginResponse, error)
}

type user struct {
	repo   *repositories.UserRepository
	config *config.BootstrapConfig
}

// New User create new instance of User
func NewUser(repo *repositories.UserRepository, config *config.BootstrapConfig) UserService {
	return &user{
		repo:   repo,
		config: config,
	}
}

func (u *user) CreateUser(ctx context.Context, request model.RegisterRequest) (*model.RegisterResponse, error) {

	if request.Password != request.ConfirmPassword {
		return nil, fmt.Errorf("confirm password do not match")
	}

	usr, _ := u.repo.GetUserByUsername(ctx, request.Username)
	if usr != nil {
		return nil, fmt.Errorf("username already exists")
	}

	hashedPassword, err := helper.HashPassword(request.Password)
	if err != nil {
		return nil, err
	}

	usr, err = u.repo.CreateUser(ctx, entity.User{
		Username: request.Username,
		Email:    request.Email,
		Password: hashedPassword,
	})
	if err != nil {
		return nil, err
	}

	newUser := model.RegisterResponse{
		ID:        usr.ID,
		Username:  usr.Username,
		Email:     usr.Email,
		CreatedAt: usr.CreatedAt.String(),
		UpdatedAt: usr.UpdatedAt.String(),
	}

	return &newUser, nil
}

func (u *user) LoginUser(ctx context.Context, request model.LoginRequest) (*model.LoginResponse, error) {

	usr, err := u.repo.GetUserByUsername(ctx, request.Username)
	if err != nil {
		return nil, err
	}

	if usr == nil {
		return nil, fmt.Errorf("invalid username or password")
	}

	if !helper.ComparePassword(usr.Password, request.Password) {
		return nil, fmt.Errorf("invalid username or password")
	}

	jwt := middleware.NewJWT(u.config)

	token, expTime, err := jwt.GenerateJWT(usr)
	if err != nil {
		return nil, err
	}

	return &model.LoginResponse{
		Token:   token,
		Expired: expTime,
	}, nil

}
