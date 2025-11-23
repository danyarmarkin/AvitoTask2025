package repository

import (
	"AvitoTask2025/internal/entity"
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/lib/pq"
)

var _ PRRepository = (*Impl)(nil)

func (i *Impl) CreatePullRequest(ctx context.Context, pr entity.PullRequest) (entity.PullRequest, error) {
	var (
		prID, prName, authorID, status string
		createdAt, mergedAt            sql.NullTime
		assignedReviewers              pgtype.Array[string]
	)

	const query = `
		INSERT INTO pr (pull_request_id, pull_request_name, author_id, assigned_reviewers, status)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (pull_request_id) DO NOTHING
		RETURNING pull_request_id, pull_request_name, author_id, assigned_reviewers, status, created_at, merged_at
	`

	err := i.db.QueryRow(ctx, query,
		pr.PullRequestId,
		pr.PullRequestName,
		pr.AuthorId,
		pq.Array(pr.AssignedReviewers),
		string(pr.Status),
	).Scan(&prID, &prName, &authorID, &assignedReviewers, &status, &createdAt, &mergedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return entity.PullRequest{}, entity.ErrPRAlreadyExists
	}
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23503" {
			return entity.PullRequest{}, entity.ErrAuthorNotFound
		}
		return entity.PullRequest{}, err
	}

	pr.PullRequestId = prID
	pr.PullRequestName = prName
	pr.AuthorId = authorID
	pr.AssignedReviewers = assignedReviewers.Elements
	pr.Status = entity.PullRequestStatus(status)

	if createdAt.Valid {
		t := createdAt.Time
		pr.CreatedAt = &t
	} else {
		pr.CreatedAt = nil
	}
	if mergedAt.Valid {
		t := mergedAt.Time
		pr.MergedAt = &t
	} else {
		pr.MergedAt = nil
	}

	return pr, nil
}

// UpdateAssignedReviewers updates the assigned reviewers of a pull request.
// Returns entity.ErrPRMerged if the pull request is already merged.
func (i *Impl) UpdateAssignedReviewers(ctx context.Context, pr entity.PullRequest) (entity.PullRequest, error) {
	const query = `
			UPDATE pr
			SET assigned_reviewers = $2
			WHERE pull_request_id = $1 AND status != $3
			RETURNING pull_request_id, pull_request_name, author_id, assigned_reviewers, status, created_at, merged_at
		`
	var (
		prID, prName, authorID, status string
		createdAt, mergedAt            sql.NullTime
		assignedReviewers              pgtype.Array[string]
	)
	err := i.db.QueryRow(
		ctx,
		query,
		pr.PullRequestId,
		pq.Array(pr.AssignedReviewers),
		entity.PullRequestStatusMERGED,
	).Scan(&prID, &prName, &authorID, &assignedReviewers, &status, &createdAt, &mergedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return entity.PullRequest{}, entity.ErrPRMerged
	}
	if err != nil {
		return entity.PullRequest{}, err
	}
	pr.PullRequestId = prID
	pr.PullRequestName = prName
	pr.AuthorId = authorID
	pr.AssignedReviewers = assignedReviewers.Elements
	pr.Status = entity.PullRequestStatus(status)
	if createdAt.Valid {
		t := createdAt.Time
		pr.CreatedAt = &t
	} else {
		pr.CreatedAt = nil
	}
	if mergedAt.Valid {
		t := mergedAt.Time
		pr.MergedAt = &t
	} else {
		pr.MergedAt = nil
	}
	return pr, nil
}

func (i *Impl) Merge(ctx context.Context, pr entity.PullRequest) (entity.PullRequest, error) {
	const query = `
			UPDATE pr
			SET status = $2,
			    merged_at = CASE WHEN status = $3 THEN NOW() ELSE merged_at END
			WHERE pull_request_id = $1
			RETURNING pull_request_id, pull_request_name, author_id, assigned_reviewers, status, created_at, merged_at
		`
	var (
		prID, prName, authorID, status string
		createdAt, mergedAt            sql.NullTime
		assignedReviewers              pgtype.Array[string]
	)
	err := i.db.QueryRow(
		ctx,
		query,
		pr.PullRequestId,
		entity.PullRequestStatusMERGED,
		entity.PullRequestStatusOPEN,
	).Scan(&prID, &prName, &authorID, &assignedReviewers, &status, &createdAt, &mergedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return entity.PullRequest{}, entity.ErrPrNotFound
	}
	if err != nil {
		return entity.PullRequest{}, err
	}
	pr.PullRequestId = prID
	pr.PullRequestName = prName
	pr.AuthorId = authorID
	pr.AssignedReviewers = assignedReviewers.Elements
	pr.Status = entity.PullRequestStatus(status)
	if createdAt.Valid {
		t := createdAt.Time
		pr.CreatedAt = &t
	} else {
		pr.CreatedAt = nil
	}
	if mergedAt.Valid {
		t := mergedAt.Time
		pr.MergedAt = &t
	} else {
		pr.MergedAt = nil
	}
	return pr, nil
}

func (i *Impl) GetPRByID(ctx context.Context, prID string) (entity.PullRequest, error) {
	const query = `
			SELECT pull_request_id, pull_request_name, author_id, assigned_reviewers, status, created_at, merged_at
			FROM pr
			WHERE pull_request_id = $1
		`
	var (
		prName, authorID, status string
		createdAt, mergedAt      sql.NullTime
		assignedReviewers        pgtype.Array[string]
	)
	err := i.db.QueryRow(ctx, query, prID).Scan(&prID, &prName, &authorID, &assignedReviewers, &status, &createdAt, &mergedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return entity.PullRequest{}, entity.ErrPrNotFound
	}
	if err != nil {
		return entity.PullRequest{}, err
	}
	pr := entity.PullRequest{
		PullRequestId:     prID,
		PullRequestName:   prName,
		AuthorId:          authorID,
		AssignedReviewers: assignedReviewers.Elements,
		Status:            entity.PullRequestStatus(status),
	}
	if createdAt.Valid {
		t := createdAt.Time
		pr.CreatedAt = &t
	} else {
		pr.CreatedAt = nil
	}
	if mergedAt.Valid {
		t := mergedAt.Time
		pr.MergedAt = &t
	} else {
		pr.MergedAt = nil
	}
	return pr, nil
}

func (i *Impl) GetUsersPRs(ctx context.Context, userID string) ([]entity.PullRequest, error) {
	const query = ` 
SELECT pull_request_id, pull_request_name, author_id, status
FROM pr
WHERE assigned_reviewers @> ARRAY[$1]::text[];
`
	rows, err := i.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prs []entity.PullRequest
	for rows.Next() {
		var pr entity.PullRequest
		err := rows.Scan(&pr.PullRequestId, &pr.PullRequestName, &pr.AuthorId, &pr.Status)
		if err != nil {
			return nil, err
		}
		prs = append(prs, pr)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return prs, nil
}
