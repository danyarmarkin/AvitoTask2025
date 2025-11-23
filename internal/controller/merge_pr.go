package controller

import (
	"AvitoTask2025/generated/api/pr_service"
	"AvitoTask2025/internal/entity"
	"context"
	"errors"
	"fmt"
)

func (i Impl) PostPullRequestMerge(ctx context.Context, request pr_service.PostPullRequestMergeRequestObject) (pr_service.PostPullRequestMergeResponseObject, error) {
	pr := entity.PullRequest{
		PullRequestId: request.Body.PullRequestId,
	}

	mergedPR, err := i.PRUseCase.MergePR(ctx, pr)
	if errors.Is(err, entity.ErrPrNotFound) {
		return pr_service.PostPullRequestMerge404JSONResponse{
			Error: struct {
				Code    pr_service.ErrorResponseErrorCode `json:"code"`
				Message string                            `json:"message"`
			}{Code: pr_service.NOTFOUND, Message: fmt.Sprintf("pull request with id=%s not found", pr.PullRequestId)},
		}, nil
	}

	if err != nil {
		return nil, err
	}

	return pr_service.PostPullRequestMerge200JSONResponse{
		Pr: &pr_service.PullRequest{
			PullRequestId:     mergedPR.PullRequestId,
			PullRequestName:   mergedPR.PullRequestName,
			AuthorId:          mergedPR.AuthorId,
			Status:            pr_service.PullRequestStatus(mergedPR.Status),
			AssignedReviewers: mergedPR.AssignedReviewers,
			CreatedAt:         mergedPR.CreatedAt,
			MergedAt:          mergedPR.MergedAt,
		},
	}, nil
}
