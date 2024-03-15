CREATE TABLE IF NOT EXISTS actors
(
    id         SERIAL PRIMARY KEY,
    name       VARCHAR NOT NULL,
    gender     int     NOT NULL,
    birth_date DATE    NOT NULL
);

CREATE TABLE IF NOT EXISTS movies
(
    id           SERIAL PRIMARY KEY,
    title        VARCHAR(150)  NOT NULL CHECK (LENGTH(title) BETWEEN 1 AND 150),
    description  VARCHAR(1000) NOT NULL CHECK (LENGTH(description) < 1000),
    release_date DATE          NOT NULL,
    rating       DECIMAL(2, 1) NOT NULL CHECK (rating BETWEEN 0 AND 10)
);

CREATE TABLE IF NOT EXISTS movie_actors
(
    id       SERIAL PRIMARY KEY,
    movie_id INT NOT NULL,
    actor_id INT NOT NULL,
    FOREIGN KEY (movie_id) REFERENCES movies (id) ON DELETE CASCADE,
    FOREIGN KEY (actor_id) REFERENCES actors (id) ON DELETE CASCADE
);