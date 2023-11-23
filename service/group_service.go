package service

import (
	"github.com/google/uuid"
	"noda"
	"noda/data/model"
	"noda/data/transfer"
	"noda/data/types"
	"noda/repository"
)

type GroupService struct {
	r repository.GroupRepository
}

func NewGroupService(repository repository.GroupRepository) *GroupService {
	return &GroupService{repository}
}

func (s *GroupService) SaveGroup(ownerID uuid.UUID, newGroup *transfer.GroupCreation) (insertedID string, err error) {
	if len(newGroup.Name) > 50 {
		return "", noda.ErrTooLong.Clone().FormatDetails("name", "group", 50)
	}
	return s.r.Save(ownerID.String(), newGroup)
}

func (s *GroupService) FindGroupByID(ownerID, groupID uuid.UUID) (group *model.Group, err error) {
	return s.r.FetchByID(ownerID.String(), groupID.String())
}

func (s *GroupService) FindGroups(
	ownerID uuid.UUID,
	pag *types.Pagination,
	needle, sortExpr string) (result *types.Result[model.Group], err error) {
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

func (s *GroupService) UpdateGroup(ownerID, groupID uuid.UUID, up *transfer.GroupUpdate) (ok bool, err error) {
	if len(up.Name) > 50 {
		return false, noda.ErrTooLong.Clone().FormatDetails("name", "group", 50)
	}
	return s.r.Update(ownerID.String(), groupID.String(), up)
}

func (s *GroupService) DeleteGroup(ownerID, groupID uuid.UUID) (ok bool, err error) {
	return s.r.Remove(ownerID.String(), groupID.String())
}
