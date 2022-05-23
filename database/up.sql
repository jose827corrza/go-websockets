DROP TABLE IF EXISTS users;

CREATE TABLE users (
    id VARCHAR(32) PRIMARY KEY,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    username VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

DROP TABLE IF EXISTS posts;

CREATE TABLE posts(
    id VARCHAR(32) PRIMARY KEY,
    post_content VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    user_id VARCHAR(32) NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
)

-- tag = go-rest-db

-- Estando parado en la carpeta database..
-- docker build . -t go-rest-db

-- para luego correr el docker es..
-- docker run -p 54321:5432 go-rest-db
-- donde es: PC:Postgres

-- Modificar el .env