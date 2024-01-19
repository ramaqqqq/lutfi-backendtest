package repo

import (
	"context"
	"fmt"
	"folkatech-customerIdentity/src/config"
	"folkatech-customerIdentity/src/modules/user/model"
	"folkatech-customerIdentity/src/pkg/db"
	"folkatech-customerIdentity/src/pkg/helpers"
	"folkatech-customerIdentity/src/pkg/utils"
	"math"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	DeleteUserRedisKey  = "user:*"
	GetListUserRedisKey = "user:getlist"
	GetByIdRedisKey     = "user:get:%d"
	GetByEmailRedisKey  = "user:getemail"
)

type UserRepoImpl struct {
	client *mongo.Client
	redis  config.Redis
	config config.MongoDBConfig
}

func NewUserRepository(client *mongo.Client, redis config.Redis, config config.MongoDBConfig) UserRepo {
	return &UserRepoImpl{client: client, redis: redis, config: config}
}

func (r *UserRepoImpl) Save(ctx context.Context, user model.User) (*model.User, error) {
	nextID, err := r.getNextSequence(ctx, "users")
	if err != nil {
		helpers.Logger("error", "In Server: "+err.Error())
		return nil, err
	}

	user.ID = nextID

	_, err = db.MgoCollection("users", r.client).InsertOne(ctx, user)
	if err != nil {
		helpers.Logger("error", "In Server: "+err.Error())
		return nil, err
	}

	err = r.redis.DelWithPattern(ctx, DeleteUserRedisKey)
	if err != nil {
		helpers.Logger("error", fmt.Sprintf("error when deleted with pattern redis: %s", err))
		return nil, err
	}
	return &user, nil
}

func (r *UserRepoImpl) Update(ctx context.Context, userID int64, user model.User) error {
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "username", Value: user.Username},
			{Key: "account_number", Value: user.AccountNumber},
			{Key: "email_address", Value: user.EmailAddress},
			{Key: "identity_number", Value: user.IdentityNumber},
		}},
	}

	filter := bson.D{{Key: "_id", Value: userID}}
	_, err := db.MgoCollection("users", r.client).UpdateOne(ctx, filter, update)
	if err != nil {
		helpers.Logger("error", "In Server: "+err.Error())
		return err
	}

	err = r.redis.DelWithPattern(ctx, DeleteUserRedisKey)
	if err != nil {
		helpers.Logger("error", fmt.Sprintf("error when deleted with pattern redis: %s", err))
		return err
	}

	return nil
}

func (r *UserRepoImpl) Delete(ctx context.Context, userID int64) error {
	filter := bson.D{{Key: "_id", Value: userID}}
	_, err := db.MgoCollection("users", r.client).DeleteOne(ctx, filter)
	if err != nil {
		helpers.Logger("error", "In Server: "+err.Error())
		return err
	}

	err = r.redis.DelWithPattern(ctx, DeleteUserRedisKey)
	if err != nil {
		helpers.Logger("error", fmt.Sprintf("error when deleted with pattern redis: %s", err))
		return err
	}

	return nil
}

func (r *UserRepoImpl) Find(ctx context.Context, filter model.FilterUser, pg *utils.PaginateQueryOffset) ([]model.User, int64, int64, error) {
	var (
		data      []model.User
		totalData int64
		totalPage int64
		err       error
	)

	filterBSON, err := buildFilterBSON(filter)
	if err != nil {
		helpers.Logger("error", "Error building filter BSON: "+err.Error())
		return nil, 0, 0, err
	}

	if err := r.redis.WithCache(ctx, GetListUserRedisKey, &data, func() (interface{}, error) {
		var user []model.User
		sortField := pg.Order.Field
		sortDirection := 1
		if sortField == "" {
			sortField = "_id"
		}

		options := options.Find()
		options.SetSort(bson.D{{Key: sortField, Value: sortDirection}})
		options.SetLimit(int64(pg.Limit))
		options.SetSkip(int64(pg.Offset))

		cursor, err := db.MgoCollection("users", r.client).Find(ctx, filterBSON, options)
		if err != nil {
			helpers.Logger("error", "Error querying MongoDB: "+err.Error())
			return nil, err
		}

		defer cursor.Close(ctx)

		if err := cursor.All(ctx, &user); err != nil {
			helpers.Logger("error", "Error decoding MongoDB cursor: "+err.Error())
			return nil, err
		}

		totalData, err = r.getTotal(ctx, filter)
		if err != nil {
			return nil, err
		}

		if pg.Limit > 0 {
			totalPage = int64(math.Ceil(float64(totalData) / float64(pg.Limit)))
		} else {
			totalPage = 1
		}

		return user, nil

	}); err != nil {
		return data, 0, 0, err
	}

	return data, totalData, totalPage, nil
}

func (r *UserRepoImpl) FindByID(ctx context.Context, id int64) (model.User, error) {
	var user model.User

	err := r.redis.WithCache(ctx, fmt.Sprintf(GetByIdRedisKey, id), &user, func() (interface{}, error) {
		var data model.User
		filter := bson.D{{Key: "_id", Value: id}}
		err := r.findUser(ctx, filter, &data)
		if err != nil {
			return data, err
		}
		return data, nil
	})
	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *UserRepoImpl) FindByEmail(ctx context.Context, email string) (model.User, error) {
	var data model.User
	filter := bson.D{{Key: "email_address", Value: email}}
	err := r.findUser(ctx, filter, &data)
	if err != nil {
		return data, nil
	}

	return data, nil
}

func (r *UserRepoImpl) findUser(ctx context.Context, filter interface{}, user *model.User) error {
	err := db.MgoCollection("users", r.client).FindOne(ctx, filter).Decode(user)
	if err != nil {
		helpers.Logger("error", "In Server: "+err.Error())
		return err
	}

	return nil
}

func (r *UserRepoImpl) getNextSequence(ctx context.Context, collectionName string) (int64, error) {
	counterCollection := r.client.Database("lutfi_folkatechdb").Collection("counters")

	filter := bson.M{"_id": collectionName}
	update := bson.M{"$inc": bson.M{"value": 1}}

	after := options.After
	opt := options.FindOneAndUpdateOptions{ReturnDocument: &after}
	result := counterCollection.FindOneAndUpdate(ctx, filter, update, &opt)

	var counter model.Counter
	err := result.Decode(&counter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			initialCounter := model.Counter{ID: collectionName, Value: 1}
			_, err := counterCollection.InsertOne(ctx, initialCounter)
			if err != nil {
				return 0, err
			}
			return initialCounter.Value, nil
		}
		return 0, err
	}

	return counter.Value, nil
}

func (r *UserRepoImpl) getTotal(ctx context.Context, filter model.FilterUser) (int64, error) {
	filterBSON, err := buildFilterBSON(filter)
	if err != nil {
		helpers.Logger("error", "Error building filter BSON: "+err.Error())
		return 0, err
	}

	totalData, err := db.MgoCollection("users", r.client).CountDocuments(ctx, filterBSON)
	if err != nil {
		helpers.Logger("error", "Error querying MongoDB for total count: "+err.Error())
		return 0, err
	}

	return totalData, nil
}

func buildFilterBSON(filter model.FilterUser) (bson.M, error) {
	filterBSON := bson.M{}

	if filter.AccountNumber != "" {
		filterBSON["account_number"] = filter.AccountNumber
	}

	if filter.IdentityNumber != "" {
		filterBSON["identity_number"] = filter.IdentityNumber
	}

	if filter.Search != "" {
		searchRegex := bson.M{"$regex": primitive.Regex{Pattern: filter.Search, Options: "i"}}
		filterBSON["$or"] = []bson.M{
			{"email_address": searchRegex},
			{"username": searchRegex},
		}
	}

	return filterBSON, nil
}
