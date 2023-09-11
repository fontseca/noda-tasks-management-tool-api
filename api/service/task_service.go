package service

import (
	"noda/api/data/model"
	"noda/api/repository"

	"github.com/google/uuid"
)

type TaskService struct {
	repository *repository.TaskRepository
}

func NewTaskService(repository *repository.TaskRepository) *TaskService {
	return &TaskService{repository}
}

func (s *TaskService) GetByID(id uuid.UUID) (*model.Task, error) {
	return s.repository.GetByID(id)
}

func (s *TaskService) GetAll() (*[]*model.Task, error) {
	return s.repository.GetAll()
}
