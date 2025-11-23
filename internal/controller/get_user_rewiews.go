package controller

import (
	generated "AvitoTask2025/generated/api/pr_service"
	"context"
)

func (i Impl) GetUsersGetReview(ctx context.Context, request generated.GetUsersGetReviewRequestObject) (generated.GetUsersGetReviewResponseObject, error) {
	return generated.GetUsersGetReview200JSONResponse{
		PullRequests: make([]generated.PullRequestShort, 0),
		UserId:       "it works",
	}, nil
}
