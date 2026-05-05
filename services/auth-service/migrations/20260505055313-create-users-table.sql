
-- +migrate Up
CREATE TABLE users (
                       id VARCHAR(36) PRIMARY KEY,
                       email VARCHAR(255) NOT NULL,
                       password_hash VARCHAR(255) NOT NULL,
                       role VARCHAR(50) NOT NULL,
                       created_at DATETIME NOT NULL,
                       updated_at DATETIME NOT NULL,
                       CONSTRAINT uk_users_email UNIQUE (email)
);
-- Index
CREATE INDEX idx_users_email_lookup ON users(email);

-- +migrate Down
DROP INDEX idx_users_email_lookup ON users;
DROP TABLE users;