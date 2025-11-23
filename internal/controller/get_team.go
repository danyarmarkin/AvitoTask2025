package controller

import (
	"AvitoTask2025/generated/api/pr_service"
	"AvitoTask2025/internal/entity"
	"context"
	"errors"
)

func (i Impl) GetTeamGet(ctx context.Context, request pr_service.GetTeamGetRequestObject) (pr_service.GetTeamGetResponseObject, error) {
	teamName := request.Params.TeamName

	team, err := i.TeamsUseCase.GetTeam(ctx, teamName)
	if errors.Is(err, entity.ErrTeamNotFound) {
		return pr_service.GetTeamGet404JSONResponse{
			Error: struct {
				Code    pr_service.ErrorResponseErrorCode `json:"code"`
				Message string                            `json:"message"`
			}{Code: pr_service.NOTFOUND, Message: "resource not found"},
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
	return pr_service.GetTeamGet200JSONResponse{
		TeamName: team.TeamName,
		Members:  members,
	}, nil
}
