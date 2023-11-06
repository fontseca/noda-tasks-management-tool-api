package service

import (
	"noda/data/model"
	"noda/repository"

	"github.com/google/uuid"
)

type TaskService struct {
	r *repository.TaskRepository
}

func NewTaskService(repository *repository.TaskRepository) *TaskService {
	return &TaskService{repository}
}

func (s *TaskService) GetByID(id uuid.UUID) (*model.Task, error) {
	return s.r.SelectByID(id)
}

func (s *TaskService) GetByUserID(userID uuid.UUID) (*[]*model.Task, error) {
	return s.r.SelectByOwnerID(userID)
}

func (s *TaskService) GetAll() (*[]*model.Task, error) {
	return s.r.SelectAll()
}
