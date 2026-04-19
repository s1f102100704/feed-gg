package cache

import (
	"time"

	"feed-gg/backend/internal/usecase"

	gocache "github.com/patrickmn/go-cache"
)

const (
	defaultPlayerSearchCacheTTL     = 5 * time.Minute
	defaultPlayerSearchCacheCleanup = 10 * time.Minute
)

type PlayerSearchCache struct {
	cache *gocache.Cache
}

func NewPlayerSearchCache() *PlayerSearchCache {
	return &PlayerSearchCache{
		cache: gocache.New(defaultPlayerSearchCacheTTL, defaultPlayerSearchCacheCleanup),
	}
}

func (c *PlayerSearchCache) Get(key usecase.PlayerSearchKey) (*usecase.PlayerSearchResult, bool) {
	if c == nil || c.cache == nil {
		return nil, false
	}

	cached, found := c.cache.Get(string(key))
	if !found {
		return nil, false
	}

	result, ok := cached.(*usecase.PlayerSearchResult)
	if !ok {
		return nil, false
	}

	return result, true
}

func (c *PlayerSearchCache) Set(key usecase.PlayerSearchKey, value *usecase.PlayerSearchResult) {
	if c == nil || c.cache == nil || value == nil {
		return
	}

	c.cache.Set(string(key), value, gocache.DefaultExpiration)
}

var _ usecase.PlayerSearchCache = (*PlayerSearchCache)(nil)
