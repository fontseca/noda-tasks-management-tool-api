package service

import (
	"errors"
	"github.com/google/uuid"
	"log"
	"noda"
	"noda/data/model"
	"noda/data/transfer"
	"noda/data/types"
	"noda/repository"
	"strings"
)

type ListService interface {
	SaveList(ownerID, groupID uuid.UUID, next *transfer.ListCreation) (insertedID uuid.UUID, err error)
	FindListByID(ownerID, groupID, listID uuid.UUID) (list *model.List, err error)
	GetTodayListID(ownerID uuid.UUID) (listID uuid.UUID, err error)
	GetTomorrowListID(ownerID uuid.UUID) (listID uuid.UUID, err error)
	FindLists(ownerID uuid.UUID, pagination *types.Pagination, needle, sortBy string) (lists *types.Result[model.List], err error)
	FindGroupedLists(ownerID, groupID uuid.UUID, pagination *types.Pagination, needle, sortBy string) (result *types.Result[model.List], err error)
	FindScatteredLists(ownerID uuid.UUID, pagination *types.Pagination, needle, sortBy string) (result *types.Result[model.List], err error)
	DeleteList(ownerID, groupID, listID uuid.UUID) error
	DuplicateList(ownerID, listID uuid.UUID) (replicaID uuid.UUID, err error)
	ConvertToScatteredList(ownerID, listID uuid.UUID) (ok bool, err error)
	MoveList(ownerID, listID, targetGroupID uuid.UUID) (ok bool, err error)
	UpdateList(ownerID, groupID, listID uuid.UUID, up *transfer.ListUpdate) (ok bool, err error)
}

type listService struct {
	r repository.IListRepository
}

func NewListService(r repository.IListRepository) ListService {
	return &listService{r}
}

func (s *listService) SaveList(ownerID, groupID uuid.UUID, next *transfer.ListCreation) (insertedID uuid.UUID, err error) {
	var groupIDStr = ""
	switch {
	case uuid.Nil == ownerID:
		err = noda.NewNilParameterError("SaveList", "ownerID")
		log.Println(err)
		return uuid.Nil, err
	case nil == next:
		err = noda.NewNilParameterError("SaveList", "next")
		log.Println(err)
		return uuid.Nil, err
	}
	next.Name = strings.Trim(next.Name, " \t\n")
	next.Description = strings.Trim(next.Description, " \t\n")
	switch {
	case "" == next.Name:
		return uuid.Nil, errors.New("name cannot be an empty string")
	case 50 < len(next.Name):
		return uuid.Nil, errors.New("name too long for list: max length must be 50")
	}
	if uuid.Nil != groupID {
		groupIDStr = groupID.String()
	}
	id, err := s.r.InsertList(ownerID.String(), groupIDStr, next)
	if nil != err {
		return uuid.Nil, err
	}
	return uuid.Parse(id)
}

func (s *listService) FindListByID(ownerID, groupID, listID uuid.UUID) (list *model.List, err error) {
	var groupIDStr = ""
	switch {
	case uuid.Nil == ownerID:
		err = noda.NewNilParameterError("FindListByID", "ownerID")
		log.Println(err)
		return nil, err
	case uuid.Nil == listID:
		err = noda.NewNilParameterError("FindListByID", "listID")
		log.Println(err)
		return nil, err
	}
	if uuid.Nil != groupID {
		groupIDStr = groupID.String()
	}
	return s.r.FetchListByID(ownerID.String(), groupIDStr, listID.String())
}

func (s *listService) GetTodayListID(ownerID uuid.UUID) (listID uuid.UUID, err error) {
	if uuid.Nil == ownerID {
		err = noda.NewNilParameterError("GetTodayListID", "ownerID")
		log.Println(err)
		return uuid.Nil, err
	}
	id, err := s.r.GetTodayListID(ownerID.String())
	if nil != err {
		return uuid.Nil, err
	}
	return uuid.Parse(id)
}

func (s *listService) GetTomorrowListID(ownerID uuid.UUID) (listID uuid.UUID, err error) {
	if uuid.Nil == ownerID {
		err = noda.NewNilParameterError("GetTomorrowListID", "ownerID")
		log.Println(err)
		return uuid.Nil, err
	}
	id, err := s.r.GetTomorrowListID(ownerID.String())
	if nil != err {
		return uuid.Nil, err
	}
	return uuid.Parse(id)
}

func (s *listService) FindLists(
	ownerID uuid.UUID, pagination *types.Pagination, needle, sortBy string) (lists *types.Result[model.List], err error) {
	switch {
	case uuid.Nil == ownerID:
		err = noda.NewNilParameterError("FindLists", "ownerID")
		log.Println(err)
		return nil, err
	case nil == pagination:
		err = noda.NewNilParameterError("FindLists", "pagination")
		log.Println(err)
		return nil, err
	}
	res, err := s.r.FetchLists(ownerID.String(), pagination.Page, pagination.RPP, needle, sortBy)
	if nil != err {
		return nil, err
	}
	lists = &types.Result[model.List]{
		Page:      pagination.Page,
		RPP:       pagination.RPP,
		Retrieved: int64(len(res)),
		Payload:   res,
	}
	return lists, nil
}

func (s *listService) FindGroupedLists(
	ownerID, groupID uuid.UUID,
	pagination *types.Pagination, needle, sortBy string) (result *types.Result[model.List], err error) {
	switch {
	case uuid.Nil == ownerID:
		err = noda.NewNilParameterError("FindGroupedLists", "ownerID")
		log.Println(err)
		return nil, err
	case uuid.Nil == groupID:
		err = noda.NewNilParameterError("FindGroupedLists", "groupID")
		log.Println(err)
		return nil, err
	case nil == pagination:
		err = noda.NewNilParameterError("FindGroupedLists", "pagination")
		log.Println(err)
		return nil, err
	}
	setToDefaultValues(pagination, &needle, &sortBy)
	res, err := s.r.FetchGroupedLists(ownerID.String(), groupID.String(), pagination.Page, pagination.RPP, needle, sortBy)
	if nil != err {
		return nil, err
	}
	result = &types.Result[model.List]{
		Page:      pagination.Page,
		RPP:       pagination.RPP,
		Retrieved: int64(len(res)),
		Payload:   res,
	}
	return result, nil
}

func (s *listService) FindScatteredLists(
	ownerID uuid.UUID, pagination *types.Pagination, needle, sortBy string) (result *types.Result[model.List], err error) {
	switch {
	case uuid.Nil == ownerID:
		err = noda.NewNilParameterError("FindScatteredLists", "ownerID")
		log.Println(err)
		return nil, err
	case nil == pagination:
		err = noda.NewNilParameterError("FindScatteredLists", "pagination")
		log.Println(err)
		return nil, err
	}
	setToDefaultValues(pagination, &needle, &sortBy)
	res, err := s.r.FetchScatteredLists(ownerID.String(), pagination.Page, pagination.RPP, needle, sortBy)
	if nil != err {
		return nil, err
	}
	result = &types.Result[model.List]{
		Page:      pagination.Page,
		RPP:       pagination.RPP,
		Retrieved: int64(len(res)),
		Payload:   res,
	}
	return result, nil
}

func setToDefaultValues(pagination *types.Pagination, needle, sortBy *string) {
	switch {
	case "" != *needle:
		*needle = strings.Trim(*needle, " \n\t")
	case "" != *sortBy:
		*sortBy = strings.Trim(*sortBy, " \n\t")
	case 0 >= pagination.Page:
		pagination.Page = 1
	case 0 >= pagination.RPP:
		pagination.RPP = 10
	}
}

func (s *listService) DeleteList(ownerID, groupID, listID uuid.UUID) error {
	var (
		err        error
		groupIDStr = ""
	)
	switch {
	case uuid.Nil == ownerID:
		err = noda.NewNilParameterError("DeleteList", "ownerID")
		log.Println(err)
		return err
	case uuid.Nil == listID:
		err = noda.NewNilParameterError("DeleteList", "listID")
		log.Println(err)
		return err
	case uuid.Nil != groupID:
		groupIDStr = groupID.String()
	}
	_, err = s.r.DeleteList(ownerID.String(), groupIDStr, listID.String())
	return err
}

func (s *listService) DuplicateList(ownerID, listID uuid.UUID) (replicaID uuid.UUID, err error) {
	switch {
	case uuid.Nil == ownerID:
		err = noda.NewNilParameterError("DuplicateList", "ownerID")
		log.Println(err)
		return uuid.Nil, err
	case uuid.Nil == listID:
		err = noda.NewNilParameterError("DuplicateList", "listID")
		log.Println(err)
		return uuid.Nil, err
	}
	id, err := s.r.DuplicateList(ownerID.String(), listID.String())
	if nil != err {
		return uuid.Nil, err
	}
	return uuid.Parse(id)
}

func (s *listService) ConvertToScatteredList(ownerID, listID uuid.UUID) (ok bool, err error) {
	switch {
	case uuid.Nil == ownerID:
		err = noda.NewNilParameterError("ConvertToScatteredList", "ownerID")
		log.Println(err)
		return false, err
	case uuid.Nil == listID:
		err = noda.NewNilParameterError("ConvertToScatteredList", "listID")
		log.Println(err)
		return false, err
	}
	return s.r.ConvertToScatteredList(ownerID.String(), listID.String())
}

func (s *listService) MoveList(ownerID, listID, targetGroupID uuid.UUID) (ok bool, err error) {
	switch {
	case uuid.Nil == ownerID:
		err = noda.NewNilParameterError("MoveList", "ownerID")
		log.Println(err)
		return false, err
	case uuid.Nil == listID:
		err = noda.NewNilParameterError("MoveList", "listID")
		log.Println(err)
		return false, err
	case uuid.Nil == targetGroupID:
		err = noda.NewNilParameterError("MoveList", "targetGroupID")
		log.Println(err)
		return false, err
	}
	return s.r.MoveList(ownerID.String(), listID.String(), targetGroupID.String())
}

func (s *listService) UpdateList(ownerID, groupID, listID uuid.UUID, up *transfer.ListUpdate) (ok bool, err error) {
	var groupIDStr = ""
	switch {
	case uuid.Nil == ownerID:
		err = noda.NewNilParameterError("UpdateList", "ownerID")
		log.Println(err)
		return false, err
	case uuid.Nil == listID:
		err = noda.NewNilParameterError("UpdateList", "listID")
		log.Println(err)
		return false, err
	case nil == up:
		err = noda.NewNilParameterError("UpdateList", "up")
		log.Println(err)
		return false, err
	case uuid.Nil != groupID:
		groupIDStr = groupID.String()
	}
	up.Name = strings.Trim(up.Name, " \t\n")
	up.Description = strings.Trim(up.Description, " \t\n")
	switch {
	case 50 < len(up.Name):
		return false, noda.ErrTooLong.Clone().FormatDetails("name", "list", 50)
	case 1<<9 < len(up.Description):
		return false, noda.ErrTooLong.Clone().FormatDetails("description", "list", 1<<9)
	}
	return s.r.UpdateList(ownerID.String(), groupIDStr, listID.String(), up)
}
