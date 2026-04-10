package riot

type Account struct {
	PUUID    string `json:"puuid"`
	GameName string `json:"gameName"`
	TagLine  string `json:"tagLine"`
}

type Summoner struct {
	ProfileIconID int    `json:"profileIconId"`
	RevisionDate  int64  `json:"revisionDate"`
	PUUID         string `json:"puuid"`
	SummonerLevel int64  `json:"summonerLevel"`
}

type RankedQueue struct {
	Tier         string `json:"tier"`
	Rank         string `json:"rank"`
	LeaguePoints int    `json:"leaguePoints"`
	Wins         int    `json:"wins"`
	Losses       int    `json:"losses"`
}

type LeagueEntry struct {
	QueueType    string `json:"queueType"`
	Tier         string `json:"tier"`
	Rank         string `json:"rank"`
	LeaguePoints int    `json:"leaguePoints"`
	Wins         int    `json:"wins"`
	Losses       int    `json:"losses"`
}

type MatchDTO struct {
	Metadata MatchMetadataDTO `json:"metadata"`
	Info     MatchInfoDTO     `json:"info"`
}

type MatchMetadataDTO struct {
	MatchID string `json:"matchId"`
}

type MatchInfoDTO struct {
	GameCreation       int64                 `json:"gameCreation"`
	GameStartTimestamp int64                 `json:"gameStartTimestamp"`
	GameEndTimestamp   int64                 `json:"gameEndTimestamp"`
	GameDuration       int64                 `json:"gameDuration"`
	GameVersion        string                `json:"gameVersion"`
	GameMode           string                `json:"gameMode"`
	QueueID            int                   `json:"queueId"`
	Participants       []MatchParticipantDTO `json:"participants"`
}

type MatchParticipantDTO struct {
	RiotIDGameName     string `json:"riotIdGameName"`
	RiotIDTagline      string `json:"riotIdTagline"`
	PUUID              string `json:"puuid"`
	ChampionName       string `json:"championName"`
	TeamPosition       string `json:"teamPosition"`
	IndividualPosition string `json:"individualPosition"`
	Role               string `json:"role"`
	Kills              int    `json:"kills"`
	Deaths             int    `json:"deaths"`
	Assists            int    `json:"assists"`
	Win                bool   `json:"win"`
	Summoner1ID        int    `json:"summoner1Id"`
	Summoner2ID        int    `json:"summoner2Id"`
}

type MatchParticipantSummary struct {
	PUUID            string `json:"puuid"`
	GameName         string `json:"gameName"`
	TagLine          string `json:"tagLine"`
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
	MatchID          string                    `json:"matchId"`
	PlayedAt         int64                     `json:"playedAt"`
	GameVersion      string                    `json:"gameVersion"`
	GameMode         string                    `json:"gameMode"`
	QueueID          int                       `json:"queueId"`
	ChampionName     string                    `json:"championName"`
	Role             string                    `json:"role"`
	Win              bool                      `json:"win"`
	Kills            int                       `json:"kills"`
	Deaths           int                       `json:"deaths"`
	Assists          int                       `json:"assists"`
	SummonerSpell1ID int                       `json:"summonerSpell1Id"`
	SummonerSpell2ID int                       `json:"summonerSpell2Id"`
	DurationSeconds  int64                     `json:"durationSeconds"`
	Participants     []MatchParticipantSummary `json:"participants,omitempty"`
}

type PlayerProfile struct {
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

type riotErrorResponse struct {
	Status struct {
		StatusCode int    `json:"status_code"`
		Message    string `json:"message"`
	} `json:"status"`
}
