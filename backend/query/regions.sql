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
