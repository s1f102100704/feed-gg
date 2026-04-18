package cache

import "feed-gg/backend/internal/usecase"

type NoopPlayerSearchCache struct{}

func NewNoopPlayerSearchCache() *NoopPlayerSearchCache {
	return &NoopPlayerSearchCache{}
}

func (c *NoopPlayerSearchCache) Get(key usecase.PlayerSearchKey) (*usecase.PlayerSearchResult, bool) {
	return nil, false
}

func (c *NoopPlayerSearchCache) Set(key usecase.PlayerSearchKey, value *usecase.PlayerSearchResult) {}

var _ usecase.PlayerSearchCache = (*NoopPlayerSearchCache)(nil)
