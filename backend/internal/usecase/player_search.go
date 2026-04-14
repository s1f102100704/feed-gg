package usecase

import (
	"context"
	"errors"
	"net/http"

	"feed-gg/backend/internal/infrastructure/riot"
)

var (
	ErrUnsupportedRegion       = errors.New("unsupported region")
	ErrRegionMasterUnavailable = errors.New("failed to load region master")
)

type PlayerSearch struct {
	riotClient    PlayerSearcher
	regionChecker RegionChecker
}

type PlayerSearcher interface {
	SearchPlayerByRiotID(
		ctx context.Context,
		platformRegion string,
		gameName string,
		tagLine string,
	) (*riot.PlayerProfile, int, error)
}

type RegionChecker interface {
	RegionExists(ctx context.Context, name string) (bool, error)
}

type PlayerSearchInput struct {
	Region   string
	GameName string
	TagLine  string
}

type PlayerSearchResult struct {
	Region         string         `json:"region"`
	PUUID          string         `json:"puuid"`
	GameName       string         `json:"gameName"`
	TagLine        string         `json:"tagLine"`
	SummonerLevel  int64          `json:"summonerLevel"`
	ProfileIconID  int            `json:"profileIconId"`
	ProfileIconURL string         `json:"profileIconUrl"`
	RevisionDate   int64          `json:"revisionDate"`
	SoloRank       *RankedQueue   `json:"soloRank,omitempty"`
	FlexRank       *RankedQueue   `json:"flexRank,omitempty"`
	Matches        []MatchSummary `json:"matches,omitempty"`
}

type RankedQueue struct {
	Tier         string `json:"tier"`
	Rank         string `json:"rank"`
	LeaguePoints int    `json:"leaguePoints"`
	Wins         int    `json:"wins"`
	Losses       int    `json:"losses"`
}

type MatchParticipant struct {
	PUUID            string `json:"puuid"`
	GameName         string `json:"gameName"`
	TagLine          string `json:"tagLine"`
	TeamID           int    `json:"teamId"`
	ChampionName     string `json:"championName"`
	Role             string `json:"role"`
	Win              bool   `json:"win"`
	Kills            int    `json:"kills"`
	Deaths           int    `json:"deaths"`
	Assists          int    `json:"assists"`
	SummonerSpell1ID int    `json:"summonerSpell1Id"`
	SummonerSpell2ID int    `json:"summonerSpell2Id"`
}

type MatchSummary struct {
	MatchID          string             `json:"matchId"`
	PlayedAt         int64              `json:"playedAt"`
	GameVersion      string             `json:"gameVersion"`
	GameMode         string             `json:"gameMode"`
	QueueID          int                `json:"queueId"`
	TeamID           int                `json:"teamId"`
	ChampionName     string             `json:"championName"`
	Role             string             `json:"role"`
	Win              bool               `json:"win"`
	Kills            int                `json:"kills"`
	Deaths           int                `json:"deaths"`
	Assists          int                `json:"assists"`
	SummonerSpell1ID int                `json:"summonerSpell1Id"`
	SummonerSpell2ID int                `json:"summonerSpell2Id"`
	DurationSeconds  int64              `json:"durationSeconds"`
	Participants     []MatchParticipant `json:"participants,omitempty"`
}

func NewPlayerSearch(riotClient PlayerSearcher, regionChecker RegionChecker) *PlayerSearch {
	return &PlayerSearch{
		riotClient:    riotClient,
		regionChecker: regionChecker,
	}
}

func (u *PlayerSearch) Execute(
	ctx context.Context,
	input PlayerSearchInput,
) (*PlayerSearchResult, int, error) {
	exists, err := u.regionChecker.RegionExists(ctx, input.Region)
	if err != nil {
		return nil, http.StatusInternalServerError, ErrRegionMasterUnavailable
	}
	if !exists {
		return nil, http.StatusBadRequest, ErrUnsupportedRegion
	}

	player, statusCode, err := u.riotClient.SearchPlayerByRiotID(
		ctx,
		input.Region,
		input.GameName,
		input.TagLine,
	)
	if err != nil {
		if errors.Is(err, riot.ErrInvalidRegion) {
			return nil, http.StatusBadRequest, ErrUnsupportedRegion
		}
		return nil, statusCode, err
	}

	return mapPlayerSearchResult(player), http.StatusOK, nil
}

func mapPlayerSearchResult(player *riot.PlayerProfile) *PlayerSearchResult {
	if player == nil {
		return nil
	}

	return &PlayerSearchResult{
		Region:         player.Region,
		PUUID:          player.PUUID,
		GameName:       player.GameName,
		TagLine:        player.TagLine,
		SummonerLevel:  player.SummonerLevel,
		ProfileIconID:  player.ProfileIconID,
		ProfileIconURL: player.ProfileIconURL,
		RevisionDate:   player.RevisionDate,
		SoloRank:       mapRankedQueue(player.SoloRank),
		FlexRank:       mapRankedQueue(player.FlexRank),
		Matches:        mapMatchSummaries(player.Matches),
	}
}

func mapRankedQueue(rank *riot.RankedQueue) *RankedQueue {
	if rank == nil {
		return nil
	}

	return &RankedQueue{
		Tier:         rank.Tier,
		Rank:         rank.Rank,
		LeaguePoints: rank.LeaguePoints,
		Wins:         rank.Wins,
		Losses:       rank.Losses,
	}
}

func mapMatchSummaries(matches []riot.MatchSummary) []MatchSummary {
	summaries := make([]MatchSummary, 0, len(matches))
	for _, match := range matches {
		summaries = append(summaries, MatchSummary{
			MatchID:          match.MatchID,
			PlayedAt:         match.PlayedAt,
			GameVersion:      match.GameVersion,
			GameMode:         match.GameMode,
			QueueID:          match.QueueID,
			TeamID:           match.TeamID,
			ChampionName:     match.ChampionName,
			Role:             match.Role,
			Win:              match.Win,
			Kills:            match.Kills,
			Deaths:           match.Deaths,
			Assists:          match.Assists,
			SummonerSpell1ID: match.SummonerSpell1ID,
			SummonerSpell2ID: match.SummonerSpell2ID,
			DurationSeconds:  match.DurationSeconds,
			Participants:     mapMatchParticipants(match.Participants),
		})
	}

	return summaries
}

func mapMatchParticipants(participants []riot.MatchParticipantSummary) []MatchParticipant {
	result := make([]MatchParticipant, 0, len(participants))
	for _, participant := range participants {
		result = append(result, MatchParticipant{
			PUUID:            participant.PUUID,
			GameName:         participant.GameName,
			TagLine:          participant.TagLine,
			TeamID:           participant.TeamID,
			ChampionName:     participant.ChampionName,
			Role:             participant.Role,
			Win:              participant.Win,
			Kills:            participant.Kills,
			Deaths:           participant.Deaths,
			Assists:          participant.Assists,
			SummonerSpell1ID: participant.SummonerSpell1ID,
			SummonerSpell2ID: participant.SummonerSpell2ID,
		})
	}

	return result
}
