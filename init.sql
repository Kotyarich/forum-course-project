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
  slug INTEGER[] NOT NULL,
  rootId INT
);

DROP TRIGGER IF EXISTS add_root_id ON posts;

CREATE OR REPLACE FUNCTION init_post() RETURNS TRIGGER AS
$add_root_id$
BEGIN
  IF new.parent = 0 THEN
    UPDATE posts
    SET rootId = new.id
    WHERE id = new.id;
  ELSE
    UPDATE posts
    SET rootId = (SELECT rootId FROM posts WHERE id = new.parent)
    WHERE id = NEW.id;
  END IF;
  RETURN new;
END;
$add_root_id$
LANGUAGE plpgsql;

CREATE TRIGGER add_root_id AFTER INSERT ON posts FOR EACH ROW EXECUTE PROCEDURE init_post();

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