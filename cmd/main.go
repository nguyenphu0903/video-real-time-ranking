package main

import (
	"go-server/config"
	"go-server/internal/infrastructure/router"
	"go-server/internal/registry"
	redis "go-server/pkg"
	"go-server/pkg/mongo"
)

func main() {
	config.LoadConfig()
	mongoDB := mongo.NewMongo(config.C.MongoDB.URLString, config.C.MongoDB.DatabaseName)
	redisClient := redis.InitRedis(config.C.Redis.Host + ":" + config.C.Redis.Port)

	rg := registry.NewInteractor(mongoDB, redisClient)

	masterHandler := rg.NewAppHandler()
	router.Initialize(masterHandler)
}
