package db

import (
	"folkatech-customerIdentity/src/config"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

type Conn struct {
	Redis *redis.Client
	Mongo *mongo.Client
}

func NewDbConnection(cfg *config.Config) *Conn {
	return &Conn{
		Redis: InitRedis(cfg),
		Mongo: InitMongoDB(cfg),
	}
}
