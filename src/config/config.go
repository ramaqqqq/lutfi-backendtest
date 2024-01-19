package config

import (
	"folkatech-customerIdentity/src/pkg/helpers"

	"github.com/spf13/viper"
)

type Config struct {
	MongoDB MongoDBConfig
	Cache   CacheConfig
}

func NewConfig() *Config {
	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	if err := viper.ReadInConfig(); err != nil {
		helpers.Logger("error", "Failed to read config file: "+err.Error())
		panic(err)
	}

	return &Config{
		MongoDB: LoadMongoDBConfig(),
		Cache:   LoadCacheConfig(),
	}
}
