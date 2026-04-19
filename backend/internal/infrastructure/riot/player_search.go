package riot

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const (
	accountPathTemplate  = "/riot/account/v1/accounts/by-riot-id/%s/%s"
	summonerPathTemplate = "/lol/summoner/v4/summoners/by-puuid/%s"
)

func (c *Client) SearchPlayerByRiotID(
	ctx context.Context,
	platformRegion string,
	gameName string,
	tagLine string,
) (*PlayerProfile, int, error) {
	regionalHost, platformHost, normalizedRegion, err := routeHosts(platformRegion)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	account, statusCode, err := c.fetchAccount(ctx, regionalHost, gameName, tagLine)
	if err != nil {
		return nil, statusCode, err
	}

	summoner, statusCode, err := c.fetchSummoner(ctx, platformHost, account.PUUID)
	if err != nil {
		return nil, statusCode, err
	}

	leagueEntries, statusCode, err := c.fetchRankEntries(ctx, platformHost, account.PUUID)
	if err != nil {
		return nil, statusCode, err
	}
	soloRank, flexRank := mapRanksByQueue(leagueEntries)

	matches, statusCode, err := c.fetchRecentMatchSummaries(ctx, regionalHost, account.PUUID)
	if err != nil {
		return nil, statusCode, err
	}

	iconURL := ""
	resolvedIconURL, err := c.ProfileIconURL(ctx, summoner.ProfileIconID)
	if err == nil {
		iconURL = resolvedIconURL
	}

	return &PlayerProfile{
		Region:         normalizedRegion,
		PUUID:          account.PUUID,
		GameName:       account.GameName,
		TagLine:        account.TagLine,
		SummonerLevel:  summoner.SummonerLevel,
		ProfileIconID:  summoner.ProfileIconID,
		ProfileIconURL: iconURL,
		RevisionDate:   summoner.RevisionDate,
		SoloRank:       soloRank,
		FlexRank:       flexRank,
		Matches:        matches,
	}, http.StatusOK, nil
}

func (c *Client) fetchAccount(ctx context.Context, host string, gameName string, tagLine string) (*Account, int, error) {
	path := fmt.Sprintf(
		accountPathTemplate,
		url.PathEscape(strings.TrimSpace(gameName)),
		url.PathEscape(strings.TrimSpace(tagLine)),
	)

	var account Account
	statusCode, err := c.doJSON(ctx, host+path, &account)
	if err != nil {
		return nil, statusCode, err
	}
	return &account, statusCode, nil
}

func (c *Client) fetchSummoner(ctx context.Context, host string, puuid string) (*Summoner, int, error) {
	path := fmt.Sprintf(summonerPathTemplate, url.PathEscape(puuid))

	var summoner Summoner
	statusCode, err := c.doJSON(ctx, host+path, &summoner)
	if err != nil {
		return nil, statusCode, err
	}
	return &summoner, statusCode, nil
}
