package handler

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"log"
	"net/http"
	"noda"
	"noda/api/data/types"
	"noda/api/service"
)

type TaskHandler struct {
	s *service.TaskService
}

func NewTaskHandler(service *service.TaskService) *TaskHandler {
	return &TaskHandler{service}
}

func (h *TaskHandler) RetrieveTaskByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "task_id"))
	if err != nil {
		//noda.Emit(w, http.StatusBadRequest, "failure with \"task_id\"", err)
		return
	}

	task, err := h.s.GetByID(id)
	if err != nil {
		switch {
		default:
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		case errors.Is(err, noda.ErrUserNotFound):
			//noda.Emit(w, http.StatusNotFound,
			//	"record not found", fmt.Sprintf("could not find any task with ID %q", id))
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

func (h *TaskHandler) RetrieveAll(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.s.GetAll()
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

func (h *TaskHandler) RetrieveTasksFromUser(w http.ResponseWriter, r *http.Request) {
	jwtPayload := r.Context().Value(types.ContextKey{}).(types.JWTPayload)
	tasks, err := h.s.GetByUserID(jwtPayload.UserID)
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
