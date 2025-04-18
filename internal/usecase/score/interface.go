package score

import (
	"context"
	"go-server/internal/entity"
)

type Action interface {
	GetByVideo(ctx context.Context, videoID string) (float64, error)
	IncrementScore(ctx context.Context, videoID string, increment float64) error
	InsertOne(ctx context.Context, videoID string, score float64) error
	GetPersonalScore(ctx context.Context, userID string, videoID string) (float64, error)
	IncrementPersonalScore(ctx context.Context, userID string, videoID string, increment float64) error
	InsertPersonalScore(ctx context.Context, userID string, videoID string, score float64) error
}

type Cache interface {
	UpdateCachedScore(ctx context.Context, videoID string, score float64) error
	GetTopRankedVideos(ctx context.Context, limit int64) ([]string, error)
	GetPersonalTopRankedVideos(ctx context.Context, userID string, limit int64) ([]string, error)
	UpdatePersonalizedRankingCache(ctx context.Context, userID string, videoID string, score float64) error
}

type Repository interface {
	Action
	Cache
}

type UseCase interface {
	StartEventConsumer(ctx context.Context)
	UpdateVideoScoreInDB(ctx context.Context, event *entity.InteractionEvent) error
	ListTopRankedVideos(ctx context.Context, limit int) ([]string, error)
	ListPersonalTopRankedVideos(ctx context.Context, userID string, limit int) ([]string, error)
}
