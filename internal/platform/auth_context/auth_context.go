package authcontext

import (
	"context"

	"github.com/danilobml/workstream/internal/platform/jwt"
)

type ctxKey string

const claimsKey ctxKey = "claims"

func WithClaims(ctx context.Context, claims *jwt.Claims) context.Context {
	return context.WithValue(ctx, claimsKey, claims)
}

func GetClaims(ctx context.Context) (*jwt.Claims, bool) {
	claims, ok := ctx.Value(claimsKey).(*jwt.Claims)
	return claims, ok
}
