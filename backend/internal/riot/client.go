package riot

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

const (
	accountPathTemplate   = "/riot/account/v1/accounts/by-riot-id/%s/%s"
	summonerPathTemplate  = "/lol/summoner/v4/summoners/by-puuid/%s"
	ddragonVersionsURL    = "https://ddragon.leagueoflegends.com/api/versions.json"
	defaultRequestTimeout = 10 * time.Second
)

var (
	ErrInvalidRegion = errors.New("invalid region")
)

type Client struct {
	apiKey         string
	httpClient     *http.Client
	versionMu      sync.RWMutex
	ddragonVersion string
}

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

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: defaultRequestTimeout,
		},
	}
}

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

	iconURL := ""
	version, err := c.ddragonLatestVersion(ctx)
	if err == nil && version != "" {
		iconURL = fmt.Sprintf(
			"https://ddragon.leagueoflegends.com/cdn/%s/img/profileicon/%d.png",
			version,
			summoner.ProfileIconID,
		)
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

func (c *Client) ddragonLatestVersion(ctx context.Context) (string, error) {
	c.versionMu.RLock()
	if c.ddragonVersion != "" {
		version := c.ddragonVersion
		c.versionMu.RUnlock()
		return version, nil
	}
	c.versionMu.RUnlock()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ddragonVersionsURL, nil)
	if err != nil {
		return "", err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ddragon returned status %d", resp.StatusCode)
	}

	var versions []string
	if err := json.NewDecoder(resp.Body).Decode(&versions); err != nil {
		return "", err
	}
	if len(versions) == 0 {
		return "", errors.New("ddragon returned no versions")
	}

	c.versionMu.Lock()
	c.ddragonVersion = versions[0]
	c.versionMu.Unlock()

	return versions[0], nil
}

func (c *Client) doJSON(ctx context.Context, rawURL string, dest any) (int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	req.Header.Set("X-Riot-Token", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return http.StatusBadGateway, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		var riotErr riotErrorResponse
		if err := json.Unmarshal(body, &riotErr); err == nil && riotErr.Status.Message != "" {
			return resp.StatusCode, errors.New(riotErr.Status.Message)
		}
		if len(body) > 0 {
			return resp.StatusCode, fmt.Errorf("riot api returned status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
		}
		return resp.StatusCode, fmt.Errorf("riot api returned status %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(dest); err != nil {
		return http.StatusBadGateway, err
	}

	return resp.StatusCode, nil
}

func routeHosts(platformRegion string) (regionalHost string, platformHost string, normalizedRegion string, err error) {
	region := strings.ToUpper(strings.TrimSpace(platformRegion))

	platformToRegional := map[string]string{
		"BR1":  "americas",
		"LA1":  "americas",
		"LA2":  "americas",
		"NA1":  "americas",
		"JP1":  "asia",
		"KR":   "asia",
		"EUN1": "europe",
		"EUW1": "europe",
		"ME1":  "europe",
		"RU":   "europe",
		"TR1":  "europe",
		"OC1":  "sea",
		"PH2":  "sea",
		"SG2":  "sea",
		"TH2":  "sea",
		"TW2":  "sea",
		"VN2":  "sea",
	}

	regionalKey, ok := platformToRegional[region]
	if !ok {
		return "", "", "", ErrInvalidRegion
	}

	platformHost = fmt.Sprintf("https://%s.api.riotgames.com", strings.ToLower(region))
	regionalHost = fmt.Sprintf("https://%s.api.riotgames.com", regionalKey)
	return regionalHost, platformHost, region, nil
}
