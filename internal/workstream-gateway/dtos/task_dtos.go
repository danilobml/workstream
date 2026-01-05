package dtos

type CreateTaskRequest struct {
	Title string `json:"title"`
}

type CreateTaskResponse struct {
	Id        string `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

type GetTaskResponse struct {
	Id        string `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}
