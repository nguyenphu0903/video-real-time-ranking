package entity

import (
	interactionConstant "go-server/internal/common/constant"
	"time"
)

type UserInteractionReq struct {
	InteractionType interactionConstant.InteractionType `json:"reaction_type" validate:"required"`
	UserID          string                              `json:"user_id" validate:"required"`
	VideoID         string                              `json:"video_id" validate:"required"`
	ReactionAt      time.Time                           `json:"reaction_at" validate:"required"`
}

type InteractionEvent struct {
	UserID          string                              `json:"user_id"`
	VideoID         string                              `json:"video_id"`
	InteractionType interactionConstant.InteractionType `json:"reaction_type"`
	PublishedAt     time.Time                           `json:"published_at"`
}

type Interaction struct {
	UserID          string                              `bson:"user_id" json:"user_id"`
	VideoID         string                              `bson:"video_id" json:"video_id"`
	InteractionType interactionConstant.InteractionType `bson:"interaction_type" json:"interaction_type"`
	CreatedAt       time.Time                           `bson:"created_at" json:"created_at"`
}
