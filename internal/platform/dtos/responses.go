package dtos

import "github.com/danilobml/workstream/internal/platform/models"

type RegisterResponse struct {
	Token string `json:"token,omitempty"`
}

type LoginResponse struct {
	Token string `json:"token,omitempty"`
}

type CheckUserResponse struct {
	IsValid bool       `json:"is_valid"`
	User    models.User `json:"user"`
}

type GetAllUsersResponse = []ResponseUser
