package models

import (
	"github.com/danilobml/workstream/internal/platform/errs"
	"github.com/google/uuid"
)

type Role int

const (
	Admin Role = iota
	AppUser
)

var roleName = map[Role]string{
	Admin:   "admin",
	AppUser: "user",
}

func (r Role) GetName() string {
	return roleName[r]
}

func ParseRole(s string) (Role, error) {
	for r, name := range roleName {
		if name == s {
			return r, nil
		}
	}
	return 0, errs.ErrParsingRoles
}

type User struct {
	ID             uuid.UUID `json:"id"`
	Email          string    `json:"email"`
	HashedPassword string    `json:"-"`
	Roles          []Role    `json:"roles"`
	IsActive       bool      `json:"is_active"`
}

type DbRole struct {
	ID     uuid.UUID `json:"id"`
	UserId uuid.UUID `json:"user_id"`
	Role   string    `json:"role"`
}
