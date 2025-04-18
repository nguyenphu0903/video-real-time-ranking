package repository

import (
	"context"
	"go-server/internal/entity"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
)

var InteractionCollectionName = "interactions"

type InteractionRepository struct {
	dbMongo *mongo.Database
}

// NewInteractionRepository initializes the repository
func NewInteractionRepository(mg *mongo.Database) *InteractionRepository {
	return &InteractionRepository{dbMongo: mg}
}

// Insert Interaction inserts a new interaction into the database
func (repo *InteractionRepository) InsertOne(ctx context.Context, interactionData *entity.Interaction) error {
	if _, err := repo.dbMongo.Collection(InteractionCollectionName).InsertOne(ctx, interactionData); err != nil {
		log.Printf("Error inserting interaction: %v", err)
		return err
	}

	return nil
}
