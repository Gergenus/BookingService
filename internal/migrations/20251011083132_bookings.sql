-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS booking(
    id SERIAL PRIMARY KEY,
    equipment_id int REFERENCES equipment(id) ON DELETE CASCADE,
    user_id uuid REFERENCES users(uid) ON DELETE CASCADE,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS booking;
-- +goose StatementEnd
