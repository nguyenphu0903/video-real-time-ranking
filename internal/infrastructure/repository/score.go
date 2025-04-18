package repository

import (
	"context"
	"log"

	"go-server/internal/common/constant"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ScoreRepository struct {
	collection         *mongo.Collection
	personalCollection *mongo.Collection
	redisClient        *redis.Client
}

func NewScoreRepository(db *mongo.Database, redisClient *redis.Client) *ScoreRepository {
	return &ScoreRepository{
		collection:         db.Collection("video_scores"),
		personalCollection: db.Collection("personal_scores"),
		redisClient:        redisClient,
	}
}

// GetByVideo retrieves the score of a video by its ID
func (r *ScoreRepository) GetByVideo(
	ctx context.Context, videoID string,
) (float64, error) {
	filter := bson.M{"video_id": videoID}
	var result struct {
		Score float64 `bson:"score"`
	}
	if err := r.collection.FindOne(ctx, filter).Decode(&result); err != nil {
		log.Printf("Failed to get score for video %s: %v", videoID, err)
		return 0, err
	}
	log.Printf("Successfully retrieved score for video %s: %f", videoID, result.Score)
	return result.Score, nil
}

// IncrementScore increments the score of a video by a given value
func (r *ScoreRepository) IncrementScore(ctx context.Context, videoID string, increment float64) error {
	filter := bson.M{"video_id": videoID}
	update := bson.M{
		"$inc": bson.M{
			"score": increment,
		},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("Failed to increment score for video %s: %v", videoID, err)
		return err
	}
	log.Printf("Successfully incremented score for video %s by %f", videoID, increment)
	return nil
}

// InsertOne inserts a new video score into the collection
func (r *ScoreRepository) InsertOne(ctx context.Context, videoID string, score float64) error {
	document := bson.M{
		"video_id": videoID,
		"score":    score,
	}
	_, err := r.collection.InsertOne(ctx, document)
	if err != nil {
		log.Printf("Failed to insert score for video %s: %v", videoID, err)
		return err
	}
	log.Printf("Successfully inserted score for video %s", videoID)
	return nil
}

// GetPersonalScore retrieves the personal score of a user for a specific video
func (r *ScoreRepository) GetPersonalScore(ctx context.Context, userID string, videoID string) (float64, error) {
	filter := bson.M{"user_id": userID, "video_id": videoID}
	var result struct {
		Score float64 `bson:"score"`
	}
	if err := r.personalCollection.FindOne(ctx, filter).Decode(&result); err != nil {
		log.Printf("Failed to get personal score for user %s and video %s: %v", userID, videoID, err)
		return 0, err
	}
	log.Printf("Successfully retrieved personal score for user %s and video %s: %f", userID, videoID, result.Score)
	return result.Score, nil
}

// IncrementPersonalScore increments the personal score of a user for a specific video by a given value
func (r *ScoreRepository) IncrementPersonalScore(ctx context.Context, userID string, videoID string, increment float64) error {
	filter := bson.M{"user_id": userID, "video_id": videoID}

	update := bson.M{
		"$inc": bson.M{
			"score": increment,
		},
	}
	_, err := r.personalCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("Failed to increment personal score for user %s and video %s: %v", userID, videoID, err)
		return err
	}
	log.Printf("Successfully incremented personal score for user %s and video %s by %f", userID, videoID, increment)

	return nil
}

// InsertPersonalScore inserts a new personal score for a user and video into the collection
func (r *ScoreRepository) InsertPersonalScore(ctx context.Context, userID string, videoID string, score float64) error {
	document := bson.M{
		"user_id":  userID,
		"video_id": videoID,
		"score":    score,
	}
	_, err := r.personalCollection.InsertOne(ctx, document)
	if err != nil {
		log.Printf("Failed to insert personal score for user %s and video %s: %v", userID, videoID, err)
		return err
	}
	log.Printf("Successfully inserted personal score for user %s and video %s", userID, videoID)
	return nil
}

// UpdateCachedScore updates the score of a video in Redis for quick ranking
func (r *ScoreRepository) UpdateCachedScore(ctx context.Context, videoID string, score float64) error {
	if err := r.redisClient.ZAdd(ctx, constant.VideoRanking, &redis.Z{
		Score:  score,
		Member: videoID,
	}).Err(); err != nil {
		log.Printf("Failed to cache score for video %s: %v", videoID, err)
		return err
	}
	log.Printf("Cached score for video %s: %f", videoID, score)
	return nil
}

// GetTopRankedVideos retrieves the top N videos from the Redis Sorted Set
func (r *ScoreRepository) GetTopRankedVideos(ctx context.Context, limit int64) ([]string, error) {
	videos, err := r.redisClient.ZRevRange(ctx, constant.VideoRanking, 0, limit-1).Result()
	if err != nil {
		log.Printf("Failed to get top-ranked videos: %v", err)
		return nil, err
	}
	log.Printf("Successfully retrieved top %d ranked videos", limit)
	return videos, nil
}

// UpdatePersonalizedRankingCache updates the personalized ranking cache for a user
func (r *ScoreRepository) UpdatePersonalizedRankingCache(ctx context.Context, userID string, videoID string, score float64) error {
	if err := r.redisClient.ZAdd(ctx, constant.PersonalRankingPrefix+userID, &redis.Z{
		Score:  score,
		Member: videoID,
	}).Err(); err != nil {
		log.Printf("Failed to cache personalized score for user %s and video %s: %v", userID, videoID, err)
		return err
	}
	log.Printf("Cached personalized score for user %s and video %s: %f", userID, videoID, score)
	return nil
}

// GetPersonalTopRankedVideos retrieves the top N personalized
func (r *ScoreRepository) GetPersonalTopRankedVideos(ctx context.Context, userID string, limit int64) ([]string, error) {
	videos, err := r.redisClient.ZRevRange(ctx, constant.PersonalRankingPrefix+userID, 0, limit-1).Result()
	if err != nil {
		log.Printf("Failed to get personalized top-ranked videos for user %s: %v", userID, err)
		return nil, err
	}
	log.Printf("Successfully retrieved personalized top %d ranked videos for user %s", limit, userID)
	return videos, nil
}
