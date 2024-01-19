package db

import (
	"context"
	"fmt"
	"folkatech-customerIdentity/src/config"
	"folkatech-customerIdentity/src/pkg/helpers"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitMongoDB(cfg *config.Config) *mongo.Client {
	mongoDBConfig := cfg.MongoDB

	db, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoDBConfig.DSN))
	if err != nil {
		helpers.Logger("error", "Failed to connect to MongoDB: "+err.Error())
		panic(err)
	}

	if err := db.Ping(context.TODO(), nil); err != nil {
		helpers.Logger("error", "Failed to ping MongoDB: "+err.Error())
		panic(err)
	}

	return db
}

func MgoCollection(coll string, client *mongo.Client) *mongo.Collection {
	dbName := viper.GetString("MONGODB_NAME")
	helpers.Logger("info", fmt.Sprintf("Using MongoDB database: %s", dbName))
	return client.Database(dbName).Collection(coll)
}
