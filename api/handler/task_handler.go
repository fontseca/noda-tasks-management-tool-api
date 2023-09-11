package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"noda/api/service"
	"noda/failure"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type TaskHandler struct {
	service *service.TaskService
}

func NewTaskHandler(service *service.TaskService) *TaskHandler {
	return &TaskHandler{service}
}

func (h *TaskHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "task_id"))
	if err != nil {
		failure.Emit(w, http.StatusBadRequest,
			"failure with `task_id'", err.Error())
		return
	}

	task, err := h.service.GetByID(id)
	if err != nil {
		switch {
		default:
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		case errors.Is(err, failure.ErrNotFound):
			failure.Emit(w, http.StatusNotFound,
				"record not found", fmt.Sprintf("could not find task with ID `%s'", id))
			return
		}
	}

	res, err := json.Marshal(task)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(res)
}

func (h *TaskHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.service.GetAll()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(*tasks)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(res)
}

func (h *TaskHandler) CreateOneForLoggedInUser(w http.ResponseWriter, r *http.Request) {}

func (h *TaskHandler) GetAllFromLoggedInUser(w http.ResponseWriter, r *http.Request) {}

func (h *TaskHandler) GetByIDFromLoggedInUser(w http.ResponseWriter, r *http.Request) {}

func (h *TaskHandler) UpdateByIDFromLoggedInUser(w http.ResponseWriter, r *http.Request) {}

func (h *TaskHandler) DeleteByIDFromLoggedInUser(w http.ResponseWriter, r *http.Request) {}
