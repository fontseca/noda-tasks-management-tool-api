package routes

import (
	"github.com/go-chi/chi/v5"
	"noda/api/handler"
	"noda/engine/internal/injector"
	"noda/engine/internal/middleware"
)

func InitializeForGroups(router chi.Router) {
	s := injector.GroupService()
	h := handler.NewGroupHandler(s)
	router.
		With(middleware.Authorization).
		Group(func(r chi.Router) {
			r.Get("/me/groups", h.HandleGroupsRetrieval)
			r.Post("/me/groups", h.HandleGroupCreation)
			r.Get("/me/groups/{group_id}", h.HandleRetrieveGroupByID)
			r.Patch("/me/groups/{group_id}", h.HandleGroupUpdate)
			r.Delete("/me/groups/{group_id}", h.HandleGroupDeletion)
		})
}
