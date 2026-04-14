-- name: GetPlayerRankByTierDivision :one
SELECT
  id,
  tier,
  division,
  created_at,
  updated_at
FROM player_ranks
WHERE tier = UPPER(BTRIM(sqlc.arg(tier)))
  AND division = UPPER(BTRIM(sqlc.arg(division)));
