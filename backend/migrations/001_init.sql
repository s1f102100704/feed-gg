CREATE TABLE area_master (
  id BIGSERIAL PRIMARY KEY,
  name TEXT NOT NULL UNIQUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE player_rank_master (
  id BIGSERIAL PRIMARY KEY,
  rank_division TEXT NOT NULL UNIQUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE tag_master (
  id BIGSERIAL PRIMARY KEY,
  name TEXT NOT NULL UNIQUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE match_history (
  id BIGSERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  played_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE player_profile (
  id BIGSERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  tag TEXT NOT NULL,
  area_id BIGINT NOT NULL REFERENCES area_master(id),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE (name, tag)
);

CREATE TABLE player_rank_history (
  id BIGSERIAL PRIMARY KEY,
  player_profile_id BIGINT NOT NULL REFERENCES player_profile(id),
  player_rank_master_id BIGINT NOT NULL REFERENCES player_rank_master(id),
  league_points INTEGER NOT NULL,
  recorded_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE player_match (
  player_profile_id BIGINT NOT NULL REFERENCES player_profile(id),
  match_history_id BIGINT NOT NULL REFERENCES match_history(id),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  PRIMARY KEY (player_profile_id, match_history_id)
);

CREATE TABLE player_profile_tag (
  player_profile_id BIGINT NOT NULL REFERENCES player_profile(id),
  tag_id BIGINT NOT NULL REFERENCES tag_master(id),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  PRIMARY KEY (player_profile_id, tag_id)
);

CREATE INDEX idx_player_profile_area_id ON player_profile(area_id);
CREATE INDEX idx_player_rank_history_player_profile_id ON player_rank_history(player_profile_id);
CREATE INDEX idx_player_rank_history_player_rank_master_id ON player_rank_history(player_rank_master_id);
CREATE INDEX idx_player_match_match_history_id ON player_match(match_history_id);
CREATE INDEX idx_player_profile_tag_tag_id ON player_profile_tag(tag_id);
