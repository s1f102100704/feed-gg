-- +goose Up
DROP TABLE IF EXISTS player_label;

CREATE TABLE player_label_votes (
  id BIGSERIAL PRIMARY KEY,
  player_id BIGINT NOT NULL,
  label_id SMALLINT NOT NULL,
  voter_key VARCHAR(100) NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT player_label_votes_player_id_fkey FOREIGN KEY (player_id) REFERENCES player(id) ON DELETE CASCADE,
  CONSTRAINT player_label_votes_label_id_fkey FOREIGN KEY (label_id) REFERENCES label(id),
  CONSTRAINT player_label_votes_player_id_voter_key_key UNIQUE (player_id, voter_key)
);

CREATE INDEX player_label_votes_player_id_idx ON player_label_votes(player_id);
CREATE INDEX player_label_votes_label_id_idx ON player_label_votes(label_id);
CREATE INDEX player_label_votes_label_id_player_id_idx ON player_label_votes(label_id, player_id);

-- +goose Down
DROP INDEX IF EXISTS player_label_votes_label_id_player_id_idx;
DROP INDEX IF EXISTS player_label_votes_label_id_idx;
DROP INDEX IF EXISTS player_label_votes_player_id_idx;

DROP TABLE IF EXISTS player_label_votes;

CREATE TABLE player_label (
  player_id BIGINT NOT NULL REFERENCES player(id),
  label_id SMALLINT NOT NULL REFERENCES label(id),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  PRIMARY KEY (player_id, label_id)
);
