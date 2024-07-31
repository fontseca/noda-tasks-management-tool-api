package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"noda/data/transfer"
	"noda/failure"
	"noda/service"
	"strings"
)

type UserHandler struct {
	s service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{service}
}

func (h *UserHandler) HandleUsersRetrieval(w http.ResponseWriter, r *http.Request) {
	pagination := parsePagination(w, r)
	if pagination == nil {
		return
	}
	var sortExpr = extractSorting(w, r)
	var needle = extractQueryParameter(r, "search", "")
	res, err := h.s.Fetch(pagination, needle, sortExpr)
	if gotAndHandledServiceError(w, err) {
		return
	}
	data, err := json.Marshal(res)
	if nil != err {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (h *UserHandler) HandleUsersSearch(w http.ResponseWriter, r *http.Request) {
	pagination := parsePagination(w, r)
	if pagination == nil {
		return
	}
	sortExpr := extractSorting(w, r)
	if strings.Compare(sortExpr, "") == 0 {
		return
	}
	needle := extractQueryParameter(r, "q", "")
	res, err := h.s.Search(pagination, needle, sortExpr)
	if gotAndHandledServiceError(w, err) {
		return
	}
	data, err := json.Marshal(res)
	if nil != err {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func (h *UserHandler) HandleBlockedUsersRetrieval(w http.ResponseWriter, r *http.Request) {
	pagination := parsePagination(w, r)
	if pagination == nil {
		return
	}
	var sortExpr = extractSorting(w, r)
	var needle = extractQueryParameter(r, "search", "")
	res, err := h.s.FetchBlocked(pagination, needle, sortExpr)
	if gotAndHandledServiceError(w, err) {
		return
	}
	data, err := json.Marshal(res)
	if nil != err {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (h *UserHandler) HandleRetrievalOfUserByID(w http.ResponseWriter, r *http.Request) {
	var userID = parseParameterToUUID(w, r, "user_id")
	if didNotParse(userID) {
		return
	}
	user, err := h.s.FetchByID(userID)
	if gotAndHandledServiceError(w, err) {
		return
	}
	data, err := json.Marshal(user)
	if nil != err {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (h *UserHandler) HandleAdminPromotion(w http.ResponseWriter, r *http.Request) {
	var userID = parseParameterToUUID(w, r, "user_id")
	if didNotParse(userID) {
		return
	}
	userWasPromoted, err := h.s.PromoteToAdmin(userID)
	if gotAndHandledServiceError(w, err) {
		return
	}
	if userWasPromoted {
		w.WriteHeader(http.StatusNoContent)
	}
	redirect(w, r, "/users/"+userID.String())
}

func (h *UserHandler) HandleDegradeAdminToUser(w http.ResponseWriter, r *http.Request) {
	var userID = parseParameterToUUID(w, r, "user_id")
	if didNotParse(userID) {
		return
	}
	userWasPromoted, err := h.s.DegradeToUser(userID)
	if gotAndHandledServiceError(w, err) {
		return
	}
	if userWasPromoted {
		w.WriteHeader(http.StatusNoContent)
	}
	redirect(w, r, "/users/"+userID.String())
}

func (h *UserHandler) HandleBlockUser(w http.ResponseWriter, r *http.Request) {
	var userToBlock = parseParameterToUUID(w, r, "user_id")
	if didNotParse(userToBlock) {
		return
	}
	userID, _ := extractUserPayload(r)
	if userToBlock == userID {
		failure.EmitError(w, failure.ErrSelfOperation)
		return
	}
	userWasBlocked, err := h.s.Block(userToBlock)
	if gotAndHandledServiceError(w, err) {
		return
	}
	if userWasBlocked {
		w.WriteHeader(http.StatusNoContent)
	}
	redirect(w, r, "/users/"+userID.String())
}

func (h *UserHandler) HandleUnblockUser(w http.ResponseWriter, r *http.Request) {
	var userToUnblock = parseParameterToUUID(w, r, "user_id")
	if didNotParse(userToUnblock) {
		return
	}
	userID, _ := extractUserPayload(r)
	if userToUnblock == userID {
		failure.EmitError(w, failure.ErrSelfOperation)
		return
	}
	userWasUnblocked, err := h.s.Unblock(userToUnblock)
	if gotAndHandledServiceError(w, err) {
		return
	}
	if userWasUnblocked {
		w.WriteHeader(http.StatusNoContent)
	}
	redirect(w, r, "/users/"+userToUnblock.String())
}

func (h *UserHandler) HandleUserDeletion(w http.ResponseWriter, r *http.Request) {
	var userToDelete = parseParameterToUUID(w, r, "user_id")
	if didNotParse(userToDelete) {
		return
	}
	userID, _ := extractUserPayload(r)
	if userToDelete == userID {
		failure.EmitError(w, failure.ErrSelfOperation)
		return
	}
	err := h.s.RemoveHardly(userToDelete)
	if gotAndHandledServiceError(w, err) {
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) HandleRetrievalOfLoggedInUser(w http.ResponseWriter, r *http.Request) {
	userID, _ := extractUserPayload(r)
	user, err := h.s.FetchByID(userID)
	if err != nil {
		var e *failure.Error
		if errors.As(err, &e) {
			switch {
			default:
				failure.EmitError(w, e)
			case errors.Is(err, failure.ErrUserNotFound):
				failure.EmitError(w, failure.ErrUserNoLongerExists)
			}
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	data, err := json.Marshal(user)
	if nil != err {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func (h *UserHandler) HandleRetrievalOfLoggedUserSettings(w http.ResponseWriter, r *http.Request) {
	pagination := parsePagination(w, r)
	if pagination == nil {
		return
	}
	userID, _ := extractUserPayload(r)
	var sortExpr = extractSorting(w, r)
	var needle = extractQueryParameter(r, "search", "")
	settings, err := h.s.FetchSettings(userID, pagination, needle, sortExpr)
	if err != nil {
		var e *failure.Error
		if errors.As(err, &e) {
			switch {
			default:
				failure.EmitError(w, e)
			case errors.Is(err, failure.ErrUserNotFound):
				failure.EmitError(w, failure.ErrUserNoLongerExists)
			}
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	data, err := json.Marshal(settings)
	if nil != err {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (h *UserHandler) HandleRetrievalOfOneSettingOfLoggedUser(w http.ResponseWriter, r *http.Request) {
	settingKey := r.PathValue("setting_key")
	userID, _ := extractUserPayload(r)
	setting, err := h.s.FetchOneSetting(userID, settingKey)
	if err != nil {
		var e *failure.Error
		if errors.As(err, &e) {
			switch {
			default:
				failure.EmitError(w, e)
			case errors.Is(err, failure.ErrUserNotFound):
				failure.EmitError(w, failure.ErrUserNoLongerExists)
			}
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	data, err := json.Marshal(setting)
	if nil != err {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func (h *UserHandler) HandleUpdateOneSettingForLoggedUser(w http.ResponseWriter, r *http.Request) {
	up := &transfer.UserSettingUpdate{}
	var err = parseRequestBody(w, r, up)
	if nil != err {
		failure.EmitError(w, failure.ErrMalformedRequest.Clone().SetDetails(err.Error()))
		return
	}
	userID, _ := extractUserPayload(r)
	settingKey := r.PathValue("setting_key")
	wasUpdated, err := h.s.UpdateUserSetting(userID, settingKey, up)
	if err != nil {
		var e *failure.Error
		if errors.As(err, &e) {
			switch {
			default:
				failure.EmitError(w, e)
			case errors.Is(err, failure.ErrUserNotFound):
				failure.EmitError(w, failure.ErrUserNoLongerExists)
			}
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	if wasUpdated {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	redirect(w, r, "/me/settings/"+settingKey)
}

func (h *UserHandler) HandleUpdateForLoggedUser(w http.ResponseWriter, r *http.Request) {
	up := &transfer.UserUpdate{}
	var err = parseRequestBody(w, r, up)
	if nil != err {
		failure.EmitError(w, failure.ErrBadRequest.Clone().SetDetails(err.Error()))
		return
	}
	if err = up.Validate(); err != nil {
		failure.EmitError(w, failure.ErrBadRequest.Clone().SetDetails(err.Error()))
		return
	}
	userID, _ := extractUserPayload(r)
	userWasUpdated, err := h.s.Update(userID, up)
	if err != nil {
		var e *failure.Error
		if errors.As(err, &e) {
			switch {
			default:
				failure.EmitError(w, e)
			case errors.Is(err, failure.ErrUserNotFound):
				failure.EmitError(w, failure.ErrUserNoLongerExists)
			}
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	if userWasUpdated {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	redirect(w, r, r.URL.Path)
}

func (h *UserHandler) HandleRemovalOfLoggedUser(w http.ResponseWriter, r *http.Request) {
	userID, _ := extractUserPayload(r)
	err := h.s.RemoveSoftly(userID)
	if gotAndHandledServiceError(w, err) {
		return
	}
}
