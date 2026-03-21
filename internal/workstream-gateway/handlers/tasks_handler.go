package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/danilobml/workstream/internal/platform/errs"
	"github.com/danilobml/workstream/internal/workstream-gateway/dtos"
	"github.com/danilobml/workstream/internal/workstream-gateway/httputils"
	services "github.com/danilobml/workstream/internal/workstream-gateway/services/ports"
)

type TasksHandler struct {
	tasksService services.TasksServicePort
}

func NewTasksHandler(tasksService services.TasksServicePort) *TasksHandler {
	return &TasksHandler{
		tasksService: tasksService,
	}
}

func (gh *TasksHandler) CreateNewTask(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	ctx := httputils.CtxWithAuth(r.Context(), auth)

	request := &dtos.CreateTaskRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "failed to parse request", http.StatusBadRequest)
		return
	}

	if !isInputValid(w, request) {
		return
	}

	tsResp, err := gh.tasksService.CreateTask(ctx, request.Title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := dtos.SingleTaskResponse{
		Id:             tsResp.Id,
		Title:          tsResp.Title,
		Completed:      tsResp.Completed,
		OrganizationId: tsResp.OrganizationId.String(),
	}

	httputils.WriteJSONResponse(w, http.StatusCreated, resp)
}

func (gh *TasksHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	ctx := httputils.CtxWithAuth(r.Context(), auth)

	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "error: id missing in path", http.StatusBadRequest)
		return
	}

	task, err := gh.tasksService.GetTask(ctx, id)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := &dtos.SingleTaskResponse{
		Id:             task.Id,
		Title:          task.Title,
		Completed:      task.Completed,
		OrganizationId: task.OrganizationId.String(),
	}

	httputils.WriteJSONResponse(w, http.StatusOK, resp)
}

func (gh *TasksHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	ctx := httputils.CtxWithAuth(r.Context(), auth)

	tasks, err := gh.tasksService.ListTasks(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var resp dtos.ListTasksResponse

	for _, task := range tasks {
		resp = append(resp, dtos.SingleTaskResponse{
			Id:             task.Id,
			Title:          task.Title,
			Completed:      task.Completed,
			OrganizationId: task.OrganizationId.String(),
		})
	}
	httputils.WriteJSONResponse(w, http.StatusOK, resp)
}

func (gh *TasksHandler) CompleteTask(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	ctx := httputils.CtxWithAuth(r.Context(), auth)

	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "error: id missing in path", http.StatusBadRequest)
		return
	}

	err := gh.tasksService.CompleteTask(ctx, id)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	httputils.WriteJSONResponse(w, http.StatusOK, "Task successfully completed")
}
