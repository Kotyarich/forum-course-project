DROP TABLE IF EXISTS forum_users;
DROP TABLE IF EXISTS forums;

CREATE EXTENSION IF NOT EXISTS CITEXT;

CREATE TABLE forums
(
  id      SERIAL PRIMARY KEY,
  posts   BIGINT        NOT NULL DEFAULT 0,
  slug    CITEXT UNIQUE NOT NULL,
  threads INT           NOT NULL DEFAULT 0,
  title   TEXT          NOT NULL,
  author  CITEXT        NOT NULL REFERENCES users (nickname)
);

CREATE INDEX IF NOT EXISTS forums_slug ON forums USING hash (slug);

CREATE TABLE forum_users
(
  fUser CITEXT COLLATE ucs_basic NOT NULL,
  forum CITEXT                   NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS f_users_ind ON forum_users (forum, fUser);
CREATE INDEX IF NOT EXISTS forum_users_fuser ON forum_users (fUser);
