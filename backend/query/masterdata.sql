-- name: ListRegions :many
SELECT
  id,
  name,
  created_at,
  updated_at
FROM region
ORDER BY id;

-- name: RegionExists :one
SELECT EXISTS (
  SELECT 1
  FROM region
  WHERE name = UPPER(BTRIM(sqlc.arg(name)))
) AS exists;

-- name: GetRegionByName :one
SELECT
  id,
  name,
  created_at,
  updated_at
FROM region
WHERE name = UPPER(BTRIM(sqlc.arg(name)));

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
