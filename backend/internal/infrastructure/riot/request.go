package riot

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

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
