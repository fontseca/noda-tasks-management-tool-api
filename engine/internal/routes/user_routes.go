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

	r.Get("/user", nil)
	r.Patch("/user", nil)
	r.Delete("/user", nil)
	r.Get("/user/settings", nil)
	r.Post("/user/change_password", nil)

	/* For administrators.  */

	r.Get("/users", middleware.WithBearerAuthorization(h.GetAllUsers))
}
