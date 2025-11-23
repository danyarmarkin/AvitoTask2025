package repository

import (
	"AvitoTask2025/internal/entity"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

var _ UserRepository = (*Impl)(nil)

func (i Impl) CreateOrUpdateUser(ctx context.Context, user entity.User) (entity.User, error) {
	const query = `
		INSERT INTO users (user_id, username, team_name, is_active)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id) DO UPDATE
		SET username = EXCLUDED.username, team_name = EXCLUDED.team_name, is_active = EXCLUDED.is_active
		RETURNING user_id, username, team_name, is_active;`
	row := i.db.QueryRow(ctx, query, user.UserId, user.Username, user.TeamName, user.IsActive)
	var updatedUser entity.User
	err := row.Scan(&updatedUser.UserId, &updatedUser.Username, &updatedUser.TeamName, &updatedUser.IsActive)
	if err != nil {
		return entity.User{}, err
	}
	return updatedUser, nil
}

func (i Impl) UpdateUserActive(ctx context.Context, user entity.User) (entity.User, error) {
	const query = `
		UPDATE users
		SET is_active = $1
		WHERE user_id = $2
		RETURNING user_id, username, team_name, is_active;`
	row := i.db.QueryRow(ctx, query, user.IsActive, user.UserId)
	var updatedUser entity.User
	err := row.Scan(&updatedUser.UserId, &updatedUser.Username, &updatedUser.TeamName, &updatedUser.IsActive)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.User{}, entity.ErrUserNotFound
		}
		return entity.User{}, err
	}
	return updatedUser, nil
}

func (i Impl) GetTeamUsers(ctx context.Context, teamName string) ([]entity.User, error) {
	const query = `
		SELECT user_id, username, team_name, is_active
		FROM users
		WHERE team_name = $1;`
	rows, err := i.db.Query(ctx, query, teamName)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entity.ErrTeamNotFound
		}
		return nil, err
	}
	defer rows.Close()

	var users []entity.User
	for rows.Next() {
		var user entity.User
		err = rows.Scan(&user.UserId, &user.Username, &user.TeamName, &user.IsActive)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	if len(users) == 0 {
		return nil, entity.ErrTeamNotFound
	}
	return users, nil
}

func (i Impl) CheckTeamExists(ctx context.Context, teamName string) (bool, error) {
	const query = `
		SELECT 1
		FROM users
		WHERE team_name = $1
		LIMIT 1;`
	row := i.db.QueryRow(ctx, query, teamName)
	var exists int
	err := row.Scan(&exists)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (i Impl) GetTeamUsersByAuthorID(ctx context.Context, authorID string) ([]entity.User, error) {
	const query = `
		SELECT u.user_id, u.username, u.team_name, u.is_active
		FROM users u
		JOIN users a ON u.team_name = a.team_name
		WHERE a.user_id = $1;`
	rows, err := i.db.Query(ctx, query, authorID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entity.ErrUserNotFound
		}
		return nil, err
	}
	defer rows.Close()

	var users []entity.User
	for rows.Next() {
		var user entity.User
		err = rows.Scan(&user.UserId, &user.Username, &user.TeamName, &user.IsActive)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	if len(users) == 0 {
		return nil, entity.ErrUserNotFound
	}
	return users, nil
}
