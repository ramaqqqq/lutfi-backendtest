package auth

import (
	"folkatech-customerIdentity/src/modules/auth/controller"
	"folkatech-customerIdentity/src/modules/auth/service"
	"folkatech-customerIdentity/src/modules/user"

	"github.com/go-playground/validator"
)

var (
	authService    service.AuthService
	authController controller.AuthController
)

type Module interface {
	InitModule()
}

type ModuleImpl struct {
}

func New() Module {
	return &ModuleImpl{}
}

func (module ModuleImpl) InitModule() {
	authService = service.NewAuthService(user.GetService(), validator.New())
	authController = controller.NewAuthController(authService)
}

func GetAuthService() service.AuthService {
	return authService
}

func GetAuthController() controller.AuthController {
	return authController
}
