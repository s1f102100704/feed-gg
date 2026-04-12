-- +goose Up
CREATE TABLE region (
  id SMALLSERIAL PRIMARY KEY,
  name VARCHAR(20) NOT NULL UNIQUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE player_ranks (
  id SMALLSERIAL PRIMARY KEY,
  tier VARCHAR(20) NOT NULL,
  division VARCHAR(5) NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE (tier, division)
);

CREATE TABLE tag (
  id SMALLSERIAL PRIMARY KEY,
  name VARCHAR(50) NOT NULL UNIQUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE season (
  id SMALLSERIAL PRIMARY KEY,
  name VARCHAR(30) NOT NULL UNIQUE,
  start_date DATE NOT NULL,
  end_date DATE NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE player (
  id BIGSERIAL PRIMARY KEY,
  puuid VARCHAR(100) NOT NULL UNIQUE,
  game_name VARCHAR(50),
  tag_line VARCHAR(20),
  region_id SMALLINT NOT NULL REFERENCES region(id),
  profile_icon_id INTEGER,
  summoner_level BIGINT,
  revision_date TIMESTAMPTZ,
  last_synced_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
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
  queue_type VARCHAR(30) NOT NULL,
  player_ranks_id SMALLINT NOT NULL REFERENCES player_ranks(id),
  league_points INTEGER NOT NULL,
  season_id SMALLINT NOT NULL REFERENCES season(id),
  recorded_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE (player_id, season_id, queue_type)
);

CREATE TABLE match_history (
  id BIGSERIAL PRIMARY KEY,
  match_id VARCHAR(40) NOT NULL UNIQUE,
  region_id SMALLINT NOT NULL REFERENCES region(id),
  queue_id INTEGER NOT NULL,
  game_mode VARCHAR(30) NOT NULL,
  game_version VARCHAR(30) NOT NULL,
  duration_seconds INTEGER NOT NULL,
  played_at TIMESTAMPTZ NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE match_participant (
  match_history_id BIGINT NOT NULL REFERENCES match_history(id) ON DELETE CASCADE,
  player_id BIGINT NOT NULL REFERENCES player(id),
  game_name_snapshot VARCHAR(50),
  tag_line_snapshot VARCHAR(20),
  team_id SMALLINT NOT NULL,
  champion_name VARCHAR(50) NOT NULL,
  team_position VARCHAR(20),
  role VARCHAR(20),
  win BOOLEAN NOT NULL,
  kills INTEGER NOT NULL,
  deaths INTEGER NOT NULL,
  assists INTEGER NOT NULL,
  summoner_spell1_id INTEGER NOT NULL,
  summoner_spell2_id INTEGER NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  PRIMARY KEY (match_history_id, player_id)
);

CREATE TABLE player_tag (
  player_id BIGINT NOT NULL REFERENCES player(id),
  tag_id SMALLINT NOT NULL REFERENCES tag(id),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  PRIMARY KEY (player_id, tag_id)
);

CREATE UNIQUE INDEX player_game_name_tag_line_unique
  ON player(game_name, tag_line)
  WHERE game_name IS NOT NULL AND tag_line IS NOT NULL;

CREATE INDEX player_region_id ON player(region_id);
CREATE INDEX player_current_rank_player_ranks_id ON player_current_rank(player_ranks_id);
CREATE INDEX player_rank_history_player_id ON player_rank_history(player_id);
CREATE INDEX player_rank_history_player_ranks_id ON player_rank_history(player_ranks_id);
CREATE INDEX player_rank_history_queue_type ON player_rank_history(queue_type);
CREATE INDEX match_history_region_id ON match_history(region_id);
CREATE INDEX match_history_played_at ON match_history(played_at DESC);
CREATE INDEX match_participant_player_id ON match_participant(player_id);
