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

type PlayerProfile struct {
	Region         string `json:"region"`
	PUUID          string `json:"puuid"`
	GameName       string `json:"gameName"`
	TagLine        string `json:"tagLine"`
	SummonerLevel  int64  `json:"summonerLevel"`
	ProfileIconID  int    `json:"profileIconId"`
	ProfileIconURL string `json:"profileIconUrl"`
	RevisionDate   int64  `json:"revisionDate"`
}

type riotErrorResponse struct {
	Status struct {
		StatusCode int    `json:"status_code"`
		Message    string `json:"message"`
	} `json:"status"`
}
