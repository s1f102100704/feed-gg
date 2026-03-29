-- +goose Up
CREATE TABLE region (
  id SMALLSERIAL PRIMARY KEY,
  name VARCHAR(20) NOT NULL UNIQUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE player_ranks (
  id SMALLSERIAL PRIMARY KEY,
  rank_division VARCHAR(20) NOT NULL UNIQUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE tag (
  id SMALLSERIAL PRIMARY KEY,
  name VARCHAR(50) NOT NULL UNIQUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE match_history (
  id SERIAL PRIMARY KEY,
  played_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE season(
  id SMALLSERIAL PRIMARY KEY,
  name varchar(30) NOT NULL UNIQUE,
  start_date DATE NOT NULL,
  end_date DATE NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE player (
  id BIGSERIAL PRIMARY KEY,
  puuid VARCHAR(100) NOT NULL,
  game_name VARCHAR(30) NOT NULL,
  tag_line  varchar(10) NOT NULL,
  region_id SMALLINT NOT NULL REFERENCES region(id),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE (game_name, tag_line)
);

CREATE TABLE player_current_rank (
  player_id BIGINT NOT NULL REFERENCES player(id),
  queue_type VARCHAR(30) NOT NULL,
  player_ranks_id SMALLINT REFERENCES player_ranks(id),
  league_points INTEGER,
  wins INTEGER,
  losses INTEGER,
  recorded_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  PRIMARY KEY (player_id, queue_type),
  CHECK (
    (player_ranks_id IS NULL AND league_points IS NULL AND wins IS NULL AND losses IS NULL) OR
    (player_ranks_id IS NOT NULL AND league_points IS NOT NULL AND wins IS NOT NULL AND losses IS NOT NULL)
  )
);

CREATE TABLE player_rank_history (
  id BIGSERIAL PRIMARY KEY,
  player_id BIGINT NOT NULL REFERENCES player(id),
  player_ranks_id SMALLINT NOT NULL REFERENCES player_ranks(id),
  league_points INTEGER NOT NULL,
  season_id SMALLINT NOT NULL REFERENCES season(id),
  recorded_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE player_match (
  player_id BIGINT NOT NULL REFERENCES player(id),
  match_history_id INTEGER NOT NULL REFERENCES match_history(id),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  PRIMARY KEY (player_id, match_history_id)
);

CREATE TABLE player_tag (
  player_id BIGINT NOT NULL REFERENCES player(id),
  tag_id SMALLINT NOT NULL REFERENCES tag(id),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  PRIMARY KEY (player_id, tag_id)
);


CREATE  UNIQUE INDEX player_puuid ON player(puuid);
CREATE INDEX player_region_id ON player(region_id);
CREATE INDEX player_current_rank_player_ranks_id ON player_current_rank(player_ranks_id);
CREATE INDEX player_rank_history_player_id ON player_rank_history(player_id);
CREATE INDEX player_rank_history_player_ranks_id ON player_rank_history(player_ranks_id);
CREATE INDEX player_match_match_history_id ON player_match(match_history_id);
