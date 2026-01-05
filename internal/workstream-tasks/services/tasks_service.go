package services

import (
	"context"
	"errors"

	pb "github.com/danilobml/workstream/internal/gen/tasks/v1"
)

type TasksService struct {
	pb.UnimplementedTasksServiceServer
}

func NewTasksService() *TasksService {
	return &TasksService{}
}

func (ts *TasksService) CreateTask(ctx context.Context, r *pb.CreateTaskRequest) (*pb.CreateTaskResponse, error) {
	if r.Title == "" {
		return nil, errors.New("CreateTask - required parameter Title is missing")
	}

	// Dummy task for testing:
	newTask := &pb.Task{
		TaskId: "1",
		Title: r.Title,
		Completed: false,		
	}
	
	return &pb.CreateTaskResponse{
		Task: newTask,
	}, nil
}
