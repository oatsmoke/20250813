package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/oatsmoke/20250813/internal/model"
)

type taskService interface {
	Start(duration int64) string
	Get(key int64) *model.Task
	GetAll(status string) map[int64]*model.Task
}

type TaskHandler struct {
	taskService taskService
}

func NewTaskHandler(taskService taskService) *TaskHandler {
	return &TaskHandler{taskService: taskService}
}

func (h *TaskHandler) Start(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Duration int64 `json:"duration"`
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	res := h.taskService.Start(data.Duration)

	if _, err := w.Write([]byte(res)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *TaskHandler) Get(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	key, err := strconv.ParseInt(strings.TrimPrefix(r.URL.Path, "/tasks/"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res := h.taskService.Get(key)

	if err := json.NewEncoder(w).Encode(res); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *TaskHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")

	res := h.taskService.GetAll(status)

	if err := json.NewEncoder(w).Encode(res); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
