package riot

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

const (
	ddragonVersionsURL = "https://ddragon.leagueoflegends.com/api/versions.json"
	ddragonVersionTTL  = 24 * time.Hour
)

func (c *Client) ddragonLatestVersion(ctx context.Context) (string, error) {
	c.versionMu.RLock()
	if c.ddragonVersion != "" && time.Since(c.ddragonFetchedAt) < ddragonVersionTTL {
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
	defer closeBody(resp.Body)

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
	c.ddragonFetchedAt = time.Now()
	c.versionMu.Unlock()

	return versions[0], nil
}
