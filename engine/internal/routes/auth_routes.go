package routes

import (
	"noda/engine/internal/injector"
	"noda/handler"

	"github.com/go-chi/chi/v5"
)

func InitializeForAuthentication(r *chi.Mux) {
	s := injector.AuthenticationService()
	h := handler.NewAuthenticationHandler(s)

	/* For logged in user.  */

	r.Post("/user/signout", nil)

	/* For anyone.  */

	r.Post("/signup", h.SignUp)
	r.Post("/signin", h.SignIn)
}
