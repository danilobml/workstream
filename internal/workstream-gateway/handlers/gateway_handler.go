package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/danilobml/workstream/internal/platform/httputils"
	"github.com/danilobml/workstream/internal/workstream-gateway/dtos"
	services "github.com/danilobml/workstream/internal/workstream-gateway/services/ports"
)

type IGatewayHandler interface {
	CreateNewTask(w http.ResponseWriter, r *http.Request)
}

type GatewayHandler struct {
	tasksService services.TasksServicePort
}

func NewGatewayHandler(tasksService services.TasksServicePort) *GatewayHandler {
	return &GatewayHandler{
		tasksService: tasksService,
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

	tsResp, err := gh.tasksService.CreateTask(ctx, request.Title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := dtos.CreateTaskResponse{
		Id: tsResp.Id,
		Title: tsResp.Title,
		Completed: tsResp.Completed,
	}

	err = httputils.WriteJson(w, http.StatusCreated, resp)
	if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}
