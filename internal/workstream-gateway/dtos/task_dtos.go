package dtos

type CreateTaskRequest struct {
	Title string `json:"title"`
}

type SingleTaskResponse struct {
	Id        string `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

type ListTasksResponse []SingleTaskResponse
