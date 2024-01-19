package user

import (
	"folkatech-customerIdentity/src/config"
	"folkatech-customerIdentity/src/modules/user/controller"
	"folkatech-customerIdentity/src/modules/user/repo"
	"folkatech-customerIdentity/src/modules/user/service"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	userRepository repo.UserRepo
	userService    service.UserServ
	userController controller.UserController
)

type Module interface {
	InitModule()
}

type ModuleImpl struct {
	client *mongo.Client
	redis  redis.Client
}

func New(database *mongo.Client, redis redis.Client) Module {
	return &ModuleImpl{client: database, redis: redis}
}

func (module ModuleImpl) InitModule() {
	userRepository = repo.NewUserRepository(module.client, &config.RedisCfg{Conn: &module.redis}, config.MongoDBConfig{})
	userService = service.NewUserService(userRepository)
	userController = controller.NewUserController(userService)
}

func GetUserController() controller.UserController {
	return userController
}

func GetService() service.UserServ {
	return userService
}

func GetUserRepository() repo.UserRepo {
	return userRepository
}
