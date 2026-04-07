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

type PlayerProfile struct {
	Region         string       `json:"region"`
	PUUID          string       `json:"puuid"`
	GameName       string       `json:"gameName"`
	TagLine        string       `json:"tagLine"`
	SummonerLevel  int64        `json:"summonerLevel"`
	ProfileIconID  int          `json:"profileIconId"`
	ProfileIconURL string       `json:"profileIconUrl"`
	RevisionDate   int64        `json:"revisionDate"`
	SoloRank       *RankedQueue `json:"soloRank,omitempty"`
	FlexRank       *RankedQueue `json:"flexRank,omitempty"`
}

type riotErrorResponse struct {
	Status struct {
		StatusCode int    `json:"status_code"`
		Message    string `json:"message"`
	} `json:"status"`
}
