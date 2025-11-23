package controller

import (
	generated "AvitoTask2025/generated/api/pr_service"
	"AvitoTask2025/internal/usecase/pr_service"

	"go.uber.org/zap"
)

var _ generated.StrictServerInterface = (*Impl)(nil)

type Impl struct {
	logger *zap.Logger

	TeamsUseCase pr_service.TeamsUseCase
	UsersUseCase pr_service.UsersUseCase
	PRUseCase    pr_service.PRUseCase
}

func NewPrServiceController(logger *zap.Logger, teamsUseCase pr_service.TeamsUseCase, usersUseCase pr_service.UsersUseCase, prUseCase pr_service.PRUseCase) Impl {
	return Impl{
		logger:       logger,
		TeamsUseCase: teamsUseCase,
		UsersUseCase: usersUseCase,
		PRUseCase:    prUseCase,
	}
}
