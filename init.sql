DROP TABLE IF EXISTS votes;
DROP TABLE IF EXISTS posts;
DROP TABLE IF EXISTS threads;
DROP TABLE IF EXISTS forums;
DROP TABLE IF EXISTS users;

CREATE EXTENSION IF NOT EXISTS CITEXT;

CREATE TABLE users
(
  id  SERIAL PRIMARY KEY,
  about TEXT,
  email CITEXT UNIQUE NOT NULL,
  fullname VARCHAR(256) NOT NULL,
  nickname CITEXT UNIQUE
);

CREATE TABLE forums
(
  id SERIAL PRIMARY KEY,
  posts BIGINT NOT NULL DEFAULT 0,
  slug CITEXT UNIQUE NOT NULL,
  threads INT NOT NULL DEFAULT 0,
  title TEXT NOT NULL,
  author CITEXT NOT NULL REFERENCES users (nickname)
);

CREATE TABLE threads
(
  author CITEXT NOT NULL REFERENCES users (nickname),
  created TIMESTAMPTZ DEFAULT NOW(),
  forum CITEXT,
  id SERIAL PRIMARY KEY,
  message TEXT NOT NULL,
  slug CITEXT UNIQUE,
  title TEXT,
  votes INT NOT NULL DEFAULT 0
);

CREATE TABLE posts
(
  id SERIAL PRIMARY KEY,
  author CITEXT NOT NULL REFERENCES users (nickname),
  created TIMESTAMPTZ DEFAULT NOW(),
  forum CITEXT,
  isEdited BOOLEAN DEFAULT FALSE,
  message TEXT NOT NULL,
  parent INT DEFAULT 0,
  tid INT NOT NULL REFERENCES threads (id),
  slug VARCHAR(32)
);

-- DROP TRIGGER IF EXISTS init_post_slug ON posts;
--
--
-- CREATE OR REPLACE FUNCTION init_post() RETURNS TRIGGER AS
-- $init_post_slug$
-- BEGIN
--   UPDATE posts
-- --   SET slug = concat(cast(new.parent AS TEXT), cast(posts.id AS TEXT))
--   SET slug = cast(parent AS TEXT)
--   WHERE id = TG_RELID;
--   RETURN new;
-- END;
-- $init_post_slug$
-- LANGUAGE plpgsql;
--
-- CREATE TRIGGER init_post_slug BEFORE INSERT ON posts FOR EACH ROW EXECUTE PROCEDURE init_post();

CREATE TABLE votes
(
  nickname CITEXT NOT NULL REFERENCES users (nickname),
  tid INT NOT NULL,
  voice INT NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS votes_ind ON votes (nickname, tid);

DROP TRIGGER IF EXISTS vote_insertion ON votes;
DROP TRIGGER IF EXISTS vote_updating ON votes;

CREATE OR REPLACE FUNCTION insert_vote() RETURNS TRIGGER AS
$vote_insertion$
BEGIN
  UPDATE threads
  SET votes = votes + new.voice
  WHERE id = new.tid;
  RETURN NEW;
END;
$vote_insertion$
LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION update_vote() RETURNS TRIGGER AS
$vote_updating$
BEGIN
  UPDATE threads
    SET votes = votes - OLD.voice + NEW.voice
    WHERE id = new.tid;
  RETURN NEW;
END;
$vote_updating$
LANGUAGE plpgsql;

CREATE TRIGGER vote_updating BEFORE UPDATE ON votes FOR EACH ROW EXECUTE PROCEDURE update_vote();
CREATE TRIGGER vote_insertion BEFORE INSERT ON votes FOR EACH ROW EXECUTE PROCEDURE insert_vote();