package entity

import "errors"

type User struct {
	IsActive bool
	TeamName string
	UserId   string
	Username string
}

var ErrUserNotFound = errors.New("user not found")
