package adapters

import (
	"context"

	pb "github.com/danilobml/workstream/internal/gen/tasks/v1"
	"github.com/danilobml/workstream/internal/platform/grpcutils"
	"github.com/danilobml/workstream/internal/platform/models"
	"google.golang.org/grpc"
)

type Client struct {
	pb pb.TasksServiceClient
}

func NewTasksServiceClient(conn grpc.ClientConnInterface) *Client {
	return &Client{pb: pb.NewTasksServiceClient(conn)}
}

func (c *Client) CreateTask(ctx context.Context, title string) (*models.Task, error) {
	resp, err := c.pb.CreateTask(ctx, &pb.CreateTaskRequest{Title: title})
	if err != nil {
		return nil, grpcutils.ParseGrpcError(err)
	}

	t := resp.GetTask()
	return &models.Task{
		Id:        t.GetTaskId(),
		Title:     t.GetTitle(),
		Completed: t.GetCompleted(),
	}, nil
}

func (c *Client) GetTask(ctx context.Context, id string) (*models.Task, error) {
	resp, err := c.pb.GetTask(ctx, &pb.GetTaskRequest{TaskId: id})
	if err != nil {
		return nil, grpcutils.ParseGrpcError(err)
	}

	t := resp.GetTask()
	return &models.Task{
		Id:        t.GetTaskId(),
		Title:     t.GetTitle(),
		Completed: t.GetCompleted(),
	}, nil
}

func (c *Client) ListTasks(ctx context.Context) ([]*models.Task, error) {
	resp, err := c.pb.ListTasks(ctx, &pb.ListTasksRequest{})
	if err != nil {
		return nil, grpcutils.ParseGrpcError(err)
	}

	var tasks []*models.Task

	for _, task := range resp.GetTasks() {
		tasks = append(tasks, &models.Task{
			Id: task.GetTaskId(),
			Title: task.GetTitle(),
			Completed: task.GetCompleted(),
		})
	}

	return tasks, nil
}

func (c *Client) CompleteTask(ctx context.Context, id string) error {
	_, err := c.pb.CompleteTask(ctx, &pb.CompleteTaskRequest{TaskId: id})
	if err != nil {
		return grpcutils.ParseGrpcError(err)
	}

	return nil
}
