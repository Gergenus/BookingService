-- +goose Up
-- +goose StatementBegin
CREATE TYPE roles AS ENUM (
    'admin', 'scientist'
);

CREATE TABLE IF NOT EXISTS users(
    uid uuid PRIMARY KEY,
    username VARCHAR(50) NOT NULL, 
    role roles NOT NULL, 
    email VARCHAR(255) NOT NULL UNIQUE,
    hashed_password TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
DROP TYPE IF EXISTS roles;
-- +goose StatementEnd
