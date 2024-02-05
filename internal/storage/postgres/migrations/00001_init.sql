-- +goose Up

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS orders (
       id SERIAL PRIMARY KEY,
       user_id varchar,
       good varchar,
       quantity INTEGER,
       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE orders;
DROP EXTENSION IF EXISTS "uuid-ossp";