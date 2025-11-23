-- +goose Up
CREATE TABLE users (
    user_id TEXT PRIMARY KEY,
    username TEXT NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    team_name TEXT
);

-- +goose Down
DROP TABLE users;