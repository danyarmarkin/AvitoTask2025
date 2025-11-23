package pr_service

import (
	"AvitoTask2025/internal/entity"
	"context"
	"errors"
	"math/rand"
	"slices"

	"go.uber.org/zap"
)

var _ PRUseCase = (*Impl)(nil)

func (i Impl) CreatePR(ctx context.Context, pr entity.PullRequest) (entity.PullRequest, error) {
	pr.Status = entity.PullRequestStatusOPEN
	users, err := i.UserRepo.GetTeamUsersByAuthorID(ctx, pr.AuthorId)
	if err != nil {
		if !errors.Is(err, entity.ErrUserNotFound) {
			i.logger.Error("failed to get team users by author id", zap.Error(err))
		} else {
			err = entity.ErrAuthorNotFound
		}
		return entity.PullRequest{}, err
	}

	active := make([]entity.User, 0, len(users))
	for _, u := range users {
		if u.IsActive && u.UserId != pr.AuthorId {
			active = append(active, u)
		}
	}
	selected := active
	if len(active) > 2 {
		rand.Shuffle(len(active), func(a, b int) { active[a], active[b] = active[b], active[a] })
		selected = active[:2]
	}
	selectedIDs := make([]string, 0, len(selected))
	for _, u := range selected {
		selectedIDs = append(selectedIDs, u.UserId)
	}

	pr.AssignedReviewers = selectedIDs

	createdPR, err := i.PRRepo.CreatePullRequest(ctx, pr)
	if err != nil {
		if !errors.Is(err, entity.ErrPRAlreadyExists) && !errors.Is(err, entity.ErrAuthorNotFound) {
			i.logger.Error("failed to create pull request", zap.Error(err))
		}
		return entity.PullRequest{}, err
	}
	return createdPR, nil
}

func (i Impl) MergePR(ctx context.Context, pr entity.PullRequest) (entity.PullRequest, error) {
	mergedPR, err := i.PRRepo.Merge(ctx, pr)
	if err != nil {
		if !errors.Is(err, entity.ErrPrNotFound) {
			i.logger.Error("failed to merge pull request", zap.Error(err))
		}
		return entity.PullRequest{}, err
	}
	return mergedPR, nil
}

func (i Impl) ReassignPR(ctx context.Context, pr entity.PullRequest, oldReviewer entity.User) (entity.PullRequest, error) {
	users, err := i.UserRepo.GetTeamUsersByAuthorID(ctx, oldReviewer.UserId)
	if err != nil {
		if !errors.Is(err, entity.ErrUserNotFound) {
			i.logger.Error("failed to get team users by author id", zap.Error(err))
		}
		return entity.PullRequest{}, err
	}

	pr, err = i.PRRepo.GetPRByID(ctx, pr.PullRequestId)
	if err != nil {
		if !errors.Is(err, entity.ErrPrNotFound) {
			i.logger.Error("failed to get pull request by id", zap.Error(err))
		}
		return entity.PullRequest{}, err
	}

	newUserId := ""
	for idx, reviewerID := range pr.AssignedReviewers {
		if reviewerID == oldReviewer.UserId {
			pr.AssignedReviewers = append(pr.AssignedReviewers[:idx], pr.AssignedReviewers[idx+1:]...)
			break
		}
	}
	if len(pr.AssignedReviewers) < 2 {
		active := make([]entity.User, 0, len(users))
		for _, u := range users {
			if u.IsActive && u.UserId != pr.AuthorId && u.UserId != oldReviewer.UserId && !slices.Contains(pr.AssignedReviewers, u.UserId) {
				active = append(active, u)
			}
		}
		if len(active) > 0 {
			newUserId = active[rand.Intn(len(active))].UserId
			pr.AssignedReviewers = append(pr.AssignedReviewers, newUserId)
		}
	}

	pr, err = i.PRRepo.UpdateAssignedReviewers(ctx, pr)
	if err != nil {
		if !errors.Is(err, entity.ErrPrNotFound) && !errors.Is(err, entity.ErrUserNotFound) && !errors.Is(err, entity.ErrPRMerged) {
			i.logger.Error("failed to reassign pull request", zap.Error(err))
		}
		return entity.PullRequest{}, err
	}
	pr.ReplacedByReviewer = newUserId
	return pr, nil
}
