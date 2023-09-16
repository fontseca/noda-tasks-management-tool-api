package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"noda/api/service"
)

type UserHandler struct {
	s *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service}
}

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.s.GetAll()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	res, err := json.Marshal(*users)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(res)
}
