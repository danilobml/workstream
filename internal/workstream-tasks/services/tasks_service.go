package services

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	pb "github.com/danilobml/workstream/internal/gen/tasks/v1"
	"github.com/danilobml/workstream/internal/platform/errs"
	"github.com/danilobml/workstream/internal/platform/models"
	"github.com/danilobml/workstream/internal/workstream-tasks/repositories"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TasksService struct {
	pb.UnimplementedTasksServiceServer
	repo          repositories.ITaskRepository
	eventsService EventsService
}

func NewTasksService(repo repositories.ITaskRepository, eventsService EventsService) *TasksService {
	return &TasksService{
		repo: repo,
		eventsService: eventsService,
	}
}

func (ts *TasksService) CreateTask(ctx context.Context, r *pb.CreateTaskRequest) (*pb.CreateTaskResponse, error) {
	if r.GetTitle() == "" {
		return nil, status.Error(codes.InvalidArgument, "required parameter Title is missing")
	}

	id := uuid.New().String()

	task := models.Task{
		Id:        id,
		Title:     r.GetTitle(),
		Completed: false,
	}

	taskDb, err := ts.repo.Create(ctx, task)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create task")
	}

	newTask := &pb.Task{
		TaskId:    taskDb.Id,
		Title:     taskDb.Title,
		Completed: taskDb.Completed,
	}

	taskJson, err := json.Marshal(newTask)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to marshal new task for event")
	}

	event := models.Event{
		EventID: uuid.NewString(),
		EventType: "workstream.task.created.v1",
		OccurredAt: time.Now(),
		Producer: "workstream-tasks",
		TraceID: uuid.NewString(),
		Payload: taskJson,
	}

	err = ts.eventsService.Publish(ctx, event)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to send event")
	}

	return &pb.CreateTaskResponse{
		Task: newTask,
	}, nil
}

func (ts *TasksService) GetTask(ctx context.Context, r *pb.GetTaskRequest) (*pb.GetTaskResponse, error) {
	if r.GetTaskId() == "" {
		return nil, status.Error(codes.InvalidArgument, "required parameter TaskId is missing")
	}

	taskDb, err := ts.repo.GetById(ctx, r.GetTaskId())
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "task not found")
		}
		return nil, status.Error(codes.Internal, "failed to get task")
	}

	foundTask := &pb.Task{
		TaskId:    taskDb.Id,
		Title:     taskDb.Title,
		Completed: taskDb.Completed,
	}

	return &pb.GetTaskResponse{
		Task: foundTask,
	}, nil
}

func (ts *TasksService) ListTasks(ctx context.Context, r *pb.ListTasksRequest) (*pb.ListTasksResponse, error) {
	dbTasks, err := ts.repo.List(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get task list")
	}

	var tasks []*pb.Task
	for _, task := range dbTasks {
		rTask := &pb.Task{
			TaskId:    task.Id,
			Title:     task.Title,
			Completed: task.Completed,
		}
		tasks = append(tasks, rTask)
	}

	resp := &pb.ListTasksResponse{
		Tasks: tasks,
	}

	return resp, nil
}

func (ts *TasksService) CompleteTask(ctx context.Context, r *pb.CompleteTaskRequest) (*pb.CompleteTaskResponse, error) {
	if r.GetTaskId() == "" {
		return nil, status.Error(codes.InvalidArgument, "required parameter TaskId is missing")
	}

	taskToUpdate, err := ts.repo.GetById(ctx, r.GetTaskId())
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "task not found")
		}
		return nil, status.Error(codes.Internal, "failed to get task")
	}

	taskToUpdate.Completed = true
	_, err = ts.repo.Update(ctx, *taskToUpdate)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to update")
	}

	return &pb.CompleteTaskResponse{}, nil
}
