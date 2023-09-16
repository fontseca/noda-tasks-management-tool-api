package injector

import (
	"noda/api/repository"
	"noda/api/service"
	"noda/database"
	"sync"
)

var (
	taskOnce    sync.Once
	taskService *service.TaskService
)

func TaskService() *service.TaskService {
	if taskService == nil {
		taskOnce.Do(func() {
			rep := repository.NewTaskRepository(database.Get())
			taskService = service.NewTaskService(rep)
		})
	}
	return taskService
}
