package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"noda"
	"noda/data/transfer"
	"noda/service"
	"strings"

	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	s *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service}
}

func (h *UserHandler) RetrieveAllUsers(w http.ResponseWriter, r *http.Request) {
	pagination := parsePagination(w, r)
	if pagination == nil { /* Errors handled in parsePagination ocurred.  */
		return
	}
	res, err := h.s.GetAll(pagination)
	if gotAndHandledServiceError(w, err) {
		return
	}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *UserHandler) SearchUsers(w http.ResponseWriter, r *http.Request) {
	pagination := parsePagination(w, r)
	if pagination == nil {
		return
	}
	sortExpr := extractSorting(w, r)
	if strings.Compare(sortExpr, "") == 0 {
		return
	}
	needle := extractQueryParameter(r, "q", "")
	res, err := h.s.SearchUsers(pagination, needle, sortExpr)
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

func (h *UserHandler) RetrieveAllBlockedUsers(w http.ResponseWriter, r *http.Request) {
	pagination := parsePagination(w, r)
	if pagination == nil {
		return
	}

	res, err := h.s.GetAllBlocked(pagination)
	if gotAndHandledServiceError(w, err) {
		return
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *UserHandler) RetrieveUserByID(w http.ResponseWriter, r *http.Request) {
	var userID = parseParameterToUUID(w, r, "user_id")
	if didNotParse(userID) {
		return
	}
	user, err := h.s.GetByID(userID)
	if gotAndHandledServiceError(w, err) {
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

func (h *UserHandler) PromoteUserToAdmin(w http.ResponseWriter, r *http.Request) {
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
	} else {
		var (
			scheme = "http://"
			host   = r.Host
			path   = fmt.Sprintf("/users/%s", userID)
		)
		if r.TLS != nil { /* Running on HTTPS.  */
			scheme = "https://"
		}
		w.Header().Set("Location", fmt.Sprintf("%s%s%s", scheme, host, path))
		w.WriteHeader(http.StatusSeeOther)
	}
}

func (h *UserHandler) DegradeAdminUser(w http.ResponseWriter, r *http.Request) {
	var userID = parseParameterToUUID(w, r, "user_id")
	if didNotParse(userID) {
		return
	}
	userWasPromoted, err := h.s.DegradeToNormalUser(userID)
	if gotAndHandledServiceError(w, err) {
		return
	}
	if userWasPromoted {
		w.WriteHeader(http.StatusNoContent)
	} else {
		var (
			scheme = "http://"
			host   = r.Host
			path   = fmt.Sprintf("/users/%s", userID)
		)
		if r.TLS != nil { /* Running on HTTPS.  */
			scheme = "https://"
		}
		w.Header().Set("Location", fmt.Sprintf("%s%s%s", scheme, host, path))
		w.WriteHeader(http.StatusSeeOther)
	}
}

func (h *UserHandler) BlockUser(w http.ResponseWriter, r *http.Request) {
	var userToBlock = parseParameterToUUID(w, r, "user_id")
	if didNotParse(userToBlock) {
		return
	}
	userID, _ := extractUserPayload(r)
	if userToBlock == userID {
		noda.EmitError(w, noda.ErrSelfOperation)
		return
	}
	userWasBlocked, err := h.s.Block(userToBlock)
	if gotAndHandledServiceError(w, err) {
		return
	}
	if userWasBlocked {
		w.WriteHeader(http.StatusNoContent)
	} else {
		var (
			scheme = "http://"
			host   = r.Host
			path   = fmt.Sprintf("/users/%s", userID)
		)
		if r.TLS != nil { /* Running on HTTPS.  */
			scheme = "https://"
		}
		w.Header().Set("Location", fmt.Sprintf("%s%s%s", scheme, host, path))
		w.WriteHeader(http.StatusSeeOther)
	}
}

func (h *UserHandler) UnblockUser(w http.ResponseWriter, r *http.Request) {
	var userToUnblock = parseParameterToUUID(w, r, "user_id")
	if didNotParse(userToUnblock) {
		return
	}
	userID, _ := extractUserPayload(r)
	if userToUnblock == userID {
		noda.EmitError(w, noda.ErrSelfOperation)
		return
	}
	userWasUnblocked, err := h.s.Unblock(userToUnblock)
	if gotAndHandledServiceError(w, err) {
		return
	}
	if userWasUnblocked {
		w.WriteHeader(http.StatusNoContent)
	} else {
		var (
			scheme = "http://"
			host   = r.Host
			path   = fmt.Sprintf("/users/%s", userToUnblock)
		)
		if r.TLS != nil { /* Running on HTTPS.  */
			scheme = "https://"
		}
		w.Header().Set("Location", fmt.Sprintf("%s%s%s", scheme, host, path))
		w.WriteHeader(http.StatusSeeOther)
	}
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	var userToDelete = parseParameterToUUID(w, r, "user_id")
	if didNotParse(userToDelete) {
		return
	}
	userID, _ := extractUserPayload(r)
	if userToDelete == userID {
		noda.EmitError(w, noda.ErrSelfOperation)
		return
	}
	err := h.s.HardDelete(userToDelete)
	if gotAndHandledServiceError(w, err) {
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) RetrieveCurrentUser(w http.ResponseWriter, r *http.Request) {
	userID, _ := extractUserPayload(r)
	user, err := h.s.GetByID(userID)
	if err != nil {
		var e *noda.Error
		if errors.As(err, &e) {
			switch {
			default:
				noda.EmitError(w, e)
			case errors.Is(err, noda.ErrUserNotFound):
				noda.EmitError(w, noda.ErrUserNoLongerExists)
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

func (h *UserHandler) RetrieveCurrentUserSettings(w http.ResponseWriter, r *http.Request) {
	pagination := parsePagination(w, r)
	if pagination == nil {
		return
	}
	userID, _ := extractUserPayload(r)
	settings, err := h.s.GetUserSettings(pagination, userID)
	if err != nil {
		var e *noda.Error
		if errors.As(err, &e) {
			switch {
			default:
				noda.EmitError(w, e)
			case errors.Is(err, noda.ErrUserNotFound):
				noda.EmitError(w, noda.ErrUserNoLongerExists)
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
	w.Write(data)
}

func (h *UserHandler) RetrieveOneSettingOfCurrentUser(w http.ResponseWriter, r *http.Request) {
	settingKey := chi.URLParam(r, "setting_key")
	userID, _ := extractUserPayload(r)
	setting, err := h.s.GetOneSetting(userID, settingKey)
	if err != nil {
		var e *noda.Error
		if errors.As(err, &e) {
			switch {
			default:
				noda.EmitError(w, e)
			case errors.Is(err, noda.ErrUserNotFound):
				noda.EmitError(w, noda.ErrUserNoLongerExists)
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

func (h *UserHandler) UpdateOneSettingForCurrentUser(w http.ResponseWriter, r *http.Request) {
	up := &transfer.UserSettingUpdate{}
	var err = parseRequestBody(w, r, up)
	if nil != err {
		noda.EmitError(w, noda.ErrMalformedRequest.Clone().SetDetails(err.Error()))
		return
	}
	userID, _ := extractUserPayload(r)
	settingKey := chi.URLParam(r, "setting_key")
	wasUpdated, err := h.s.UpdateUserSetting(userID, settingKey, up)
	if err != nil {
		var e *noda.Error
		if errors.As(err, &e) {
			switch {
			default:
				noda.EmitError(w, e)
			case errors.Is(err, noda.ErrUserNotFound):
				noda.EmitError(w, noda.ErrUserNoLongerExists)
			}
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	if wasUpdated {
		w.WriteHeader(http.StatusNoContent)
	} else {
		var (
			scheme = "http://"
			host   = r.Host
			path   = fmt.Sprintf("/me/settings/%s", settingKey)
		)
		if r.TLS != nil { /* Running on HTTPS.  */
			scheme = "https://"
		}
		w.Header().Set("Location", fmt.Sprintf("%s%s%s", scheme, host, path))
		w.WriteHeader(http.StatusSeeOther)
	}
}

func (h *UserHandler) UpdateCurrentUser(w http.ResponseWriter, r *http.Request) {
	up := &transfer.UserUpdate{}
	var err = parseRequestBody(w, r, up)
	if nil != err {
		noda.EmitError(w, noda.ErrBadRequest.Clone().SetDetails(err.Error()))
		return
	}
	if err = up.Validate(); err != nil {
		noda.EmitError(w, noda.ErrBadRequest.Clone().SetDetails(err.Error()))
		return
	}
	userID, _ := extractUserPayload(r)
	userWasUpdated, err := h.s.Update(userID, up)
	if err != nil {
		var e *noda.Error
		if errors.As(err, &e) {
			switch {
			default:
				noda.EmitError(w, e)
			case errors.Is(err, noda.ErrUserNotFound):
				noda.EmitError(w, noda.ErrUserNoLongerExists)
			}
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	if userWasUpdated {
		w.WriteHeader(http.StatusNoContent)
	} else {
		var (
			scheme = "http://"
			host   = r.Host
			self   = r.URL.Path
		)
		if r.TLS != nil { /* Running on HTTPS.  */
			scheme = "https://"
		}
		w.Header().Set("Location", fmt.Sprintf("%s%s%s", scheme, host, self))
		w.WriteHeader(http.StatusSeeOther)
	}
}

func (h *UserHandler) RemoveCurrentUser(w http.ResponseWriter, r *http.Request) {
	userID, _ := extractUserPayload(r)
	id, err := h.s.SoftDelete(userID)
	if gotAndHandledServiceError(w, err) {
		return
	}
	data, err := json.Marshal(id)
	if nil != err {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}
