-- +goose Up
ALTER TABLE tag RENAME TO label;
ALTER TABLE label RENAME CONSTRAINT tag_pkey TO label_pkey;
ALTER TABLE label RENAME CONSTRAINT tag_name_key TO label_name_key;

ALTER SEQUENCE tag_id_seq RENAME TO label_id_seq;

ALTER TABLE player_tag RENAME TO player_label;
ALTER TABLE player_label RENAME COLUMN tag_id TO label_id;
ALTER TABLE player_label RENAME CONSTRAINT player_tag_pkey TO player_label_pkey;
ALTER TABLE player_label RENAME CONSTRAINT player_tag_player_id_fkey TO player_label_player_id_fkey;
ALTER TABLE player_label RENAME CONSTRAINT player_tag_tag_id_fkey TO player_label_label_id_fkey;

-- +goose Down
ALTER TABLE player_label RENAME CONSTRAINT player_label_label_id_fkey TO player_tag_tag_id_fkey;
ALTER TABLE player_label RENAME CONSTRAINT player_label_player_id_fkey TO player_tag_player_id_fkey;
ALTER TABLE player_label RENAME CONSTRAINT player_label_pkey TO player_tag_pkey;
ALTER TABLE player_label RENAME COLUMN label_id TO tag_id;
ALTER TABLE player_label RENAME TO player_tag;

ALTER SEQUENCE label_id_seq RENAME TO tag_id_seq;

ALTER TABLE label RENAME CONSTRAINT label_name_key TO tag_name_key;
ALTER TABLE label RENAME CONSTRAINT label_pkey TO tag_pkey;
ALTER TABLE label RENAME TO tag;
