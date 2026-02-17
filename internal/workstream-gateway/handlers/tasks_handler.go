package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/danilobml/workstream/internal/platform/errs"
	"github.com/danilobml/workstream/internal/platform/httputils"
	"github.com/danilobml/workstream/internal/workstream-gateway/dtos"
	services "github.com/danilobml/workstream/internal/workstream-gateway/services/ports"
)

type ITasksHandler interface {
	CreateNewTask(w http.ResponseWriter, r *http.Request)
	GetTask(w http.ResponseWriter, r *http.Request)
	GetTasks(w http.ResponseWriter, r *http.Request)
	CompleteTask(w http.ResponseWriter, r *http.Request)
}

type TasksHandler struct {
	tasksService services.TasksServicePort
}

func NewTasksHandler(tasksService services.TasksServicePort) *TasksHandler {
	return &TasksHandler{
		tasksService: tasksService,
	}
}

func (gh *TasksHandler) CreateNewTask(w http.ResponseWriter, r *http.Request) {
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

	resp := dtos.SingleTaskResponse{
		Id:        tsResp.Id,
		Title:     tsResp.Title,
		Completed: tsResp.Completed,
	}

	err = httputils.WriteJson(w, http.StatusCreated, resp)
	if err != nil {
		log.Println(err)
	}
}

func (gh *TasksHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

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
		Id:        task.Id,
		Title:     task.Title,
		Completed: task.Completed,
	}

	err = httputils.WriteJson(w, http.StatusOK, resp)
	if err != nil {
		log.Println(err)
	}
}

func (gh *TasksHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	tasks, err := gh.tasksService.ListTasks(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var resp dtos.ListTasksResponse

	for _, task := range tasks {
		resp = append(resp, dtos.SingleTaskResponse{
			Id:        task.Id,
			Title:     task.Title,
			Completed: task.Completed,
		})
	}

	err = httputils.WriteJson(w, http.StatusOK, resp)
	if err != nil {
		log.Println(err)
	}
}

func (gh *TasksHandler) CompleteTask(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

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

	err = httputils.WriteJson(w, http.StatusOK, "Task successfully completed")
	if err != nil {
		log.Println(err)
	}
}
