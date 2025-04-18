package score

import (
	"context"
	"encoding/json"
	"log"

	"go-server/internal/common/constant"
	"go-server/internal/entity"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
)

// ScoreService handles score-related business logic
// It interacts with Redis for caching and MongoDB for persistence
type ScoreService struct {
	redisClient *redis.Client
	repo        Repository
}

// NewScoreService creates a new instance of ScoreService
func NewScoreService(r Repository, redisClient *redis.Client) *ScoreService {
	return &ScoreService{
		redisClient: redisClient,
		repo:        r,
	}
}

// StartEventConsumer listens to Redis Pub/Sub for interaction events and processes them
func (s *ScoreService) StartEventConsumer(ctx context.Context) {
	sub := s.redisClient.Subscribe(ctx, constant.InteractionEventsChannel)
	ch := sub.Channel()

	log.Println("Started Redis consumer for interaction events")

	for msg := range ch {
		go func(payload string) {
			var event entity.InteractionEvent
			if err := json.Unmarshal([]byte(payload), &event); err != nil {
				log.Printf("Failed to unmarshal event: %v", err)
				return
			}

			log.Printf("Processing event: %v", event)
			if err := s.UpdateVideoScoreInDB(ctx, &event); err != nil {
				log.Printf("Failed to update video score in DB: %v", err)
			}
			if err := s.UpdatePersonalizedScore(ctx, &event); err != nil {
				log.Printf("Failed to update personalized score: %v", err)
			}
			log.Printf("Successfully processed event for user %s on video %s", event.UserID, event.VideoID)
		}(msg.Payload)
	}
}

// UpdatePersonalizedScore updates the personalized score for a user and video
// It updates both the database and the Redis cache
func (s *ScoreService) UpdatePersonalizedScore(ctx context.Context, event *entity.InteractionEvent) error {
	videoID := event.VideoID
	userID := event.UserID
	newScore := event.InteractionType.GetScore()

	currentScore, err := s.repo.GetPersonalScore(ctx, userID, videoID)
	if err != nil && err != mongo.ErrNoDocuments {
		log.Printf("Failed to get current personalized score for user %s on video %s: %v", userID, videoID, err)
		return err
	}

	if err == mongo.ErrNoDocuments {
		log.Printf("No current personalized score found for user %s on video %s, inserting new score", userID, videoID)
		if err := s.repo.InsertPersonalScore(ctx, userID, videoID, newScore); err != nil {
			log.Printf("Failed to insert new personalized score for user %s on video %s: %v", userID, videoID, err)
			return err
		}
	} else {
		log.Printf("Updating personalized score for user %s on video %s", userID, videoID)
		if err := s.repo.IncrementPersonalScore(ctx, userID, videoID, newScore); err != nil {
			log.Printf("Failed to update personalized score for user %s on video %s: %v", userID, videoID, err)
			return err
		}
		newScore += currentScore
	}

	log.Printf("Updating personalized ranking cache for user %s on video %s", userID, videoID)
	go func() {
		if err := s.repo.UpdatePersonalizedRankingCache(ctx, userID, videoID, newScore); err != nil {
			log.Printf("Failed to update personalized ranking cache for user %s on video %s: %v", userID, videoID, err)
		}
	}()

	return nil
}

// UpdateVideoScoreInDB updates the global score for a video in the database and cache
func (s *ScoreService) UpdateVideoScoreInDB(ctx context.Context, event *entity.InteractionEvent) error {
	videoID := event.VideoID
	newScore := event.InteractionType.GetScore()

	currentScore, err := s.repo.GetByVideo(ctx, videoID)
	if err != nil && err != mongo.ErrNoDocuments {
		log.Printf("Failed to get current score for video %s: %v", videoID, err)
		return err
	}

	if err == mongo.ErrNoDocuments {
		log.Printf("No current score found for video %s, inserting new score", videoID)
		if err := s.repo.InsertOne(ctx, videoID, newScore); err != nil {
			log.Printf("Failed to insert new score for video %s: %v", videoID, err)
			return err
		}
	} else {
		log.Printf("Updating score for video %s", videoID)
		if err := s.repo.IncrementScore(ctx, videoID, newScore); err != nil {
			log.Printf("Failed to update score for video %s: %v", videoID, err)
			return err
		}
		newScore += currentScore
	}

	go func() {
		if err := s.repo.UpdateCachedScore(ctx, videoID, newScore); err != nil {
			log.Printf("Failed to update cached score for video %s: %v", videoID, err)
		}
	}()

	return nil
}

// ListTopRankedVideos retrieves a list of top-ranked video IDs from the repository
func (s *ScoreService) ListTopRankedVideos(ctx context.Context, limit int) ([]string, error) {
	videos, err := s.repo.GetTopRankedVideos(ctx, int64(limit))
	if err != nil {
		log.Printf("Failed to get top ranked videos: %v", err)
		return nil, err
	}
	return videos, nil
}

// ListPersonalTopRankedVideos retrieves a list of top-ranked video IDs for a specific user
func (s *ScoreService) ListPersonalTopRankedVideos(ctx context.Context, userID string, limit int) ([]string, error) {
	videos, err := s.repo.GetPersonalTopRankedVideos(ctx, userID, int64(limit))
	if err != nil {
		log.Printf("Failed to get personalized top ranked videos for user %s: %v", userID, err)
		return nil, err
	}
	return videos, nil
}
