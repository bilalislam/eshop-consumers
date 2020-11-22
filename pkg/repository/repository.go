package repository

import (
	"context"
	"github.com/bilalislam/torc/storage"
	"github.com/bilalislam/torc/storage/clients"
	"github.com/go-redis/redis/v7"
	"github.com/spf13/viper"
	"os"
	"strings"
)

type storeType string

const (
	EventStore     storeType = "eventstore"
	SnapshotStore  storeType = "snapshotstore"
	AggregateStore storeType = "aggregatestore"
)

func NewMongoRepository(ctx context.Context, v *viper.Viper, enum storeType) (error, storage.IRepository) {
	connString := v.GetString("mongodb." + string(enum) + ".connstring")
	database := v.GetString("mongodb." + string(enum) + ".database")
	collection := v.GetString("mongodb." + string(enum) + ".collection")

	var options clients.MongoClientOptions
	options = clients.MongoClientOptions{
		Collection:       collection,
		ConnectionString: connString,
		DbName:           database,
	}

	repository := new(clients.MongoRepository)
	err := repository.Connect(ctx, options)

	if err != nil {
		panic(err)
	}
	return nil, repository
}

func NewRedisClient(ctx context.Context, v *viper.Viper) storage.IRepository {
	var environment = os.Getenv("ENV_FILE")

	var addr = v.GetString("redis.addr")
	var db = v.GetInt("redis.db")
	var password = v.GetString("redis.password")
	var sentinelMasterName = v.GetString("redis.sentinelMasterName")
	repository := new(clients.RedisRepository)

	var err error
	if environment == "prod" {
		err = repository.ConnectFailoverCluster(ctx, redis.FailoverOptions{
			MasterName:    sentinelMasterName,
			SentinelAddrs: strings.Split(addr, ","),
			Password:      password,
			DB:            db,
		})

	} else {
		err = repository.Connect(ctx, redis.Options{
			Addr:     addr,
			Password: password,
			DB:       db,
		})
	}
	if err != nil {
		panic(err.Error())
	}

	return repository
}
