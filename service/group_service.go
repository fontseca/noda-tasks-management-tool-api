package service

import (
	"github.com/google/uuid"
	"noda"
	"noda/data/model"
	"noda/data/transfer"
	"noda/data/types"
	"noda/repository"
)

type GroupService interface {
	Save(ownerID uuid.UUID, creation *transfer.GroupCreation) (insertedID string, err error)
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

func (s *groupService) Save(ownerID uuid.UUID, newGroup *transfer.GroupCreation) (insertedID string, err error) {
	if len(newGroup.Name) > 50 {
		return "", noda.ErrTooLong.Clone().FormatDetails("name", "group", 50)
	}
	return s.r.Save(ownerID.String(), newGroup)
}

func (s *groupService) FetchByID(ownerID, groupID uuid.UUID) (group *model.Group, err error) {
	return s.r.FetchByID(ownerID.String(), groupID.String())
}

func (s *groupService) Fetch(
	ownerID uuid.UUID,
	pag *types.Pagination,
	needle, sortExpr string,
) (result *types.Result[model.Group], err error) {
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
	if len(up.Name) > 50 {
		return false, noda.ErrTooLong.Clone().FormatDetails("name", "group", 50)
	}
	return s.r.Update(ownerID.String(), groupID.String(), up)
}

func (s *groupService) Remove(ownerID, groupID uuid.UUID) (ok bool, err error) {
	return s.r.Remove(ownerID.String(), groupID.String())
}
