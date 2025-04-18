package registry

import (
	"go-server/internal/api/handler"
	"go-server/pkg/mongo"

	"github.com/go-redis/redis/v8"
)

type interactor struct {
	mongo mongo.MongoDB
	redis *redis.Client
}

// Interactor Interactor interface
type Interactor interface {
	NewAppHandler() handler.AppHandler
	NewScoreHandler() handler.ScoreHandler
}

// NewInteractor Constructs new interactor
func NewInteractor(mg mongo.MongoDB, redisClient *redis.Client) Interactor {
	return &interactor{mongo: mg, redis: redisClient}
}

func (i *interactor) NewAppHandler() handler.AppHandler {
	return handler.AppHandler{
		InteractionHandler: i.NewInteractionHandler(),
		ScoreHandler:       i.NewScoreHandler(),
	}
}
