DROP TABLE IF EXISTS posts_creations;
DROP TABLE IF EXISTS users_creations;
DROP TABLE IF EXISTS threads_creations;
DROP TABLE IF EXISTS votes_creations;
DROP TABLE IF EXISTS forums_creations;

CREATE TABLE posts_creations
(
  id         SERIAL PRIMARY KEY,
  created_at TIMESTAMP
);

CREATE TABLE users_creations
(
  id         text unique ,
  created_at TIMESTAMP
);

CREATE TABLE threads_creations
(
  id         SERIAL PRIMARY KEY,
  created_at TIMESTAMP
);

CREATE TABLE votes_creations
(
  id         text unique,
  created_at TIMESTAMP
);

CREATE TABLE forums_creations
(
  id         SERIAL PRIMARY KEY,
  created_at TIMESTAMP
);


