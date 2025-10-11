-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS equipment(
    id SERIAL PRIMARY KEY,
    equipment_name VARCHAR(50) NOT NULL, 
    manufacturer VARCHAR(255),
    description VARCHAR(500) NOT NULL,
    image_url VARCHAR(255) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS equipment;
-- +goose StatementEnd
