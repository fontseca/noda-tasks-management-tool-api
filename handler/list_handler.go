package handler

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/http"
	"noda/data/model"
	"noda/data/transfer"
	"noda/data/types"
	"noda/failure"
	"noda/service"
	"strings"
)

type ListHandler struct {
	s service.ListService
}

func NewListHandler(service service.ListService) *ListHandler {
	return &ListHandler{s: service}
}

type listType uint8

const (
	scattered listType = 1
	grouped   listType = 2
)

func (h *ListHandler) doCreateList(t listType, w http.ResponseWriter, r *http.Request) {
	var next = new(transfer.ListCreation)
	var err = parseRequestBody(w, r, next)
	if nil != err {
		failure.EmitError(w, failure.ErrMalformedRequest.Clone().SetDetails(err.Error()))
		return
	}
	err = next.Validate()
	if nil != err {
		failure.EmitError(w, failure.ErrBadRequest.Clone().SetDetails(err.Error()))
		return
	}
	var userID, _ = extractUserPayload(r)
	var insertedID uuid.UUID
	if grouped == t {
		groupID := parseParameterToUUID(w, r, "group_id")
		if didNotParse(groupID) {
			return
		}
		insertedID, err = h.s.Save(userID, groupID, next)
		if gotAndHandledServiceError(w, err) {
			return
		}
	} else {
		insertedID, err = h.s.Save(userID, uuid.Nil, next)
		if gotAndHandledServiceError(w, err) {
			return
		}
	}
	var result = map[string]string{"inserted_id": insertedID.String()}
	data, err := json.Marshal(result)
	if nil != err {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(data)
}

func (h *ListHandler) HandleGroupedListCreation(w http.ResponseWriter, r *http.Request) {
	h.doCreateList(grouped, w, r)
}

func (h *ListHandler) HandleScatteredListCreation(w http.ResponseWriter, r *http.Request) {
	h.doCreateList(scattered, w, r)
}

func (h *ListHandler) doRetrieveListByID(t listType, w http.ResponseWriter, r *http.Request) {
	var (
		userID, _ = extractUserPayload(r)
		groupID   = uuid.Nil
		err       error
	)
	if grouped == t {
		groupID = parseParameterToUUID(w, r, "group_id")
		if didNotParse(groupID) {
			return
		}
	}
	var listID = parseParameterToUUID(w, r, "list_id")
	if didNotParse(listID) {
		return
	}
	list, err := h.s.FetchByID(userID, groupID, listID)
	if gotAndHandledServiceError(w, err) {
		return
	}
	data, err := json.Marshal(list)
	if nil != err {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (h *ListHandler) HandleGroupedListRetrievalByID(w http.ResponseWriter, r *http.Request) {
	h.doRetrieveListByID(grouped, w, r)
}

func (h *ListHandler) HandleScatteredListRetrievalByID(w http.ResponseWriter, r *http.Request) {
	h.doRetrieveListByID(scattered, w, r)
}

func (h *ListHandler) HandleGroupedListsRetrieval(w http.ResponseWriter, r *http.Request) {
	var ownerID, _ = extractUserPayload(r)
	groupID := parseParameterToUUID(w, r, "group_id")
	if didNotParse(groupID) {
		return
	}
	var pagination = parsePagination(w, r)
	if nil == pagination {
		return
	}
	var search, sortExpr = extractQueryParameter(r, "search", ""), extractSorting(w, r)
	if "?" == sortExpr {
		return
	}
	result, err := h.s.FetchGrouped(ownerID, groupID, pagination, search, sortExpr)
	if gotAndHandledServiceError(w, err) {
		return
	}
	data, err := json.Marshal(result)
	if nil != err {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (h *ListHandler) HandleRetrievalOfLists(w http.ResponseWriter, r *http.Request) {
	var ownerID, _ = extractUserPayload(r)
	var pagination = parsePagination(w, r)
	if nil == pagination {
		return
	}
	var search = extractQueryParameter(r, "search", "")
	var sortExpr = extractSorting(w, r)
	if "?" == sortExpr {
		return
	}
	var (
		all    = extractQueryParameter(r, "all", "")
		result *types.Result[model.List]
		err    error
	)
	if 0 == strings.Compare(all, "true") {
		result, err = h.s.Fetch(ownerID, pagination, search, sortExpr)
	} else {
		result, err = h.s.FetchScattered(ownerID, pagination, search, sortExpr)
	}
	if gotAndHandledServiceError(w, err) {
		return
	}
	data, err := json.Marshal(result)
	if nil != err {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (h *ListHandler) doUpdateList(t listType, w http.ResponseWriter, r *http.Request) {
	var (
		ownerID, _ = extractUserPayload(r)
		groupID    = uuid.Nil
		err        error
		target     string
	)
	var listID = parseParameterToUUID(w, r, "list_id")
	if didNotParse(listID) {
		return
	}
	if grouped == t {
		groupID = parseParameterToUUID(w, r, "group_id")
		if didNotParse(groupID) {
			return
		}
		target = fmt.Sprintf("/me/groups/%s/lists/%s", groupID, listID)
	} else {
		target = "/me/lists/" + listID.String()
	}
	var up = new(transfer.ListUpdate)
	err = parseRequestBody(w, r, up)
	if nil != err {
		failure.EmitError(w, failure.ErrMalformedRequest.Clone().SetDetails(err.Error()))
		return
	}
	if "" == up.Name && "" == up.Description {
		redirect(w, r, target)
		return
	}
	ok, err := h.s.Update(ownerID, groupID, listID, up)
	if gotAndHandledServiceError(w, err) {
		return
	}
	if ok {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	redirect(w, r, target)
}

func (h *ListHandler) HandlePartialUpdateOfGroupedList(w http.ResponseWriter, r *http.Request) {
	h.doUpdateList(grouped, w, r)
}

func (h *ListHandler) HandlePartialUpdateOfScatteredList(w http.ResponseWriter, r *http.Request) {
	h.doUpdateList(scattered, w, r)
}

func (h *ListHandler) doDeleteList(t listType, w http.ResponseWriter, r *http.Request) {
	var (
		userID, _ = extractUserPayload(r)
		groupID   = uuid.Nil
	)
	if grouped == t {
		groupID = parseParameterToUUID(w, r, "group_id")
		if didNotParse(groupID) {
			return
		}
	}
	var listID = parseParameterToUUID(w, r, "list_id")
	if didNotParse(listID) {
		return
	}
	err := h.s.Remove(userID, groupID, listID)
	if gotAndHandledServiceError(w, err) {
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *ListHandler) HandleGroupedListDeletion(w http.ResponseWriter, r *http.Request) {
	h.doDeleteList(grouped, w, r)
}

func (h *ListHandler) HandleScatteredListDeletion(w http.ResponseWriter, r *http.Request) {
	h.doDeleteList(scattered, w, r)
}
