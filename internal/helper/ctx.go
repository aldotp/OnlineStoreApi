package helper

import (
	"context"
	"errors"

	"github.com/aldotp/OnlineStore/internal/middleware"
	"github.com/aldotp/OnlineStore/internal/model"
)

func GetUserCtx(ctx context.Context) (*model.UserCtx, error) {
	claims, ok := ctx.Value("claims").(*middleware.Claims)
	if !ok {
		return nil, errors.New("claims not found in context")
	}
	return &model.UserCtx{
		ID:       claims.ID,
		Username: claims.Username,
	}, nil
}
