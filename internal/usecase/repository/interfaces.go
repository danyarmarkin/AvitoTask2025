package repository

import (
	"AvitoTask2025/internal/entity"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	CreateOrUpdateUser(ctx context.Context, user entity.User) (entity.User, error)
	CheckTeamExists(ctx context.Context, teamName string) (bool, error)
	UpdateUserActive(ctx context.Context, user entity.User) (entity.User, error)
	GetTeamUsers(ctx context.Context, teamName string) ([]entity.User, error)
	GetTeamUsersByAuthorID(ctx context.Context, authorID string) ([]entity.User, error)
}

type PRRepository interface {
	CreatePullRequest(ctx context.Context, pr entity.PullRequest) (entity.PullRequest, error)
	UpdateAssignedReviewers(ctx context.Context, pr entity.PullRequest) (entity.PullRequest, error)
	Merge(ctx context.Context, pr entity.PullRequest) (entity.PullRequest, error)
	GetPRByID(ctx context.Context, prID string) (entity.PullRequest, error)
	GetUsersPRs(ctx context.Context, userID string) ([]entity.PullRequest, error)
}

type Impl struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Impl {
	return &Impl{
		db: db,
	}
}
