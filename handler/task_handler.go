package handler

import (
	"noda/service"
)

type TaskHandler struct {
	s service.TaskService
}

func NewTaskHandler(service service.TaskService) *TaskHandler {
	return &TaskHandler{s: service}
}
