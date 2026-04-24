-- name: ListLabels :many
SELECT
  id,
  name,
  created_at,
  updated_at
FROM label
ORDER BY id;

-- name: GetLabelByID :one
SELECT
  id,
  name,
  created_at,
  updated_at
FROM label
WHERE id = sqlc.arg(id)
LIMIT 1;

-- name: GetPlayerLabelVoteByPuuidAndVoterKey :one
SELECT
  plv.id,
  plv.player_id,
  plv.label_id,
  plv.voter_key,
  plv.created_at,
  plv.updated_at
FROM player_label_votes AS plv
JOIN player AS p ON p.id = plv.player_id
WHERE p.puuid = sqlc.arg(puuid)
  AND plv.voter_key = sqlc.arg(voter_key)
LIMIT 1;

-- name: ListPlayerLabelVoteSummariesByPuuid :many
SELECT
  l.id,
  l.name,
  COUNT(*)::bigint AS vote_count
FROM player_label_votes AS plv
JOIN player AS p ON p.id = plv.player_id
JOIN label AS l ON l.id = plv.label_id
WHERE p.puuid = sqlc.arg(puuid)
GROUP BY l.id, l.name
ORDER BY vote_count DESC, l.id ASC;

-- name: ListTopPlayerLabelVoteSummariesByPuuid :many
SELECT
  l.id,
  l.name,
  COUNT(*)::bigint AS vote_count
FROM player_label_votes AS plv
JOIN player AS p ON p.id = plv.player_id
JOIN label AS l ON l.id = plv.label_id
WHERE p.puuid = sqlc.arg(puuid)
GROUP BY l.id, l.name
ORDER BY vote_count DESC, l.id ASC
LIMIT sqlc.arg(limit_count);

-- name: CountPlayerLabelVotesByPuuid :one
SELECT
  COUNT(*)::bigint AS vote_count
FROM player_label_votes AS plv
JOIN player AS p ON p.id = plv.player_id
WHERE p.puuid = sqlc.arg(puuid);

-- name: UpsertPlayerLabelVoteByPuuid :one
INSERT INTO player_label_votes (
  player_id,
  label_id,
  voter_key
)
SELECT
  p.id,
  sqlc.arg(label_id),
  sqlc.arg(voter_key)
FROM player AS p
WHERE p.puuid = sqlc.arg(puuid)
ON CONFLICT (player_id, voter_key) DO UPDATE SET
  label_id = EXCLUDED.label_id,
  updated_at = NOW()
RETURNING
  id,
  player_id,
  label_id,
  voter_key,
  created_at,
  updated_at;
