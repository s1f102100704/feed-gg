-- name: UpsertMatchParticipant :exec
INSERT INTO match_participant (
  match_history_id,
  player_id,
  game_name_snapshot,
  tag_line_snapshot,
  team_id,
  champion_name,
  team_position,
  role,
  win,
  kills,
  deaths,
  assists,
  summoner_spell1_id,
  summoner_spell2_id
) VALUES (
  sqlc.arg(match_history_id),
  sqlc.arg(player_id),
  sqlc.arg(game_name_snapshot),
  sqlc.arg(tag_line_snapshot),
  sqlc.arg(team_id),
  sqlc.arg(champion_name),
  sqlc.arg(team_position),
  sqlc.arg(role),
  sqlc.arg(win),
  sqlc.arg(kills),
  sqlc.arg(deaths),
  sqlc.arg(assists),
  sqlc.arg(summoner_spell1_id),
  sqlc.arg(summoner_spell2_id)
)
ON CONFLICT (match_history_id, player_id) DO UPDATE SET
  game_name_snapshot = EXCLUDED.game_name_snapshot,
  tag_line_snapshot = EXCLUDED.tag_line_snapshot,
  team_id = EXCLUDED.team_id,
  champion_name = EXCLUDED.champion_name,
  team_position = EXCLUDED.team_position,
  role = EXCLUDED.role,
  win = EXCLUDED.win,
  kills = EXCLUDED.kills,
  deaths = EXCLUDED.deaths,
  assists = EXCLUDED.assists,
  summoner_spell1_id = EXCLUDED.summoner_spell1_id,
  summoner_spell2_id = EXCLUDED.summoner_spell2_id,
  updated_at = NOW();

-- name: ListMatchParticipantsByMatchHistoryID :many
SELECT
  mp.match_history_id,
  mp.player_id,
  p.puuid,
  mp.game_name_snapshot,
  mp.tag_line_snapshot,
  mp.team_id,
  mp.champion_name,
  mp.team_position,
  mp.role,
  mp.win,
  mp.kills,
  mp.deaths,
  mp.assists,
  mp.summoner_spell1_id,
  mp.summoner_spell2_id,
  mp.created_at AS match_participant_created_at,
  mp.updated_at AS match_participant_updated_at
FROM match_participant AS mp
JOIN player AS p ON p.id = mp.player_id
WHERE mp.match_history_id = sqlc.arg(match_history_id)
ORDER BY mp.team_id, mp.player_id;
