DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS users;

CREATE EXTENSION IF NOT EXISTS CITEXT;

CREATE TABLE users
(
  id       SERIAL PRIMARY KEY,
  about    TEXT,
  email    CITEXT UNIQUE NOT NULL,
  fullname VARCHAR(256)  NOT NULL,
  nickname CITEXT UNIQUE NOT NULL,
  password TEXT,
  is_admin BOOLEAN DEFAULT FALSE
);

CREATE INDEX IF NOT EXISTS users_nickname_and_email ON users (nickname, email);
