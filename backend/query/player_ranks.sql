-- name: GetPlayerRankByTierDivision :one
SELECT
  id,
  tier,
  division,
  created_at,
  updated_at
FROM player_ranks
WHERE tier = sqlc.arg(tier)
  AND division = sqlc.arg(division);
