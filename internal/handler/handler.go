package handler

import (
	"net/http"
)

type Handler struct {
	taskHandler *TaskHandler
}

func NewHandler(taskService taskService) *Handler {
	return &Handler{
		taskHandler: NewTaskHandler(taskService),
	}
}

func (h *Handler) InitRouts() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			h.taskHandler.GetAll(w, r)
		case "POST":
			h.taskHandler.Start(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/tasks/", h.taskHandler.Get)
	return mux
}
