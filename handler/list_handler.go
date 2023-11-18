package handler

import (
	"encoding/json"
	"errors"
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
	var err = decodeJSONRequestBody(w, r, next)
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
		groupID, err := parsePathParameterToUUID(r, "group_id")
		if nil != err {
			var e *noda.Error
			if errors.As(err, &e) {
				noda.EmitError(w, e)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}
		insertedID, err = h.s.SaveList(userID, groupID, next)
		if nil != err {
			var e *noda.Error
			if errors.As(err, &e) {
				noda.EmitError(w, e)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}
	} else {
		insertedID, err = h.s.SaveList(userID, uuid.Nil, next)
		if nil != err {
			var e *noda.Error
			if errors.As(err, &e) {
				noda.EmitError(w, e)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
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
		groupID, err = parsePathParameterToUUID(r, "group_id")
		if nil != err {
			var e *noda.Error
			if errors.As(err, &e) {
				noda.EmitError(w, e)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}
	}
	listID, err := parsePathParameterToUUID(r, "list_id")
	if nil != err {
		var e *noda.Error
		if errors.As(err, &e) {
			noda.EmitError(w, e)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	list, err := h.s.FindListByID(userID, groupID, listID)
	if nil != err {
		var e *noda.Error
		if errors.As(err, &e) {
			noda.EmitError(w, e)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
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
	groupID, err := parsePathParameterToUUID(r, "group_id")
	if nil != err {
		var e *noda.Error
		if errors.As(err, &e) {
			noda.EmitError(w, e)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	var pagination = parsePagination(w, r)
	if nil == pagination {
		return
	}
	var search, sortExpr = parseQueryParameter(r, "search", ""), parseSorting(w, r)
	if "?" == sortExpr {
		return
	}
	result, err := h.s.FindGroupedLists(ownerID, groupID, pagination, search, sortExpr)
	if nil != err {
		var e *noda.Error
		if errors.As(err, &e) {
			noda.EmitError(w, e)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
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
	var search = parseQueryParameter(r, "search", "")
	var sortExpr = parseSorting(w, r)
	if "?" == sortExpr {
		return
	}
	var (
		all    = parseQueryParameter(r, "all", "")
		result *types.Result[model.List]
		err    error
	)
	if 0 == strings.Compare(all, "true") {
		result, err = h.s.FindLists(ownerID, pagination, search, sortExpr)
	} else {
		result, err = h.s.FindScatteredLists(ownerID, pagination, search, sortExpr)
	}
	if nil != err {
		var e *noda.Error
		if errors.As(err, &e) {
			noda.EmitError(w, e)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
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
	listID, err := parsePathParameterToUUID(r, "list_id")
	if nil != err {
		var e *noda.Error
		if errors.As(err, &e) {
			noda.EmitError(w, e)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	if grouped == t {
		groupID, err = parsePathParameterToUUID(r, "group_id")
		if nil != err {
			var e *noda.Error
			if errors.As(err, &e) {
				noda.EmitError(w, e)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}
		target = fmt.Sprintf("/me/groups/%s/lists/%s", groupID, listID)
	} else {
		target = fmt.Sprintf("/me/lists/%s", listID)
	}
	var up = new(transfer.ListUpdate)
	err = decodeJSONRequestBody(w, r, up)
	if nil != err {
		noda.EmitError(w, noda.ErrMalformedRequest.Clone().SetDetails(err.Error()))
		return
	}
	if "" == up.Name && "" == up.Description {
		redirect(w, r, target)
		return
	}
	ok, err := h.s.UpdateList(ownerID, groupID, listID, up)
	if nil != err {
		var e *noda.Error
		if errors.As(err, &e) {
			noda.EmitError(w, e)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
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
