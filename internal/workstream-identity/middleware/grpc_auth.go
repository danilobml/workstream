package middleware

import (
	"context"
	"fmt"
	"strings"

	authcontext "github.com/danilobml/workstream/internal/platform/auth_context"
	"github.com/danilobml/workstream/internal/platform/errs"
	"github.com/danilobml/workstream/internal/platform/jwt"
	"google.golang.org/grpc/metadata"
)

type ctxKey string

const claimsCtxKey ctxKey = "claims"

func AuthenticateGRPC(ctx context.Context, jwtManager *jwt.JwtManager) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx, errs.ErrParsingToken
	}

	vals := md.Get("authorization")
	if len(vals) == 0 {
		return ctx, errs.ErrParsingToken
	}

	parts := strings.Fields(vals[0])
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return ctx, errs.ErrParsingToken
	}

	tokenString := parts[1]
	claims, err := jwtManager.ParseAndValidateToken(tokenString)
	if err != nil {
		return ctx, err
	}
	fmt.Println("identity middleware - claims", claims)


	return authcontext.WithClaims(ctx, claims), nil
}

func GetClaimsFromContext(ctx context.Context) (*jwt.Claims, bool) {
	claims, ok := ctx.Value(claimsCtxKey).(*jwt.Claims)
	return claims, ok
}
