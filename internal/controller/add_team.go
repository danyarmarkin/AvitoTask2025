package controller

import (
	"AvitoTask2025/generated/api/pr_service"
	"AvitoTask2025/internal/entity"
	"context"
	"errors"
	"fmt"
)

func (i Impl) PostTeamAdd(ctx context.Context, request pr_service.PostTeamAddRequestObject) (pr_service.PostTeamAddResponseObject, error) {
	users := make([]entity.User, len(request.Body.Members))
	for index, member := range request.Body.Members {
		users[index] = entity.User{
			UserId:   member.UserId,
			TeamName: request.Body.TeamName,
			Username: member.Username,
			IsActive: member.IsActive,
		}
	}
	team, err := i.TeamsUseCase.AddTeam(ctx, entity.Team{
		TeamName: request.Body.TeamName,
		Members:  users,
	})

	if errors.Is(err, entity.ErrTeamAlreadyExists) {
		return pr_service.PostTeamAdd400JSONResponse{
			Error: struct {
				Code    pr_service.ErrorResponseErrorCode `json:"code"`
				Message string                            `json:"message"`
			}{Code: pr_service.TEAMEXISTS, Message: fmt.Sprintf("%s already exists.", request.Body.TeamName)},
		}, nil
	}
	if err != nil {
		return nil, err
	}

	members := make([]pr_service.TeamMember, len(team.Members))
	for index, member := range team.Members {
		members[index] = pr_service.TeamMember{
			UserId:   member.UserId,
			Username: member.Username,
			IsActive: member.IsActive,
		}
	}
	return pr_service.PostTeamAdd201JSONResponse{
		Team: &pr_service.Team{
			TeamName: team.TeamName,
			Members:  members,
		},
	}, nil
}
