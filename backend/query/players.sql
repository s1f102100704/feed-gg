-- name: GetSavedPlayerKeyByRiotID :one
SELECT
  p.id AS player_id,
  p.puuid
FROM player AS p
WHERE p.region_id = sqlc.arg(region_id)
  AND p.game_name = sqlc.arg(game_name)
  AND p.tag_line = sqlc.arg(tag_line)
LIMIT 1;

-- name: GetSavedPlayerByPuuid :one
SELECT
  p.id AS player_id,
  p.puuid,
  p.game_name,
  p.tag_line,
  p.region_id,
  r.name AS region_name,
  p.profile_icon_id,
  p.summoner_level,
  p.revision_date,
  p.last_synced_at,
  p.created_at AS player_created_at,
  p.updated_at AS player_updated_at
FROM player AS p
JOIN region AS r ON r.id = p.region_id
WHERE p.puuid = sqlc.arg(puuid)
LIMIT 1;

-- name: UpsertPlayerProfile :one
INSERT INTO player (
  puuid,
  game_name,
  tag_line,
  region_id,
  profile_icon_id,
  summoner_level,
  revision_date,
  last_synced_at
) VALUES (
  sqlc.arg(puuid),
  sqlc.arg(game_name),
  sqlc.arg(tag_line),
  sqlc.arg(region_id),
  sqlc.arg(profile_icon_id),
  sqlc.arg(summoner_level),
  sqlc.arg(revision_date),
  NOW()
)
ON CONFLICT (puuid) DO UPDATE SET
  game_name = EXCLUDED.game_name,
  tag_line = EXCLUDED.tag_line,
  region_id = EXCLUDED.region_id,
  profile_icon_id = EXCLUDED.profile_icon_id,
  summoner_level = EXCLUDED.summoner_level,
  revision_date = EXCLUDED.revision_date,
  last_synced_at = EXCLUDED.last_synced_at,
  updated_at = NOW()
RETURNING *;

-- name: UpsertParticipantPlayer :one
INSERT INTO player (
  puuid,
  game_name,
  tag_line,
  region_id
) VALUES (
  sqlc.arg(puuid),
  sqlc.arg(game_name),
  sqlc.arg(tag_line),
  sqlc.arg(region_id)
)
ON CONFLICT (puuid) DO UPDATE SET
  game_name = COALESCE(EXCLUDED.game_name, player.game_name),
  tag_line = COALESCE(EXCLUDED.tag_line, player.tag_line),
  region_id = EXCLUDED.region_id,
  updated_at = NOW()
RETURNING *;
