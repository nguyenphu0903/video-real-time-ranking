package registry

import (
	"go-server/internal/api/handler"
	"go-server/internal/infrastructure/repository"
	"go-server/internal/usecase/score"
)

func (i *interactor) NewScoreRepository() *repository.ScoreRepository {
	return repository.NewScoreRepository(i.mongo, i.redis)
}

func (i *interactor) NewScoreService() *score.ScoreService {
	return score.NewScoreService(i.NewScoreRepository(), i.redis)
}

func (i *interactor) NewScoreHandler() handler.ScoreHandler {
	return handler.NewScoreHandler(i.NewScoreService())
}
