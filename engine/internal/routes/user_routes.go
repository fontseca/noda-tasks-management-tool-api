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
			r.Get("/me", h.GetLoggedInUser)
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
			r.Get("/users", h.GetAllUsers)
			r.Get("/users/{user_id}", h.GetUserByID)
			r.Get("/users/search", nil)
			r.Delete("/users/{user_id}", nil)
			r.Put("/users/{user_id}/block", nil)
			r.Delete("/users/{user_id}/block", nil)
			r.Get("/users/blocked", nil)
			r.Put("/users/{user_id}/make_admin", h.PromoteUserToAdmin)
			r.Delete("/users/{user_id}/make_admin", h.DegradeUserToAdmin)
		})
}
