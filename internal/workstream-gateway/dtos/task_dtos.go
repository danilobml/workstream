package dtos

type CreateTaskRequest struct {
	Title string `json:"title" validate:"required,min=3"`
}

type SingleTaskResponse struct {
	Id             string `json:"id"`
	Title          string `json:"title"`
	Completed      bool   `json:"completed"`
	OrganizationId string `json:"organization_id"`
}

type ListTasksResponse []SingleTaskResponse
