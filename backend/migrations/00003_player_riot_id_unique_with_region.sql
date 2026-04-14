-- +goose Up
DROP INDEX IF EXISTS player_game_name_tag_line_unique;

CREATE UNIQUE INDEX player_region_game_name_tag_line_unique
  ON player(region_id, game_name, tag_line)
  WHERE game_name IS NOT NULL AND tag_line IS NOT NULL;

-- +goose Down
DROP INDEX IF EXISTS player_region_game_name_tag_line_unique;

CREATE UNIQUE INDEX player_game_name_tag_line_unique
  ON player(game_name, tag_line)
  WHERE game_name IS NOT NULL AND tag_line IS NOT NULL;
