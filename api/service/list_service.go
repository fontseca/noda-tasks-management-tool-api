package service

import (
	"github.com/google/uuid"
	"noda/api/data/transfer"
	"noda/api/repository"
)

type ListService struct {
	r repository.IListRepository
}

func NewListService(r repository.IListRepository) *ListService {
	return &ListService{r}
}

func (s *ListService) SaveList(ownerID, groupID uuid.UUID, next *transfer.ListCreation) (insertedID uuid.UUID, err error) {
	id, err := s.r.InsertList(ownerID.String(), groupID.String(), next)
	if nil != err {
		return uuid.Nil, err
	}
	return uuid.Parse(id)
}
