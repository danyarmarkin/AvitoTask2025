package controller

import (
	generated "AvitoTask2025/generated/api/pr_service"
	"AvitoTask2025/internal/entity"
	"context"
)

func (i Impl) GetUsersGetReview(ctx context.Context, request generated.GetUsersGetReviewRequestObject) (generated.GetUsersGetReviewResponseObject, error) {
	user := entity.User{
		UserId: request.Params.UserId,
	}

	reviews, err := i.UsersUseCase.GetUserReviews(ctx, user)
	if err != nil {
		return nil, err
	}
	reviewResponses := make([]generated.PullRequestShort, len(reviews))
	for index, review := range reviews {
		reviewResponses[index] = generated.PullRequestShort{
			PullRequestId:   review.PullRequestId,
			PullRequestName: review.PullRequestName,
			AuthorId:        review.AuthorId,
			Status:          generated.PullRequestShortStatus(review.Status),
		}
	}

	return generated.GetUsersGetReview200JSONResponse{
		UserId:       user.UserId,
		PullRequests: reviewResponses,
	}, nil
}
