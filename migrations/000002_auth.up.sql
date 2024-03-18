CREATE TYPE role AS ENUM ('admin', 'user');

CREATE TABLE IF NOT EXISTS users
(
    id       SERIAL PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    role     role NOT NULL DEFAULT 'user'
);
