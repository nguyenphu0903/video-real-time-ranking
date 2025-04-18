package registry

import (
	"go-server/internal/api/handler"
	"go-server/internal/infrastructure/repository"
	"go-server/internal/usecase/interaction"
)

func (i *interactor) NewInteractionRepository() *repository.InteractionRepository {
	return repository.NewInteractionRepository(i.mongo)
}

func (i *interactor) NewInteractionService() *interaction.Service {
	return interaction.NewService(i.NewInteractionRepository(), i.redis)
}

func (i *interactor) NewInteractionHandler() handler.InteractionHandler {
	return handler.NewInteractionHandler(i.NewInteractionService())
}
