package service

import (
	"github.com/google/uuid"
	"noda/data/model"
	"noda/data/transfer"
	"noda/data/types"
	"noda/failure"
	"noda/repository"
)

type GroupService interface {
	Save(ownerID uuid.UUID, creation *transfer.GroupCreation) (insertedID uuid.UUID, err error)
	FetchByID(ownerID, groupID uuid.UUID) (group *model.Group, err error)
	Fetch(ownerID uuid.UUID, pagination *types.Pagination, needle, sortExpr string) (result *types.Result[model.Group], err error)
	Update(ownerID, groupID uuid.UUID, update *transfer.GroupUpdate) (ok bool, err error)
	Remove(ownerID, groupID uuid.UUID) (ok bool, err error)
}

type groupService struct {
	r repository.GroupRepository
}

func NewGroupService(repository repository.GroupRepository) GroupService {
	return &groupService{repository}
}

func (s *groupService) Save(ownerID uuid.UUID, newGroup *transfer.GroupCreation) (insertedID uuid.UUID, err error) {
	doTrim(&newGroup.Name, &newGroup.Description)
	switch {
	case 1<<5 < len(newGroup.Name):
		return uuid.Nil, failure.ErrTooLong.Clone().FormatDetails("name", "group", 1<<5)
	case 1<<9 < len(newGroup.Description):
		return uuid.Nil, failure.ErrTooLong.Clone().FormatDetails("description", "group", 1<<9)
	}
	id, err := s.r.Save(ownerID.String(), newGroup)
	if nil != err {
		return uuid.Nil, err
	}
	return uuid.Parse(id)
}

func (s *groupService) FetchByID(ownerID, groupID uuid.UUID) (group *model.Group, err error) {
	return s.r.FetchByID(ownerID.String(), groupID.String())
}

func (s *groupService) Fetch(
	ownerID uuid.UUID,
	pag *types.Pagination,
	needle, sortExpr string,
) (result *types.Result[model.Group], err error) {
	doTrim(&needle, &sortExpr)
	doDefaultPagination(pag)
	groups, err := s.r.Fetch(ownerID.String(), pag.Page, pag.RPP, needle, sortExpr)
	if nil != err {
		return nil, err
	}
	return &types.Result[model.Group]{
		Page:      pag.Page,
		RPP:       pag.RPP,
		Retrieved: int64(len(groups)),
		Payload:   groups,
	}, nil
}

func (s *groupService) Update(ownerID, groupID uuid.UUID, up *transfer.GroupUpdate) (ok bool, err error) {
	doTrim(&up.Name, &up.Description)
	switch {
	case 1<<5 < len(up.Name):
		return false, failure.ErrTooLong.Clone().FormatDetails("name", "group", 1<<5)
	case 1<<9 < len(up.Description):
		return false, failure.ErrTooLong.Clone().FormatDetails("description", "group", 1<<9)
	}
	return s.r.Update(ownerID.String(), groupID.String(), up)
}

func (s *groupService) Remove(ownerID, groupID uuid.UUID) (ok bool, err error) {
	return s.r.Remove(ownerID.String(), groupID.String())
}
