package entity

import "errors"

type Team struct {
	TeamName string
	Members  []User
}

var (
	ErrTeamAlreadyExists = errors.New("team already exists")
	ErrTeamNotFound      = errors.New("team not found")
)
