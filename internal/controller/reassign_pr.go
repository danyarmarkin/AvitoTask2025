package controller

import (
	"AvitoTask2025/generated/api/pr_service"
	"AvitoTask2025/internal/entity"
	"context"
	"errors"
	"fmt"
)

func (i Impl) PostPullRequestReassign(ctx context.Context, request pr_service.PostPullRequestReassignRequestObject) (pr_service.PostPullRequestReassignResponseObject, error) {
	pr := entity.PullRequest{
		PullRequestId: request.Body.PullRequestId,
	}
	user := entity.User{
		UserId: request.Body.OldUserId,
	}

	prNew, err := i.PRUseCase.ReassignPR(ctx, pr, user)
	if errors.Is(err, entity.ErrPrNotFound) {
		return pr_service.PostPullRequestReassign404JSONResponse{
			Error: struct {
				Code    pr_service.ErrorResponseErrorCode `json:"code"`
				Message string                            `json:"message"`
			}{Code: pr_service.NOTFOUND, Message: fmt.Sprintf("pull request with id=%s not found", pr.PullRequestId)},
		}, nil
	}
	if errors.Is(err, entity.ErrUserNotFound) {
		return pr_service.PostPullRequestReassign404JSONResponse{
			Error: struct {
				Code    pr_service.ErrorResponseErrorCode `json:"code"`
				Message string                            `json:"message"`
			}{Code: pr_service.NOTFOUND, Message: fmt.Sprintf("user with id=%s not found", user.UserId)},
		}, nil
	}
	if errors.Is(err, entity.ErrPRMerged) {
		return pr_service.PostPullRequestReassign409JSONResponse{
			Error: struct {
				Code    pr_service.ErrorResponseErrorCode `json:"code"`
				Message string                            `json:"message"`
			}{Code: pr_service.PRMERGED, Message: fmt.Sprintf("pull request with id=%s is already merged", pr.PullRequestId)},
		}, nil
	}

	if err != nil {
		return nil, err
	}

	return pr_service.PostPullRequestReassign200JSONResponse{
		Pr: pr_service.PullRequest{
			PullRequestId:     prNew.PullRequestId,
			PullRequestName:   prNew.PullRequestName,
			AuthorId:          prNew.AuthorId,
			Status:            pr_service.PullRequestStatus(prNew.Status),
			AssignedReviewers: prNew.AssignedReviewers,
			CreatedAt:         prNew.CreatedAt,
		},
		ReplacedBy: prNew.ReplacedByReviewer,
	}, nil
}
