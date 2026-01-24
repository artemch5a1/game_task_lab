-- SQLite schema for games and genres
-- One genre -> many games

PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS genres (
  id TEXT PRIMARY KEY,
  title TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS games (
  id TEXT PRIMARY KEY,
  title TEXT NOT NULL,
  description TEXT NOT NULL,
  release_date TEXT NOT NULL, -- RFC3339Nano
  genre_id TEXT NOT NULL,
  FOREIGN KEY (genre_id) REFERENCES genres(id) ON UPDATE CASCADE ON DELETE RESTRICT
);

CREATE INDEX IF NOT EXISTS idx_games_genre_id ON games(genre_id);
CREATE INDEX IF NOT EXISTS idx_games_release_date ON games(release_date);

CREATE TABLE IF NOT EXISTS users (
  id TEXT PRIMARY KEY,
  username TEXT NOT NULL UNIQUE,
  password TEXT NOT NULL,
  user_role TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);

