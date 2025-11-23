package pr_service

import (
	"AvitoTask2025/internal/entity"
	"AvitoTask2025/internal/usecase/repository"
	"context"

	"go.uber.org/zap"
)

type TeamsUseCase interface {
	AddTeam(ctx context.Context, team entity.Team) (entity.Team, error)
	GetTeam(ctx context.Context, teamName string) (entity.Team, error)
}

type UsersUseCase interface {
	SetUserActive(ctx context.Context, user entity.User) (entity.User, error)
	GetUserReviews(ctx context.Context, user entity.User) ([]entity.PullRequest, error)
}

type PRUseCase interface {
	CreatePR(ctx context.Context, pr entity.PullRequest) (entity.PullRequest, error)
	MergePR(ctx context.Context, pr entity.PullRequest) (entity.PullRequest, error)
	ReassignPR(ctx context.Context, pr entity.PullRequest, oldReviewer entity.User) (entity.PullRequest, error)
}

type Impl struct {
	UserRepo repository.UserRepository
	PRRepo   repository.PRRepository

	logger *zap.Logger
}

func New(logger *zap.Logger, userRepo repository.UserRepository, prRepo repository.PRRepository) *Impl {
	return &Impl{
		UserRepo: userRepo,
		PRRepo:   prRepo,
		logger:   logger,
	}
}
