package controller

import (
	generated "AvitoTask2025/generated/api/pr_service"
	"AvitoTask2025/internal/entity"
	"context"
	"errors"
	"fmt"
)

func (i Impl) PostPullRequestCreate(ctx context.Context, request generated.PostPullRequestCreateRequestObject) (generated.PostPullRequestCreateResponseObject, error) {
	pr := entity.PullRequest{
		PullRequestId:   request.Body.PullRequestId,
		PullRequestName: request.Body.PullRequestName,
		AuthorId:        request.Body.AuthorId,
	}

	createdPR, err := i.PRUseCase.CreatePR(ctx, pr)
	if errors.Is(err, entity.ErrAuthorNotFound) {
		return generated.PostPullRequestCreate404JSONResponse{
			Error: struct {
				Code    generated.ErrorResponseErrorCode `json:"code"`
				Message string                           `json:"message"`
			}{Code: generated.NOTFOUND, Message: fmt.Sprintf("author with id=%s not found", pr.AuthorId)},
		}, nil
	}
	if errors.Is(err, entity.ErrPRAlreadyExists) {
		return generated.PostPullRequestCreate409JSONResponse{
			Error: struct {
				Code    generated.ErrorResponseErrorCode `json:"code"`
				Message string                           `json:"message"`
			}{Code: generated.PREXISTS, Message: fmt.Sprintf("pull request with id=%s already exists", pr.PullRequestId)},
		}, nil
	}

	if err != nil {
		return nil, err
	}
	return generated.PostPullRequestCreate201JSONResponse{
		Pr: &generated.PullRequest{
			PullRequestId:     createdPR.PullRequestId,
			PullRequestName:   createdPR.PullRequestName,
			AuthorId:          createdPR.AuthorId,
			Status:            generated.PullRequestStatus(createdPR.Status),
			AssignedReviewers: createdPR.AssignedReviewers,
			CreatedAt:         createdPR.CreatedAt,
		},
	}, nil
}
