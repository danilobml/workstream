package models

import "github.com/google/uuid"

type Organization struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
