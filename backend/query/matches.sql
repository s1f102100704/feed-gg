-- name: UpsertMatchHistory :one
INSERT INTO match_history (
  match_id,
  region_id,
  queue_id,
  game_mode,
  game_version,
  duration_seconds,
  played_at
) VALUES (
  sqlc.arg(match_id),
  sqlc.arg(region_id),
  sqlc.arg(queue_id),
  sqlc.arg(game_mode),
  sqlc.arg(game_version),
  sqlc.arg(duration_seconds),
  sqlc.arg(played_at)
)
ON CONFLICT (match_id) DO UPDATE SET
  region_id = EXCLUDED.region_id,
  queue_id = EXCLUDED.queue_id,
  game_mode = EXCLUDED.game_mode,
  game_version = EXCLUDED.game_version,
  duration_seconds = EXCLUDED.duration_seconds,
  played_at = EXCLUDED.played_at,
  updated_at = NOW()
RETURNING *;

-- name: ListRecentMatchHistoriesByPlayerID :many
SELECT
  mh.id AS match_history_id,
  mh.match_id,
  mh.region_id,
  mh.queue_id,
  mh.game_mode,
  mh.game_version,
  mh.duration_seconds,
  mh.played_at,
  mh.created_at AS match_history_created_at,
  mh.updated_at AS match_history_updated_at,
  mp.player_id,
  mp.team_id,
  mp.champion_name,
  mp.team_position,
  mp.role,
  mp.win,
  mp.kills,
  mp.deaths,
  mp.assists,
  mp.summoner_spell1_id,
  mp.summoner_spell2_id
FROM match_participant AS mp
JOIN match_history AS mh ON mh.id = mp.match_history_id
WHERE mp.player_id = sqlc.arg(player_id)
ORDER BY mh.played_at DESC
LIMIT sqlc.arg(limit_count);
