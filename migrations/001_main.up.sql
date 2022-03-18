CREATE SCHEMA IF NOT EXISTS main;

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS orders(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    details TEXT NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);

INSERT INTO orders (details, created_at, updated_at) VALUES ('detail', '2015-01-10 00:51:14', '2015-01-10 00:51:14');
