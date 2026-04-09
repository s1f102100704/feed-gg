package riot

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const (
	matchIDsByPUUIDPathTemplate = "/lol/match/v5/matches/by-puuid/%s/ids"
	matchPathTemplate           = "/lol/match/v5/matches/%s"
	recentMatchCount            = 20
)

func (c *Client) fetchRecentMatchSummaries(
	ctx context.Context,
	host string,
	puuid string,
) ([]MatchSummary, int, error) {
	matchIDs, statusCode, err := c.fetchMatchIDs(ctx, host, puuid, recentMatchCount)
	if err != nil {
		return nil, statusCode, err
	}
	if len(matchIDs) == 0 {
		return []MatchSummary{}, http.StatusOK, nil
	}

	summaries := make([]MatchSummary, 0, len(matchIDs))
	for _, matchID := range matchIDs {
		match, statusCode, err := c.fetchMatch(ctx, host, matchID)
		if err != nil {
			return nil, statusCode, err
		}

		summary, err := mapMatchSummary(match, puuid)
		if err != nil {
			return nil, http.StatusBadGateway, err
		}

		summaries = append(summaries, summary)
	}

	return summaries, http.StatusOK, nil
}

func (c *Client) fetchMatchIDs(
	ctx context.Context,
	host string,
	puuid string,
	count int,
) ([]string, int, error) {
	path := fmt.Sprintf(matchIDsByPUUIDPathTemplate, url.PathEscape(puuid))

	query := url.Values{}
	query.Set("start", "0")
	query.Set("count", fmt.Sprintf("%d", count))

	var matchIDs []string
	statusCode, err := c.doJSON(ctx, host+path+"?"+query.Encode(), &matchIDs)
	if err != nil {
		return nil, statusCode, err
	}

	return matchIDs, statusCode, nil
}

func (c *Client) fetchMatch(ctx context.Context, host string, matchID string) (*MatchDTO, int, error) {
	path := fmt.Sprintf(matchPathTemplate, url.PathEscape(matchID))

	var match MatchDTO
	statusCode, err := c.doJSON(ctx, host+path, &match)
	if err != nil {
		return nil, statusCode, err
	}

	return &match, statusCode, nil
}

func mapMatchSummary(match *MatchDTO, puuid string) (MatchSummary, error) {
	participant, ok := findParticipant(match.Info.Participants, puuid)
	if !ok {
		return MatchSummary{}, fmt.Errorf("participant %s not found in match %s", puuid, match.Metadata.MatchID)
	}

	return MatchSummary{
		MatchID:          match.Metadata.MatchID,
		PlayedAt:         matchPlayedAt(match.Info),
		GameVersion:      match.Info.GameVersion,
		GameMode:         match.Info.GameMode,
		QueueID:          match.Info.QueueID,
		ChampionName:     participant.ChampionName,
		Role:             matchRole(participant),
		Win:              participant.Win,
		Kills:            participant.Kills,
		Deaths:           participant.Deaths,
		Assists:          participant.Assists,
		SummonerSpell1ID: participant.Summoner1ID,
		SummonerSpell2ID: participant.Summoner2ID,
		DurationSeconds:  matchDurationSeconds(match.Info),
		Participants:     mapMatchParticipants(match.Info.Participants),
	}, nil
}

func findParticipant(participants []MatchParticipantDTO, puuid string) (MatchParticipantDTO, bool) {
	for _, participant := range participants {
		if participant.PUUID == puuid {
			return participant, true
		}
	}

	return MatchParticipantDTO{}, false
}

func mapMatchParticipants(participants []MatchParticipantDTO) []MatchParticipantSummary {
	summaries := make([]MatchParticipantSummary, 0, len(participants))
	for _, participant := range participants {
		summaries = append(summaries, MatchParticipantSummary{
			PUUID:            participant.PUUID,
			GameName:         participant.RiotIDGameName,
			TagLine:          participant.RiotIDTagline,
			ChampionName:     participant.ChampionName,
			Role:             matchRole(participant),
			Win:              participant.Win,
			Kills:            participant.Kills,
			Deaths:           participant.Deaths,
			Assists:          participant.Assists,
			SummonerSpell1ID: participant.Summoner1ID,
			SummonerSpell2ID: participant.Summoner2ID,
		})
	}

	return summaries
}

func matchRole(participant MatchParticipantDTO) string {
	switch {
	case strings.TrimSpace(participant.TeamPosition) != "":
		return participant.TeamPosition
	case strings.TrimSpace(participant.IndividualPosition) != "":
		return participant.IndividualPosition
	case strings.TrimSpace(participant.Role) != "":
		return participant.Role
	default:
		return "UNKNOWN"
	}
}

func matchPlayedAt(info MatchInfoDTO) int64 {
	switch {
	case info.GameEndTimestamp > 0:
		return info.GameEndTimestamp
	case info.GameStartTimestamp > 0:
		return info.GameStartTimestamp
	default:
		return info.GameCreation
	}
}

func matchDurationSeconds(info MatchInfoDTO) int64 {
	if info.GameStartTimestamp > 0 && info.GameEndTimestamp > info.GameStartTimestamp {
		return (info.GameEndTimestamp - info.GameStartTimestamp) / 1000
	}
	if info.GameEndTimestamp == 0 && info.GameDuration >= 1000 {
		return info.GameDuration / 1000
	}
	if info.GameDuration > 0 {
		return info.GameDuration
	}
	return 0
}
