package riot

import (
	"context"
	"fmt"
	"net/url"
)

const leagueEntriesByPuuidPathTemplate = "/lol/league/v4/entries/by-puuid/%s"

func (c *Client) fetchRankEntries(ctx context.Context, host string, puuid string) ([]LeagueEntry, int, error) {
	path := fmt.Sprintf(leagueEntriesByPuuidPathTemplate, url.PathEscape(puuid))

	var entries []LeagueEntry
	statusCode, err := c.doJSON(ctx, host+path, &entries)
	if err != nil {
		return nil, statusCode, err
	}

	return entries, statusCode, nil
}

func mapRanksByQueue(entries []LeagueEntry) (soloRank *RankedQueue, flexRank *RankedQueue) {
	for _, entry := range entries {
		rankedQueue := &RankedQueue{
			Tier:         entry.Tier,
			Rank:         entry.Rank,
			LeaguePoints: entry.LeaguePoints,
			Wins:         entry.Wins,
			Losses:       entry.Losses,
		}

		switch entry.QueueType {
		case "RANKED_SOLO_5x5":
			soloRank = rankedQueue
		case "RANKED_FLEX_SR":
			flexRank = rankedQueue
		}
	}

	return soloRank, flexRank
}
