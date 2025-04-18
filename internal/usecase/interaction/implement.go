package interaction

import (
	"context"
	"log"
	"time"

	"go-server/internal/common/constant"
	"go-server/internal/entity"
	userinteraction "go-server/internal/entity"

	"github.com/go-redis/redis/v8"
)

// GoroutinePool manages a pool of goroutines to limit concurrency
type GoroutinePool struct {
	pool chan struct{}
}

// NewGoroutinePool creates a new GoroutinePool with the specified size
func NewGoroutinePool(size int) *GoroutinePool {
	return &GoroutinePool{
		pool: make(chan struct{}, size),
	}
}

// Add adds a goroutine to the pool
func (p *GoroutinePool) Add() {
	p.pool <- struct{}{}
}

// Done removes a goroutine from the pool
func (p *GoroutinePool) Done() {
	<-p.pool
}

var goroutinePool = NewGoroutinePool(10)

// Service handles interaction-related business logic
type Service struct {
	repo  Repository
	redis *redis.Client
}

// NewService creates a new Service instance
func NewService(r Repository, redis *redis.Client) *Service {
	return &Service{
		repo:  r,
		redis: redis,
	}
}

// CreateNewInteraction processes a new user interaction and publishes it to Redis
func (s *Service) CreateNewInteraction(
	ctx context.Context, req *userinteraction.UserInteractionReq,
) error {
	// Insert interaction into the database
	if err := s.repo.InsertOne(ctx, &entity.Interaction{
		UserID:          req.UserID,
		VideoID:         req.VideoID,
		InteractionType: req.InteractionType,
		CreatedAt:       time.Now(),
	}); err != nil {
		log.Printf("[CreateNewInteraction] - [InsertOne] - %v", err)
		return err
	}

	// Publish the interaction event asynchronously using GoroutinePool
	goroutinePool.Add()
	go func() {
		defer goroutinePool.Done()
		interactionEvent := entity.InteractionEvent{
			UserID:          req.UserID,
			VideoID:         req.VideoID,
			InteractionType: req.InteractionType,
			PublishedAt:     time.Now(),
		}
		if err := s.redis.Publish(ctx, constant.InteractionEventsChannel, interactionEvent); err != nil {
			log.Printf("[CreateNewInteraction] - [Publish] - %v", err)
		}
	}()

	log.Printf("[CreateNewInteraction] - Interaction created successfully for user: %s", req.UserID)
	return nil
}
