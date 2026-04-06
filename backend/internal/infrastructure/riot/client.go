package riot

import (
	"net/http"
	"sync"
	"time"
)

const (
	defaultRequestTimeout = 10 * time.Second
)

type Client struct {
	apiKey           string
	httpClient       *http.Client
	versionMu        sync.RWMutex
	ddragonVersion   string
	ddragonFetchedAt time.Time
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: defaultRequestTimeout,
		},
	}
}
