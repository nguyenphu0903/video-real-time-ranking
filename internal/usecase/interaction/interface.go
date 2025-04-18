package interaction

import (
	"context"

	"go-server/internal/entity"
	userinteraction "go-server/internal/entity"
)

type Action interface {
	InsertOne(ctx context.Context, interactionData *entity.Interaction) error
}

type Repository interface {
	Action
}

type UseCase interface {
	CreateNewInteraction(ctx context.Context, req *userinteraction.UserInteractionReq) error
}
