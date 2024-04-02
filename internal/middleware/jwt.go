package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/aldotp/OnlineStore/internal/config"
	"github.com/aldotp/OnlineStore/internal/entity"
	"github.com/dgrijalva/jwt-go"
)

// Claims adalah struktur yang digunakan untuk menyimpan klaim JWT
type Claims struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

type JWT struct {
	config *config.BootstrapConfig
}

func NewJWT(config *config.BootstrapConfig) *JWT {
	return &JWT{
		config: config,
	}
}

// GenerateJWT menghasilkan JWT menggunakan username sebagai klaim
func (j *JWT) GenerateJWT(user *entity.User) (string, string, error) {

	expirationTime := time.Now().Add(2 * time.Hour)

	claims := &Claims{
		ID:       user.ID,
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(j.config.JWTKey))
	if err != nil {
		return "", "", err
	}

	remainingTime := expirationTime.Unix() - time.Now().Unix()
	remainingTimeString := fmt.Sprintf("%d seconds", remainingTime)

	return tokenString, remainingTimeString, nil
}

func (j *JWT) ValidateJWT(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.config.JWTKey), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}

func (j *JWT) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("Authorization")
		if authorizationHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		tokenString := strings.Replace(authorizationHeader, "Bearer ", "", 1)

		claims, err := j.ValidateJWT(tokenString)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), "claims", claims))

		next.ServeHTTP(w, r)
	})
}
