package routes

import (
	"noda/api/handler"

	"github.com/go-chi/chi/v5"
)

func InitializeForTask(router *chi.Mux, taskHandler *handler.TaskHandler) {
	/* Task routes for the logged in user.  */
	router.Get("/user/tasks", taskHandler.GetAllFromLoggedInUser)
	router.Post("/user/tasks", taskHandler.CreateOneForLoggedInUser)
	router.Get("/user/tasks/{task_id}", taskHandler.GetByIDFromLoggedInUser)
	router.Patch("/user/tasks/{task_id}", taskHandler.UpdateByIDFromLoggedInUser)
	router.Delete("/user/tasks/{task_id}", taskHandler.DeleteByIDFromLoggedInUser)

	/* General routes tasks.  */
	router.Get("/tasks", taskHandler.GetAll)
	router.Get("/tasks/{task_id}", taskHandler.GetByID)
}
