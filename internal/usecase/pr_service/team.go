package pr_service

import (
	"AvitoTask2025/internal/entity"
	"context"
	"errors"

	"go.uber.org/zap"
)

var _ TeamsUseCase = (*Impl)(nil)

func (i Impl) AddTeam(ctx context.Context, team entity.Team) (entity.Team, error) {
	exists, err := i.UserRepo.CheckTeamExists(ctx, team.TeamName)
	if err != nil {
		i.logger.Error("failed to check team exists", zap.Error(err))
	}
	if exists {
		return entity.Team{}, entity.ErrTeamAlreadyExists
	}

	users := make([]entity.User, len(team.Members))

	for idx, member := range team.Members {
		user, err := i.UserRepo.CreateOrUpdateUser(ctx, member)
		if err != nil {
			i.logger.Error("failed to create user", zap.Error(err))
			return entity.Team{}, err
		}
		users[idx] = user
	}

	return entity.Team{
		TeamName: team.TeamName,
		Members:  users,
	}, nil
}

func (i Impl) GetTeam(ctx context.Context, teamName string) (entity.Team, error) {
	users, err := i.UserRepo.GetTeamUsers(ctx, teamName)
	if err != nil {
		if !errors.Is(err, entity.ErrTeamNotFound) {
			i.logger.Error("failed to get team users", zap.Error(err))
		}
		return entity.Team{}, err
	}

	return entity.Team{
		TeamName: teamName,
		Members:  users,
	}, nil
}
