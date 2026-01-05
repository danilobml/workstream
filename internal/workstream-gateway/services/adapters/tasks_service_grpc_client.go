package adapters

import (
	"context"
	"fmt"

	pb "github.com/danilobml/workstream/internal/gen/tasks/v1"
	"github.com/danilobml/workstream/internal/workstream-gateway/models"
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
		return nil, fmt.Errorf("tasks grpc CreateTask: %w", err)
	}

	t := resp.GetTask()
	return &models.Task{
		Id:        t.GetTaskId(),
		Title:     t.GetTitle(),
		Completed: t.GetCompleted(),
	}, nil
}