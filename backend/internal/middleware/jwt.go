package middleware

import (
	"context"
	"strings"

	"github.com/durianpay/fullstack-boilerplate/internal/config"
	"github.com/durianpay/fullstack-boilerplate/internal/entity"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserIDKey contextKey = "user_id"

func JWTAuth(
	ctx context.Context,
	input *openapi3filter.AuthenticationInput,
) error {

	req := input.RequestValidationInput.Request

	authHeader := req.Header.Get("Authorization")

	if authHeader == "" {
		return entity.ErrorUnauthorized("missing authorization header")
	}

	parts := strings.SplitN(authHeader, " ", 2)

	if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
		return entity.ErrorUnauthorized("invalid authorization format")
	}

	token, err := jwt.Parse(parts[1], func(t *jwt.Token) (any, error) {

		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, entity.ErrorUnauthorized("unexpected signing method")
		}

		return config.JwtSecret, nil
	})

	if err != nil || !token.Valid {
		return entity.ErrorUnauthorized("invalid or expired token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return entity.ErrorUnauthorized("invalid token claims")
	}

	sub, _ := claims["sub"].(string)

	// inject user ke context
	newCtx := context.WithValue(
		req.Context(),
		UserIDKey,
		sub,
	)

	input.RequestValidationInput.Request =
		req.WithContext(newCtx)

	return nil
}
