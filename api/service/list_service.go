package service

import "noda/api/repository"

type ListService struct {
	r repository.IListRepository
}

func NewListService(r repository.IListRepository) *ListService {
	return &ListService{r}
}
