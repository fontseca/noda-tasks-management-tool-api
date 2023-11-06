package injector

import (
	"noda/service"
	"sync"
)

var (
	authOnce    sync.Once
	authService *service.AuthenticationService
)

func AuthenticationService() *service.AuthenticationService {
	if authService == nil {
		authOnce.Do(func() {
			userService := UserService()
			authService = service.NewAuthenticationService(userService)
		})
	}
	return authService
}
