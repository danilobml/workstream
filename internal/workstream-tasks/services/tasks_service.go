package services

import (
	"context"

	pb "github.com/danilobml/workstream/internal/gen/tasks/v1"
	"github.com/danilobml/workstream/internal/platform/models"
	"github.com/danilobml/workstream/internal/workstream-tasks/repositories"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TasksService struct {
	pb.UnimplementedTasksServiceServer
	repo repositories.ITaskRepository
}

func NewTasksService(repo repositories.ITaskRepository) *TasksService {
	return &TasksService{
		repo: repo,
	}
}

func (ts *TasksService) CreateTask(ctx context.Context, r *pb.CreateTaskRequest) (*pb.CreateTaskResponse, error) {
	if r.Title == "" {
		return nil, status.Error(codes.InvalidArgument, "required parameter Title is missing")
	}

	id := uuid.New().String()

	task := models.Task{
		Id: id,
		Title: r.Title,
		Completed: false,
	}

	taskDb, err := ts.repo.Create(ctx, task)
	if err != nil {
		return nil, err
	}

	newTask := &pb.Task{
		TaskId:    taskDb.Id,
		Title:     taskDb.Title,
		Completed: taskDb.Completed,
	}

	return &pb.CreateTaskResponse{
		Task: newTask,
	}, nil
}

func (ts *TasksService) GetTask(ctx context.Context, r *pb.GetTaskRequest) (*pb.GetTaskResponse, error) {
	if r.TaskId == "" {
		return nil, status.Error(codes.InvalidArgument, "required parameter TaskId is missing")
	}

	if r.TaskId != "1" {
		return nil, status.Error(codes.NotFound, "not found")
	}

	// Dummy task for testing:
	// TODO: implement fetching logic in repo:
	foundTask := &pb.Task{
		TaskId:    r.TaskId,
		Title:     "Test task",
		Completed: false,
	}

	return &pb.GetTaskResponse{
		Task: foundTask,
	}, nil
}

func (ts *TasksService) ListTasks(ctx context.Context, r *pb.ListTasksRequest) (*pb.ListTasksResponse, error) {
	// TODO - implement fetching logic
	resp := &pb.ListTasksResponse{
		Tasks: []*pb.Task{
			{
				TaskId:    "1",
				Title:     "Test task 1",
				Completed: false,
			},
			{
				TaskId:    "2",
				Title:     "Test task 2",
				Completed: true,
			},
			{
				TaskId:    "3",
				Title:     "Test task 3",
				Completed: false,
			},
		},
	}

	return resp, nil
}

func (ts *TasksService) CompleteTask(ctx context.Context, r *pb.CompleteTaskRequest) (*pb.CompleteTaskResponse, error) {
	if r.TaskId == "" {
		return nil, status.Error(codes.InvalidArgument, "required parameter TaskId is missing")
	}

	if r.TaskId != "1" {
		return nil, status.Error(codes.NotFound, "not found")
	}
	// TODO - implement update logic

	return &pb.CompleteTaskResponse{}, nil
}
