package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"noda/handler"
)

func startRouter() *chi.Mux {
	var r = chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.SetHeader("Content-Type", "application/json"))
	r.Use(middleware.AllowContentType("application/json"))
	r.Use(middleware.SetHeader("Access-Control-Allow-Origin", "*"))
	r.Use(middleware.SetHeader("Access-Control-Allow-Methods", "GET"))
	r.Use(middleware.SetHeader("Access-Control-Allow-Headers", "*"))
	r.Use(middleware.SetHeader("Access-Control-Allow-Credentials", "true"))
	r.Use(letOptionsPassThrough)
	r.NotFound(notFound)
	routeAuthentication(r)
	routeUsers(r)
	routeGroups(r)
	routeTasks(r)
	return r
}

func routeAuthentication(r *chi.Mux) {
	s := getAuthenticationService()
	h := handler.NewAuthenticationHandler(s)
	r.Post("/signup", h.SignUp)
	r.Post("/signin", h.SignIn)
	r.Post("/me/signout", nil)
}

func routeUsers(router chi.Router) {
	s := getUserService()
	h := handler.NewUserHandler(s)
	router.
		With(authorization).
		Group(func(r chi.Router) {
			r.Get("/me", h.RetrieveCurrentUser)
			r.Patch("/me", h.UpdateCurrentUser)
			r.Delete("/me", h.RemoveCurrentUser)
			r.Get("/me/settings", h.RetrieveCurrentUserSettings)
			r.Get("/me/settings/{setting_key}", h.RetrieveOneSettingOfCurrentUser)
			r.Put("/me/settings/{setting_key}", h.UpdateOneSettingForCurrentUser)
			r.Post("/me/change_password", nil)
		})
	/* For administrator.  */
	router.
		With(authorization).
		With(adminPrivileges).
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

func routeGroups(router chi.Router) {
	s := getGroupService()
	h := handler.NewGroupHandler(s)
	router.
		With(authorization).
		Group(func(r chi.Router) {
			r.Get("/me/groups", h.HandleGroupsRetrieval)
			r.Post("/me/groups", h.HandleGroupCreation)
			r.Get("/me/groups/{group_id}", h.HandleRetrieveGroupByID)
			r.Patch("/me/groups/{group_id}", h.HandleGroupUpdate)
			r.Delete("/me/groups/{group_id}", h.HandleGroupDeletion)
		})
}

func routeTasks(router chi.Router) {
	s := getTaskService()
	h := handler.NewTaskHandler(s)
	router.
		With(authorization).
		Group(func(r chi.Router) {
			r.Get("/me/tasks", h.RetrieveTasksFromUser)
			r.Post("/me/tasks", nil)
			r.Get("/me/tasks/{task_id}", nil)
			r.Patch("/me/tasks/{task_id}", nil)
			r.Delete("/me/tasks/{task_id}", nil)
		})
}
