package injector

import (
	"noda/api/repository"
	"noda/api/service"
	"noda/database"
	"sync"
)

var (
	userOnce    sync.Once
	userService *service.UserService
)

func UserService() *service.UserService {
	if userService == nil {
		userOnce.Do(func() {
			rep := repository.NewUserRepository(database.Get())
			userService = service.NewUserService(rep)
		})
	}
	return userService
}
