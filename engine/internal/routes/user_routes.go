package routes

import (
	"noda/api/handler"
	"noda/engine/internal/injector"
	"noda/engine/internal/middleware"

	"github.com/go-chi/chi/v5"
)

func InitializeForUser(r *chi.Mux) {
	s := injector.UserService()
	h := handler.NewUserHandler(s)

	/* For the logged in user.  */

	r.Get("/me", middleware.WithBearerAuthorization(h.GetLoggedInUser))
	r.Patch("/me", middleware.WithBearerAuthorization(h.UpdateLoggedInUser))
	r.Delete("/me", middleware.WithBearerAuthorization(h.RemoveLoggedInUser))
	r.Get("/me/settings", nil)
	r.Post("/me/change_password", nil)

	/* For administrators.  */

	r.Get("/users", middleware.WithBearerAuthorization(h.GetAllUsers))
}
