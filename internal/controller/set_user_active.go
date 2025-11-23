package controller

import (
	"AvitoTask2025/generated/api/pr_service"
	"AvitoTask2025/internal/entity"
	"context"
	"errors"
	"fmt"
)

func (i Impl) PostUsersSetIsActive(ctx context.Context, request pr_service.PostUsersSetIsActiveRequestObject) (pr_service.PostUsersSetIsActiveResponseObject, error) {
	user := entity.User{
		UserId:   request.Body.UserId,
		IsActive: request.Body.IsActive,
	}
	updatedUser, err := i.UsersUseCase.SetUserActive(ctx, user)

	if errors.Is(err, entity.ErrUserNotFound) {
		return pr_service.PostUsersSetIsActive404JSONResponse{
			Error: struct {
				Code    pr_service.ErrorResponseErrorCode `json:"code"`
				Message string                            `json:"message"`
			}{Code: pr_service.NOTFOUND, Message: fmt.Sprintf("user with id=%s not found", request.Body.UserId)},
		}, nil
	}

	if err != nil {
		return nil, err
	}

	return pr_service.PostUsersSetIsActive200JSONResponse{
		User: &pr_service.User{
			UserId:   updatedUser.UserId,
			Username: updatedUser.Username,
			TeamName: updatedUser.TeamName,
			IsActive: updatedUser.IsActive,
		},
	}, nil
}
