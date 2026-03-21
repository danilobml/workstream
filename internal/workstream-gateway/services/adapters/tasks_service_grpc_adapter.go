package adapters

import (
	"context"

	pb "github.com/danilobml/workstream/internal/gen/tasks/v1"
	"github.com/danilobml/workstream/internal/platform/grpcutils"
	"github.com/danilobml/workstream/internal/platform/models"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type TasksServiceClient struct {
	pb pb.TasksServiceClient
}

func NewTasksServiceClient(conn grpc.ClientConnInterface) *TasksServiceClient {
	return &TasksServiceClient{pb: pb.NewTasksServiceClient(conn)}
}

func (c *TasksServiceClient) CreateTask(ctx context.Context, title string) (*models.Task, error) {
	resp, err := c.pb.CreateTask(ctx, &pb.CreateTaskRequest{Title: title})
	if err != nil {
		return nil, grpcutils.ParseGrpcError(err)
	}

	task := resp.GetTask()
	organizationIdParsed, err := uuid.Parse(task.GetOrganizationId())
	if err != nil {
		return nil, grpcutils.ParseGrpcError(err)
	}
	return &models.Task{
		Id:             task.GetTaskId(),
		Title:          task.GetTitle(),
		Completed:      task.GetCompleted(),
		OrganizationId: organizationIdParsed,
	}, nil
}

func (c *TasksServiceClient) GetTask(ctx context.Context, id string) (*models.Task, error) {
	resp, err := c.pb.GetTask(ctx, &pb.GetTaskRequest{TaskId: id})
	if err != nil {
		return nil, grpcutils.ParseGrpcError(err)
	}

	task := resp.GetTask()

	organizationIdParsed, err := uuid.Parse(task.GetOrganizationId())
	if err != nil {
		return nil, grpcutils.ParseGrpcError(err)
	}

	return &models.Task{
		Id:             task.GetTaskId(),
		Title:          task.GetTitle(),
		Completed:      task.GetCompleted(),
		OrganizationId: organizationIdParsed,
	}, nil
}

func (c *TasksServiceClient) ListTasks(ctx context.Context) ([]*models.Task, error) {
	resp, err := c.pb.ListTasks(ctx, &pb.ListTasksRequest{})
	if err != nil {
		return nil, grpcutils.ParseGrpcError(err)
	}

	var tasks []*models.Task

	for _, task := range resp.GetTasks() {
		organizationIdParsed, err := uuid.Parse(task.GetOrganizationId())
		if err != nil {
			return nil, grpcutils.ParseGrpcError(err)
		}
		tasks = append(tasks, &models.Task{
			Id:             task.GetTaskId(),
			Title:          task.GetTitle(),
			Completed:      task.GetCompleted(),
			OrganizationId: organizationIdParsed,
		})
	}

	return tasks, nil
}

func (c *TasksServiceClient) CompleteTask(ctx context.Context, id string) error {
	_, err := c.pb.CompleteTask(ctx, &pb.CompleteTaskRequest{TaskId: id})
	if err != nil {
		return grpcutils.ParseGrpcError(err)
	}

	return nil
}
