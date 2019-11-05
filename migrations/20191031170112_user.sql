-- +goose Up
CREATE TABLE users (
    id int NOT NULL PRIMARY KEY,
    name text
);
-- +goose Down
DROP TABLE users;