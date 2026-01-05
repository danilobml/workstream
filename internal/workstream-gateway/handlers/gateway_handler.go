package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	pb "github.com/danilobml/workstream/internal/gen/tasks/v1"
	"github.com/danilobml/workstream/internal/platform/httputils"
	"github.com/danilobml/workstream/internal/workstream-gateway/dtos"
	"google.golang.org/grpc"
)

type IGatewayHandler interface {
	CreateNewTask(w http.ResponseWriter, r *http.Request)
}

type GatewayHandler struct {
	tasksService pb.TasksServiceClient
}

func NewGatewayHandler(conn grpc.ClientConnInterface) *GatewayHandler {
	return &GatewayHandler{
		tasksService: pb.NewTasksServiceClient(conn),
	}
}

func (gh *GatewayHandler) CreateNewTask(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	request := &dtos.CreateTaskRequest{}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "failed to parse request", http.StatusBadRequest)
		return
	}

	tsResp, err := gh.tasksService.CreateTask(ctx, &pb.CreateTaskRequest{Title: request.Title})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	newTask := tsResp.GetTask()

	resp := dtos.CreateTaskResponse{
		Id: newTask.TaskId,
		Title: newTask.Title,
		Completed: newTask.Completed,
	}

	err = httputils.WriteJson(w, http.StatusCreated, resp)
	if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}
