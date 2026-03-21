package models

import "github.com/google/uuid"

type Task struct {
	Id             string
	Title          string
	Completed      bool
	OrganizationId uuid.UUID
}
