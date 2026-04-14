-- name: DeletePlayerCurrentRanksByPlayerID :exec
DELETE FROM player_current_rank
WHERE player_id = sqlc.arg(player_id);

-- name: UpsertPlayerCurrentRank :exec
INSERT INTO player_current_rank (
  player_id,
  queue_type,
  player_ranks_id,
  league_points,
  wins,
  losses,
  recorded_at
) VALUES (
  sqlc.arg(player_id),
  sqlc.arg(queue_type),
  sqlc.arg(player_ranks_id),
  sqlc.arg(league_points),
  sqlc.arg(wins),
  sqlc.arg(losses),
  sqlc.arg(recorded_at)
)
ON CONFLICT (player_id, queue_type) DO UPDATE SET
  player_ranks_id = EXCLUDED.player_ranks_id,
  league_points = EXCLUDED.league_points,
  wins = EXCLUDED.wins,
  losses = EXCLUDED.losses,
  recorded_at = EXCLUDED.recorded_at,
  updated_at = NOW();

-- name: ListPlayerCurrentRanksByPlayerID :many
SELECT
  pcr.player_id,
  pcr.queue_type,
  pcr.player_ranks_id,
  pr.tier,
  pr.division,
  pcr.league_points,
  pcr.wins,
  pcr.losses,
  pcr.recorded_at,
  pcr.created_at AS player_current_rank_created_at,
  pcr.updated_at AS player_current_rank_updated_at
FROM player_current_rank AS pcr
JOIN player_ranks AS pr ON pr.id = pcr.player_ranks_id
WHERE pcr.player_id = sqlc.arg(player_id)
ORDER BY
  CASE pcr.queue_type
    WHEN 'RANKED_SOLO_5x5' THEN 1
    WHEN 'RANKED_FLEX_SR' THEN 2
    ELSE 3
  END,
  pcr.queue_type;
