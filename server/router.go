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
	routeLists(r)
	return r
}

func routeAuthentication(r *chi.Mux) {
	s := getAuthenticationService()
	h := handler.NewAuthenticationHandler(s)
	r.Post("/signup", h.HandleSignUp)
	r.Post("/login", h.HandleSignIn)
	r.Post("/me/logout", nil)
}

func routeUsers(router chi.Router) {
	s := getUserService()
	h := handler.NewUserHandler(s)
	router.
		With(authorization).
		Group(func(r chi.Router) {
			r.Get("/me", h.HandleRetrievalOfLoggedInUser)
			r.Patch("/me", h.HandleUpdateForLoggedUser)
			r.Delete("/me", h.HandleRemovalOfLoggedUser)
			r.Get("/me/settings", h.HandleRetrievalOfLoggedUserSettings)
			r.Get("/me/settings/{setting_key}", h.HandleRetrievalOfOneSettingOfLoggedUser)
			r.Put("/me/settings/{setting_key}", h.HandleUpdateOneSettingForLoggedUser)
			r.Post("/me/change_password", nil)
		})
	/* For administrator.  */
	router.
		With(authorization).
		With(adminPrivileges).
		Group(func(r chi.Router) {
			r.Get("/users", h.HandleUsersRetrieval)
			r.Get("/users/{user_id}", h.HandleRetrievalOfUserByID)
			r.Get("/users/search", h.HandleUsersSearch)
			r.Delete("/users/{user_id}", h.HandleUserDeletion)
			r.Put("/users/{user_id}/block", h.HandleBlockUser)
			r.Delete("/users/{user_id}/block", h.HandleUnblockUser)
			r.Get("/users/blocked", h.HandleBlockedUsersRetrieval)
			r.Put("/users/{user_id}/make_admin", h.HandleAdminPromotion)
			r.Delete("/users/{user_id}/make_admin", h.HandleDegradeAdminToUser)
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

func routeLists(router chi.Router) {
	s := getListService()
	h := handler.NewListHandler(s)
	router.
		With(authorization).
		Group(func(r chi.Router) {
			r.Post("/me/lists", h.HandleScatteredListCreation)
			r.Get("/me/lists", h.HandleRetrievalOfLists)
			r.Get("/me/lists/{list_id}", h.HandleScatteredListRetrievalByID)
			r.Patch("/me/lists/{list_id}", h.HandlePartialUpdateOfScatteredList)
			r.Delete("/me/lists/{list_id}", h.HandleScatteredListDeletion)
			r.Post("/me/groups/{group_id}/lists", h.HandleGroupedListCreation)
			r.Get("/me/groups/{group_id}/lists", h.HandleGroupedListsRetrieval)
			r.Get("/me/groups/{group_id}/lists/{list_id}", h.HandleGroupedListRetrievalByID)
			r.Patch("/me/groups/{group_id}/lists/{list_id}", h.HandlePartialUpdateOfGroupedList)
			r.Delete("/me/groups/{group_id}/lists/{list_id}", h.HandleGroupedListDeletion)
		})
}
