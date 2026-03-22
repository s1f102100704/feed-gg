CREATE TABLE region (
  id SMALLSERIAL PRIMARY KEY,
  name TEXT NOT NULL UNIQUE,
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

CREATE TABLE player (
  id BIGSERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  tag_line TEXT NOT NULL,
  current_rank_id BIGINT NOT NULL REFERENCES player_ranks(id),
  current_league_points INTEGER NOT NULL,
  region_id BIGINT NOT NULL REFERENCES region(id),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE (name, tag_line)
);

CREATE TABLE player_rank_history (
  id BIGSERIAL PRIMARY KEY,
  player_id BIGINT NOT NULL REFERENCES player(id),
  player_ranks_id BIGINT NOT NULL REFERENCES player_ranks(id),
  league_points INTEGER NOT NULL,
  season_id SMALLINT NOT NULL REFERENCES season(id),
  recorded_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE player_match (
  player_profile_id BIGINT NOT NULL REFERENCES player_profile(id),
  match_history_id BIGINT NOT NULL REFERENCES match_history(id),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  PRIMARY KEY (player_profile_id, match_history_id)
);

CREATE TABLE player_tag (
  player_id BIGINT NOT NULL REFERENCES player_profile(id),
  tag_id BIGINT NOT NULL REFERENCES tag(id),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  PRIMARY KEY (player_id, tag_id)
);

CREATE TABLE patch (
  id SMALLSERIAL PRIMARY KEY,
  version VARCHAR(20) NOT NULL UNIQUE,
  release_date DATE NOT NULL,
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
)

CREATE INDEX player_profile_region_id ON player_profile(region_id);
CREATE INDEX player_rank_history_player_profile_id ON player_rank_history(player_profile_id);
CREATE INDEX player_match_match_history_id ON player_match(match_history_id);
