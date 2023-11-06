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
	r repository.IGroupRepository
}

func NewGroupService(repository repository.IGroupRepository) *GroupService {
	return &GroupService{repository}
}

func (s *GroupService) SaveGroup(ownerID uuid.UUID, newGroup *transfer.GroupCreation) (insertedID string, err error) {
	if len(newGroup.Name) > 50 {
		return "", noda.ErrTooLong.Clone().FormatDetails("name", "group", 50)
	}
	return s.r.InsertGroup(ownerID.String(), newGroup)
}

func (s *GroupService) FindGroupByID(ownerID, groupID uuid.UUID) (group *model.Group, err error) {
	return s.r.FetchGroupByID(ownerID.String(), groupID.String())
}

func (s *GroupService) FindGroups(
	ownerID uuid.UUID,
	pag *types.Pagination,
	needle, sortExpr string) (result *types.Result[model.Group], err error) {
	groups, err := s.r.FetchGroups(ownerID.String(), pag.Page, pag.RPP, needle, sortExpr)
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
	return s.r.UpdateGroup(ownerID.String(), groupID.String(), up)
}

func (s *GroupService) DeleteGroup(ownerID, groupID uuid.UUID) (ok bool, err error) {
	return s.r.DeleteGroup(ownerID.String(), groupID.String())
}
