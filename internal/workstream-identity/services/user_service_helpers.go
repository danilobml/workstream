package services

import (
	"context"

	authcontext "github.com/danilobml/workstream/internal/platform/auth_context"
	"github.com/danilobml/workstream/internal/platform/httpx/middleware"
	"github.com/danilobml/workstream/internal/platform/models"
)

// Helpers
func (us *UserService) IsUserOwner(ctx context.Context, userEmail string) bool {
	claims, ok := middleware.GetClaimsFromContext(ctx)
	if !ok {
		return false
	}

	_, err := us.userRepository.FindByEmail(ctx, claims.Email)
	if err != nil {
		return false
	}

	return claims.Email == userEmail
}

func (us *UserService) IsUserAdmin(ctx context.Context) bool {
	claims, ok := authcontext.GetClaims(ctx)
	if !ok {
		return false
	}

	_, err := us.userRepository.FindByEmail(ctx, claims.Email)
	if err != nil {
		return false
	}

	for _, role := range claims.Roles {
		if role == models.Admin {
			return true
		}
	}

	return false
}
