package pr_service

import (
	"AvitoTask2025/internal/entity"
	"context"

	"go.uber.org/zap"
)

var _ UsersUseCase = (*Impl)(nil)

func (i Impl) SetUserActive(ctx context.Context, user entity.User) (entity.User, error) {
	updatedUser, err := i.UserRepo.UpdateUserActive(ctx, user)
	if err != nil {
		i.logger.Error("failed to update user", zap.Error(err))
		return entity.User{}, err
	}
	return updatedUser, nil
}

func (i Impl) GetUserReviews(ctx context.Context, user entity.User) ([]entity.PullRequest, error) {
	prs, err := i.PRRepo.GetUsersPRs(ctx, user.UserId)
	if err != nil {
		i.logger.Error("failed to get pull requests by reviewer id", zap.Error(err))
		return nil, err
	}
	return prs, nil
}
