package handler

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/http"
	"noda"
	"noda/data/model"
	"noda/data/transfer"
	"noda/data/types"
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
		noda.EmitError(w, noda.ErrMalformedRequest.Clone().SetDetails(err.Error()))
		return
	}
	err = next.Validate()
	if nil != err {
		noda.EmitError(w, noda.ErrBadRequest.Clone().SetDetails(err.Error()))
		return
	}
	var userID, _ = extractUserPayload(r)
	var insertedID uuid.UUID
	if grouped == t {
		groupID := parseParameterToUUID(w, r, "group_id")
		if didNotParse(groupID) {
			return
		}
		insertedID, err = h.s.SaveList(userID, groupID, next)
		if gotAndHandledServiceError(w, err) {
			return
		}
	} else {
		insertedID, err = h.s.SaveList(userID, uuid.Nil, next)
		if gotAndHandledServiceError(w, err) {
			return
		}
	}
	var result = map[string]string{"insertedID": insertedID.String()}
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
	list, err := h.s.FindListByID(userID, groupID, listID)
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
	result, err := h.s.FindGroupedLists(ownerID, groupID, pagination, search, sortExpr)
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
		result, err = h.s.FindLists(ownerID, pagination, search, sortExpr)
	} else {
		result, err = h.s.FindScatteredLists(ownerID, pagination, search, sortExpr)
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
		noda.EmitError(w, noda.ErrMalformedRequest.Clone().SetDetails(err.Error()))
		return
	}
	if "" == up.Name && "" == up.Description {
		redirect(w, r, target)
		return
	}
	ok, err := h.s.UpdateList(ownerID, groupID, listID, up)
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
