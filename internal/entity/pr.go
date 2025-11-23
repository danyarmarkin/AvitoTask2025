package entity

import (
	"errors"
	"time"
)

type PullRequestStatus string

const (
	PullRequestStatusMERGED PullRequestStatus = "MERGED"
	PullRequestStatusOPEN   PullRequestStatus = "OPEN"
)

type PullRequest struct {
	AssignedReviewers []string
	AuthorId          string
	PullRequestId     string
	PullRequestName   string
	Status            PullRequestStatus

	ReplacedByReviewer string

	CreatedAt *time.Time `json:"createdAt"`
	MergedAt  *time.Time `json:"mergedAt"`
}

var (
	ErrAuthorNotFound  = errors.New("author not found")
	ErrPRAlreadyExists = errors.New("pull request already exists")
	ErrPRMerged        = errors.New("pull request already merged")
	ErrPrNotFound      = errors.New("pull request not found")
)
