package models

type Task struct {
	TaskId    string `json:"task_id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}
