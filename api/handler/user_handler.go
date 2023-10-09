package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"noda/api/data/transfer"
	"noda/api/data/types"
	"noda/api/service"
	"noda/failure"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type UserHandler struct {
	s *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service}
}

func (h *UserHandler) RetrieveAllUsers(w http.ResponseWriter, r *http.Request) {
	pagination := ParsePagination(w, r)
	if pagination == nil { /* Errors handled in ParsePagination ocurred.  */
		return
	}

	res, err := h.s.GetAll(pagination)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *UserHandler) RetrieveAllBlockedUsers(w http.ResponseWriter, r *http.Request) {
	pagination := ParsePagination(w, r)
	if pagination == nil {
		return
	}

	res, err := h.s.GetAllBlocked(pagination)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *UserHandler) RetrieveUserByID(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.Parse(chi.URLParam(r, "user_id"))
	if err != nil {
		failure.Emit(w, http.StatusBadRequest, "failure with \"user_id\"", err)
		return
	}
	user, err := h.s.GetByID(userID)
	if err != nil {
		switch {
		default:
			w.WriteHeader(http.StatusInternalServerError)
			return
		case errors.Is(err, failure.ErrNotFound):
			failure.Emit(w, http.StatusNotFound,
				"record not found", fmt.Sprintf("could not find any user with ID %q", userID))
			return
		}
	}
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *UserHandler) PromoteUserToAdmin(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.Parse(chi.URLParam(r, "user_id"))
	if err != nil {
		failure.Emit(w, http.StatusBadRequest, "failure with \"user_id\"", err)
		return
	}
	userWasPromoted, err := h.s.PromoteToAdmin(userID)
	if err != nil {
		switch {
		default:
			w.WriteHeader(http.StatusInternalServerError)
			return
		case errors.Is(err, failure.ErrNotFound):
			failure.Emit(w, http.StatusNotFound,
				"record not found", fmt.Sprintf("could not find any user with ID %q", userID))
			return
		}
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
	userID, err := uuid.Parse(chi.URLParam(r, "user_id"))
	if err != nil {
		failure.Emit(w, http.StatusBadRequest, "failure with \"user_id\"", err)
		return
	}
	userWasPromoted, err := h.s.DegradeToNormalUser(userID)
	if err != nil {
		switch {
		default:
			w.WriteHeader(http.StatusInternalServerError)
			return
		case errors.Is(err, failure.ErrNotFound):
			failure.Emit(w, http.StatusNotFound,
				"record not found", fmt.Sprintf("could not find any user with ID %q", userID))
			return
		}
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
	userID, err := uuid.Parse(chi.URLParam(r, "user_id"))
	if err != nil {
		failure.Emit(w, http.StatusBadRequest, "failure with \"user_id\"", err)
		return
	}
	jwtPayload := r.Context().Value(types.ContextKey{}).(types.JWTPayload)
	if jwtPayload.UserID == userID {
		failure.Emit(w, http.StatusBadRequest, "failure with block operation",
			"cannot block yourself")
		return
	}
	userWasBlocked, err := h.s.Block(userID)
	if err != nil {
		switch {
		default:
			w.WriteHeader(http.StatusInternalServerError)
			return
		case errors.Is(err, failure.ErrNotFound):
			failure.Emit(w, http.StatusNotFound,
				"record not found", fmt.Sprintf("could not find any user with ID %q", userID))
			return
		}
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
	userID, err := uuid.Parse(chi.URLParam(r, "user_id"))
	if err != nil {
		failure.Emit(w, http.StatusBadRequest, "failure with \"user_id\"", err)
		return
	}
	jwtPayload := r.Context().Value(types.ContextKey{}).(types.JWTPayload)
	if jwtPayload.UserID == userID {
		failure.Emit(w, http.StatusBadRequest, "failure with unblock operation",
			"cannot unblock yourself")
		return
	}
	userWasUnblocked, err := h.s.Unblock(userID)
	if err != nil {
		switch {
		default:
			w.WriteHeader(http.StatusInternalServerError)
			return
		case errors.Is(err, failure.ErrNotFound):
			failure.Emit(w, http.StatusNotFound,
				"record not found", fmt.Sprintf("could not find any user with ID %q", userID))
			return
		}
	}
	if userWasUnblocked {
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

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.Parse(chi.URLParam(r, "user_id"))
	if err != nil {
		failure.Emit(w, http.StatusBadRequest, "failure with \"user_id\"", err)
		return
	}
	jwtPayload := r.Context().Value(types.ContextKey{}).(types.JWTPayload)
	if jwtPayload.UserID == userID {
		failure.Emit(w, http.StatusBadRequest, "failure with delete operation",
			"cannot perform a self removal")
		return
	}
	if err := h.s.HardDelete(userID); err != nil {
		switch {
		default:
			w.WriteHeader(http.StatusInternalServerError)
			return
		case errors.Is(err, failure.ErrNotFound):
			failure.Emit(w, http.StatusNotFound,
				"record not found", fmt.Sprintf("could not find any user with ID %q", userID))
			return
		}
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) RetrieveCurrentUser(w http.ResponseWriter, r *http.Request) {
	jwtPayload := r.Context().Value(types.ContextKey{}).(types.JWTPayload)
	user, err := h.s.GetByID(jwtPayload.UserID)
	if err != nil {
		switch {
		case errors.Is(err, failure.ErrNotFound):
			failure.Emit(w, http.StatusNotFound, "not found", "this is user account no longer exists")
		}
		return
	}

	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *UserHandler) RetrieveCurrentUserSettings(w http.ResponseWriter, r *http.Request) {
	pagination := ParsePagination(w, r)
	if pagination == nil {
		return
	}
	jwtPayload := r.Context().Value(types.ContextKey{}).(types.JWTPayload)
	settings, err := h.s.GetUserSettings(pagination, jwtPayload.UserID)
	if err != nil {
		switch {
		default:
			w.WriteHeader(http.StatusInternalServerError)
		case errors.Is(err, failure.ErrNotFound):
			failure.Emit(w, http.StatusNotFound, "not found", "this is user account no longer exists")
		}
		return
	}
	if err := json.NewEncoder(w).Encode(settings); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *UserHandler) RetrieveOneSettingOfCurrentUser(w http.ResponseWriter, r *http.Request) {
	settingKey := chi.URLParam(r, "setting_key")
	jwtPayload := r.Context().Value(types.ContextKey{}).(types.JWTPayload)
	setting, err := h.s.GetOneSetting(jwtPayload.UserID, settingKey)
	if err != nil {
		switch {
		default:
			w.WriteHeader(http.StatusInternalServerError)
		case errors.Is(err, failure.ErrSettingNotFound):
			failure.Emit(w, http.StatusNotFound, "not found", fmt.Sprintf("could not find any user setting with key %q", settingKey))
		case errors.Is(err, failure.ErrNotFound):
			failure.Emit(w, http.StatusNotFound, "not found", "this is user account no longer exists")
		}
		return
	}
	if err := json.NewEncoder(w).Encode(setting); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *UserHandler) UpdateOneSettingForCurrentUser(w http.ResponseWriter, r *http.Request) {
	settingKey := chi.URLParam(r, "setting_key")
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	up := &transfer.UserSettingUpdate{}
	if err := decoder.Decode(up); err != nil {
		// TODO: Catch all different errors
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	jwtPayload := r.Context().Value(types.ContextKey{}).(types.JWTPayload)
	wasUpdated, err := h.s.UpdateUserSetting(jwtPayload.UserID, settingKey, up)
	if err != nil {
		switch {
		default:
			w.WriteHeader(http.StatusInternalServerError)
		case errors.Is(err, failure.ErrSettingNotFound):
			failure.Emit(w, http.StatusNotFound, "not found", fmt.Sprintf("could not find any user setting with key %q", settingKey))
		case errors.Is(err, failure.ErrNotFound):
			failure.Emit(w, http.StatusNotFound, "not found", "this is user account no longer exists")
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
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(up); err != nil {
		// TODO: Catch all different errors
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := up.Validate(); err != nil {
		failure.Emit(w, http.StatusBadRequest, "validations did not succeed", err.Dump())
		return
	}

	jwtPayload := r.Context().Value(types.ContextKey{}).(types.JWTPayload)
	userWasUpdated, err := h.s.Update(jwtPayload.UserID, up)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		switch {
		case errors.Is(err, failure.ErrNotFound):
			failure.Emit(w, http.StatusNotFound, "not found", "this is user account no longer exists")
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
	jwtPayload := r.Context().Value(types.ContextKey{}).(types.JWTPayload)
	id, err := h.s.SoftDelete(jwtPayload.UserID)
	if err != nil {
		log.Println(err)
		return
	}

	if err := json.NewEncoder(w).Encode(id); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
