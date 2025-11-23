-- +goose Up
CREATE TABLE pr
(
    pull_request_id    TEXT PRIMARY KEY,
    pull_request_name  TEXT                    NOT NULL,
    author_id          TEXT                    NOT NULL,
    assigned_reviewers TEXT[]    DEFAULT '{}'  NOT NULL,
    status             TEXT                    NOT NULL,
    created_at         TIMESTAMP DEFAULT now() NOT NULL,
    merged_at          TIMESTAMP,
    CONSTRAINT assigned_reviewers_max_len CHECK (coalesce(cardinality(assigned_reviewers), 0) <= 2),
    FOREIGN KEY (author_id) REFERENCES users (user_id) ON DELETE CASCADE

);


CREATE INDEX idx_pr_assigned_reviewers_gin ON pr USING GIN (assigned_reviewers);

-- +goose Down
DROP INDEX idx_pr_assigned_reviewers_gin;
DROP TABLE pr;