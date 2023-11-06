package injector

import (
	"noda"
	"noda/api/repository"
	"noda/api/service"
	"sync"
)

var (
	userOnce    sync.Once
	userService *service.UserService
)

func UserService() *service.UserService {
	if userService == nil {
		userOnce.Do(func() {
			rep := repository.NewUserRepository(noda.Database())
			userService = service.NewUserService(rep)
		})
	}
	return userService
}
