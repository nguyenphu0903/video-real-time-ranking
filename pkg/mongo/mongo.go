package mongo

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/mgocompat"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB *mongo.Database

func NewMongo(dbURL, dbName string) (db MongoDB) {
	var (
		ctx          = context.Background()
		currentEntry int
		// TODO: Get values from .env
		maxEntry, waitTime = 3, 5
	)

	for {
		currentEntry++

		if currentEntry > maxEntry {
			log.Fatalf("MongoDB connection errors. Exit after try %d times", maxEntry)
		}

		client, err := mongo.Connect(
			ctx,
			options.Client().ApplyURI(dbURL),
			options.Client().SetRegistry(mgocompat.Registry),
		)
		if err != nil {
			log.Printf("%s. %d times try \n", err.Error(), currentEntry)
			time.Sleep(time.Duration(waitTime) * time.Second)

			continue
		}

		db = client.Database(dbName)

		return
	}
}
