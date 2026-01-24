package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/golang-jwt/jwt/v5"

	common "sun-stockanalysis-api/internal/common"
)

type authContextKey struct{}

var userIDContextKey authContextKey

func authMiddleware(secret, issuer string) func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		if secret == "" {
			writeAuthError(ctx, http.StatusUnauthorized, "auth secret not configured")
			return
		}

		authHeader := ctx.Header("Authorization")

		if authHeader == "" || !strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
			writeAuthError(ctx, http.StatusUnauthorized, "missing or invalid authorization header")
			return
		}

		tokenString := strings.TrimSpace(authHeader[len("bearer "):])

		if tokenString == "" {
			writeAuthError(ctx, http.StatusUnauthorized, "missing token")
			return
		}

		claims := &jwt.RegisteredClaims{}

		options := []jwt.ParserOption{
			jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
		}
		if issuer != "" {
			options = append(options, jwt.WithIssuer(issuer))
		}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte(secret), nil
		}, options...)

		if err != nil || !token.Valid || claims.Subject == "" {
			log.Printf(
				"auth invalid token: err=%v valid=%v sub=%q iss=%q exp=%v iat=%v",
				err,
				token != nil && token.Valid,
				claims.Subject,
				claims.Issuer,
				claims.ExpiresAt,
				claims.IssuedAt,
			)
			writeAuthError(ctx, http.StatusUnauthorized, "invalid token")
			return
		}

		next(huma.WithValue(ctx, userIDContextKey, claims.Subject))
	}
}

func writeAuthError(ctx huma.Context, status int, message string) {
	ctx.SetStatus(status)
	ctx.SetHeader("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(ctx.BodyWriter()).Encode(common.NewErrorResponse(status, message))
}
