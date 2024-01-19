package config

import (
	"github.com/spf13/viper"
)

type MongoDBConfig struct {
	DSN string
}

func LoadMongoDBConfig() MongoDBConfig {
	return MongoDBConfig{
		DSN: viper.GetString("MONGODB_DSN"),
	}
}
