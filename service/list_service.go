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

type ListService struct {
	r repository.IListRepository
}

func NewListService(r repository.IListRepository) *ListService {
	return &ListService{r}
}

func (s *ListService) SaveList(ownerID, groupID uuid.UUID, next *transfer.ListCreation) (insertedID uuid.UUID, err error) {
	switch {
	case uuid.Nil == ownerID:
		err = noda.NewNilParameterError("SaveList", "ownerID")
		log.Println(err)
		return uuid.Nil, err
	case uuid.Nil == groupID:
		err = noda.NewNilParameterError("SaveList", "groupID")
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
	id, err := s.r.InsertList(ownerID.String(), groupID.String(), next)
	if nil != err {
		return uuid.Nil, err
	}
	return uuid.Parse(id)
}

func (s *ListService) FindListByID(ownerID, groupID, listID uuid.UUID) (list *model.List, err error) {
	switch {
	case uuid.Nil == ownerID:
		err = noda.NewNilParameterError("FindListByID", "ownerID")
		log.Println(err)
		return nil, err
	case uuid.Nil == groupID:
		err = noda.NewNilParameterError("FindListByID", "groupID")
		log.Println(err)
		return nil, err
	case uuid.Nil == listID:
		err = noda.NewNilParameterError("FindListByID", "listID")
		log.Println(err)
		return nil, err
	}
	return s.r.FetchListByID(ownerID.String(), groupID.String(), listID.String())
}

func (s *ListService) GetTodayListID(ownerID uuid.UUID) (listID uuid.UUID, err error) {
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

func (s *ListService) GetTomorrowListID(ownerID uuid.UUID) (listID uuid.UUID, err error) {
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

func (s *ListService) FindLists(
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

func (s *ListService) FindGroupedLists(
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

func (s *ListService) FindScatteredLists(
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
