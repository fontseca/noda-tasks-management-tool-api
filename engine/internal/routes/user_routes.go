package routes

import (
	"noda/engine/internal/injector"
	"noda/engine/internal/middleware"
	"noda/handler"

	"github.com/go-chi/chi/v5"
)

func InitializeForUsers(router chi.Router) {
	s := injector.UserService()
	h := handler.NewUserHandler(s)

	/* For the logged in user.  */

	router.
		With(middleware.Authorization).
		Group(func(r chi.Router) {
			r.Get("/me", h.RetrieveCurrentUser)
			r.Patch("/me", h.UpdateCurrentUser)
			r.Delete("/me", h.RemoveCurrentUser)
			r.Get("/me/settings", h.RetrieveCurrentUserSettings)
			r.Get("/me/settings/{setting_key}", h.RetrieveOneSettingOfCurrentUser)
			r.Put("/me/settings/{setting_key}", h.UpdateOneSettingForCurrentUser)
			r.Post("/me/change_password", nil)
		})

	/* For administrators.  */

	router.
		With(middleware.Authorization).
		With(middleware.AdminPrivileges).
		Group(func(r chi.Router) {
			r.Get("/users", h.RetrieveAllUsers)
			r.Get("/users/{user_id}", h.RetrieveUserByID)
			r.Get("/users/search", h.SearchUsers)
			r.Delete("/users/{user_id}", h.DeleteUser)
			r.Put("/users/{user_id}/block", h.BlockUser)
			r.Delete("/users/{user_id}/block", h.UnblockUser)
			r.Get("/users/blocked", h.RetrieveAllBlockedUsers)
			r.Put("/users/{user_id}/make_admin", h.PromoteUserToAdmin)
			r.Delete("/users/{user_id}/make_admin", h.DegradeAdminUser)
		})
}
