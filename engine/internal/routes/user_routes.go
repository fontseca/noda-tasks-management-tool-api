package routes

import (
	"noda/api/handler"
	"noda/engine/internal/injector"
	"noda/engine/internal/middleware"

	"github.com/go-chi/chi/v5"
)

func InitializeForUsers(router chi.Router) {
	s := injector.UserService()
	h := handler.NewUserHandler(s)

	/* For the logged in user.  */

	router.
		With(middleware.Authorization).
		Group(func(r chi.Router) {
			r.Get("/me", h.RetrieveLoggedInUser)
			r.Patch("/me", h.UpdateLoggedInUser)
			r.Delete("/me", h.RemoveLoggedInUser)
			r.Get("/me/settings", nil)
			r.Post("/me/change_password", nil)
		})

	/* For administrators.  */

	router.
		With(middleware.Authorization).
		With(middleware.AdminPrivileges).
		Group(func(r chi.Router) {
			r.Get("/users", h.RetrieveAllUsers)
			r.Get("/users/{user_id}", h.RetrieveUserByID)
			r.Get("/users/search", nil)
			r.Delete("/users/{user_id}", h.DeleteUser)
			r.Put("/users/{user_id}/block", h.BlockUser)
			r.Delete("/users/{user_id}/block", h.UnblockUser)
			r.Get("/users/blocked", h.RetrieveAllBlockedUsers)
			r.Put("/users/{user_id}/make_admin", h.PromoteUserToAdmin)
			r.Delete("/users/{user_id}/make_admin", h.DegradeAdminUser)
		})
}
