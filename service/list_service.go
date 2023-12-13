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
)

type ListService interface {
	Save(ownerID, groupID uuid.UUID, creation *transfer.ListCreation) (insertedID uuid.UUID, err error)
	GetTodayListID(ownerID uuid.UUID) (listID uuid.UUID, err error)
	GetTomorrowListID(ownerID uuid.UUID) (listID uuid.UUID, err error)
	FetchByID(ownerID, groupID, listID uuid.UUID) (list *model.List, err error)
	Fetch(ownerID uuid.UUID, pagination *types.Pagination, needle, sortExpr string) (result *types.Result[model.List], err error)
	FetchGrouped(ownerID, groupID uuid.UUID, pagination *types.Pagination, needle, sortExpr string) (result *types.Result[model.List], err error)
	FetchScattered(ownerID uuid.UUID, pagination *types.Pagination, needle, sortExpr string) (result *types.Result[model.List], err error)
	Update(ownerID, groupID, listID uuid.UUID, update *transfer.ListUpdate) (ok bool, err error)
	Duplicate(ownerID, listID uuid.UUID) (replicaID uuid.UUID, err error)
	Move(ownerID, listID, targetGroupID uuid.UUID) (ok bool, err error)
	Scatter(ownerID, listID uuid.UUID) (ok bool, err error)
	Remove(ownerID, groupID, listID uuid.UUID) error
}

type listService struct {
	r repository.ListRepository
}

func NewListService(r repository.ListRepository) ListService {
	return &listService{r}
}

func (s *listService) Save(ownerID, groupID uuid.UUID, creation *transfer.ListCreation) (insertedID uuid.UUID, err error) {
	var groupIDStr = ""
	switch {
	case uuid.Nil == ownerID:
		err = noda.NewNilParameterError("Save", "ownerID")
		log.Println(err)
		return uuid.Nil, err
	case nil == creation:
		err = noda.NewNilParameterError("Save", "creation")
		log.Println(err)
		return uuid.Nil, err
	}
	doTrim(&creation.Name, &creation.Description)
	switch {
	case "" == creation.Name:
		return uuid.Nil, errors.New("name cannot be an empty string") // must've been handled by validator
	case 1<<5 < len(creation.Name):
		return uuid.Nil, noda.ErrTooLong.Clone().FormatDetails("name", "list", 1<<5)
	case 1<<9 < len(creation.Description):
		return uuid.Nil, noda.ErrTooLong.Clone().FormatDetails("description", "list", 1<<9)
	}
	if uuid.Nil != groupID {
		groupIDStr = groupID.String()
	}
	id, err := s.r.Save(ownerID.String(), groupIDStr, creation)
	if nil != err {
		return uuid.Nil, err
	}
	return uuid.Parse(id)
}

func (s *listService) FetchByID(ownerID, groupID, listID uuid.UUID) (list *model.List, err error) {
	var groupIDStr = ""
	switch {
	case uuid.Nil == ownerID:
		err = noda.NewNilParameterError("FetchByID", "ownerID")
		log.Println(err)
		return nil, err
	case uuid.Nil == listID:
		err = noda.NewNilParameterError("FetchByID", "listID")
		log.Println(err)
		return nil, err
	}
	if uuid.Nil != groupID {
		groupIDStr = groupID.String()
	}
	return s.r.FetchByID(ownerID.String(), groupIDStr, listID.String())
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

func (s *listService) Fetch(
	ownerID uuid.UUID,
	pagination *types.Pagination,
	needle, sortExpr string,
) (lists *types.Result[model.List], err error) {
	switch {
	case uuid.Nil == ownerID:
		err = noda.NewNilParameterError("Fetch", "ownerID")
		log.Println(err)
		return nil, err
	case nil == pagination:
		err = noda.NewNilParameterError("Fetch", "pagination")
		log.Println(err)
		return nil, err
	}
	res, err := s.r.Fetch(ownerID.String(), pagination.Page, pagination.RPP, needle, sortExpr)
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

func (s *listService) FetchGrouped(
	ownerID, groupID uuid.UUID,
	pagination *types.Pagination,
	needle, sortExpr string,
) (result *types.Result[model.List], err error) {
	switch {
	case uuid.Nil == ownerID:
		err = noda.NewNilParameterError("FetchGrouped", "ownerID")
		log.Println(err)
		return nil, err
	case uuid.Nil == groupID:
		err = noda.NewNilParameterError("FetchGrouped", "groupID")
		log.Println(err)
		return nil, err
	case nil == pagination:
		err = noda.NewNilParameterError("FetchGrouped", "pagination")
		log.Println(err)
		return nil, err
	}
	doTrim(&needle, &sortExpr)
	doDefaultPagination(pagination)
	res, err := s.r.FetchGrouped(ownerID.String(), groupID.String(), pagination.Page, pagination.RPP, needle, sortExpr)
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

func (s *listService) FetchScattered(
	ownerID uuid.UUID,
	pagination *types.Pagination,
	needle, sortExpr string,
) (result *types.Result[model.List], err error) {
	switch {
	case uuid.Nil == ownerID:
		err = noda.NewNilParameterError("FetchScattered", "ownerID")
		log.Println(err)
		return nil, err
	case nil == pagination:
		err = noda.NewNilParameterError("FetchScattered", "pagination")
		log.Println(err)
		return nil, err
	}
	doTrim(&needle, &sortExpr)
	doDefaultPagination(pagination)
	res, err := s.r.FetchScattered(ownerID.String(), pagination.Page, pagination.RPP, needle, sortExpr)
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

func (s *listService) Remove(ownerID, groupID, listID uuid.UUID) error {
	var (
		err        error
		groupIDStr = ""
	)
	switch {
	case uuid.Nil == ownerID:
		err = noda.NewNilParameterError("Remove", "ownerID")
		log.Println(err)
		return err
	case uuid.Nil == listID:
		err = noda.NewNilParameterError("Remove", "listID")
		log.Println(err)
		return err
	case uuid.Nil != groupID:
		groupIDStr = groupID.String()
	}
	_, err = s.r.Remove(ownerID.String(), groupIDStr, listID.String())
	return err
}

func (s *listService) Duplicate(ownerID, listID uuid.UUID) (replicaID uuid.UUID, err error) {
	switch {
	case uuid.Nil == ownerID:
		err = noda.NewNilParameterError("Duplicate", "ownerID")
		log.Println(err)
		return uuid.Nil, err
	case uuid.Nil == listID:
		err = noda.NewNilParameterError("Duplicate", "listID")
		log.Println(err)
		return uuid.Nil, err
	}
	id, err := s.r.Duplicate(ownerID.String(), listID.String())
	if nil != err {
		return uuid.Nil, err
	}
	return uuid.Parse(id)
}

func (s *listService) Scatter(ownerID, listID uuid.UUID) (ok bool, err error) {
	switch {
	case uuid.Nil == ownerID:
		err = noda.NewNilParameterError("Scatter", "ownerID")
		log.Println(err)
		return false, err
	case uuid.Nil == listID:
		err = noda.NewNilParameterError("Scatter", "listID")
		log.Println(err)
		return false, err
	}
	return s.r.Scatter(ownerID.String(), listID.String())
}

func (s *listService) Move(ownerID, listID, targetGroupID uuid.UUID) (ok bool, err error) {
	switch {
	case uuid.Nil == ownerID:
		err = noda.NewNilParameterError("Move", "ownerID")
		log.Println(err)
		return false, err
	case uuid.Nil == listID:
		err = noda.NewNilParameterError("Move", "listID")
		log.Println(err)
		return false, err
	case uuid.Nil == targetGroupID:
		err = noda.NewNilParameterError("Move", "targetGroupID")
		log.Println(err)
		return false, err
	}
	return s.r.Move(ownerID.String(), listID.String(), targetGroupID.String())
}

func (s *listService) Update(ownerID, groupID, listID uuid.UUID, up *transfer.ListUpdate) (ok bool, err error) {
	var groupIDStr = ""
	switch {
	case uuid.Nil == ownerID:
		err = noda.NewNilParameterError("Update", "ownerID")
		log.Println(err)
		return false, err
	case uuid.Nil == listID:
		err = noda.NewNilParameterError("Update", "listID")
		log.Println(err)
		return false, err
	case nil == up:
		err = noda.NewNilParameterError("Update", "up")
		log.Println(err)
		return false, err
	case uuid.Nil != groupID:
		groupIDStr = groupID.String()
	}
	doTrim(&up.Name, &up.Description)
	switch {
	case 1<<5 < len(up.Name):
		return false, noda.ErrTooLong.Clone().FormatDetails("name", "list", 1<<5)
	case 1<<9 < len(up.Description):
		return false, noda.ErrTooLong.Clone().FormatDetails("description", "list", 1<<9)
	}
	return s.r.Update(ownerID.String(), groupIDStr, listID.String(), up)
}
