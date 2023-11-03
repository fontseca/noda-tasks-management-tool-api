package service

import (
	"errors"
	"github.com/google/uuid"
	"noda/api/data/transfer"
	"noda/api/repository"
	"strings"
)

type ListService struct {
	r repository.IListRepository
}

func NewListService(r repository.IListRepository) *ListService {
	return &ListService{r}
}

func (s *ListService) SaveList(ownerID, groupID uuid.UUID, next *transfer.ListCreation) (insertedID uuid.UUID, err error) {
	next.Name = strings.Trim(next.Name, " \t\n")
	next.Description = strings.Trim(next.Description, " \t\n")
	if "" == next.Name {
		return uuid.Nil, errors.New("name cannot be an empty string")
	}
	if 50 < len(next.Name) {
		return uuid.Nil, errors.New("name too long for list: max length must be 50")
	}
	id, err := s.r.InsertList(ownerID.String(), groupID.String(), next)
	if nil != err {
		return uuid.Nil, err
	}
	return uuid.Parse(id)
}