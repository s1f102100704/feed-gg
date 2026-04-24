package cache

import (
	"time"

	"feed-gg/backend/internal/usecase"

	gocache "github.com/patrickmn/go-cache"
)

const (
	labelsCacheKey            = "labels"
	defaultLabelsCacheTTL     = 7 * 24 * time.Hour
	defaultLabelsCacheCleanup = 12 * time.Hour
)

type LabelsCache struct {
	cache *gocache.Cache
}

func NewLabelsCache() *LabelsCache {
	return &LabelsCache{
		cache: gocache.New(defaultLabelsCacheTTL, defaultLabelsCacheCleanup),
	}
}

func (c *LabelsCache) Get() ([]usecase.Label, bool) {
	if c == nil || c.cache == nil {
		return nil, false
	}

	cached, found := c.cache.Get(labelsCacheKey)
	if !found {
		return nil, false
	}

	labels, ok := cached.([]usecase.Label)
	if !ok {
		return nil, false
	}

	return cloneLabels(labels), true
}

func (c *LabelsCache) Set(labels []usecase.Label) {
	if c == nil || c.cache == nil {
		return
	}

	c.cache.Set(labelsCacheKey, cloneLabels(labels), gocache.DefaultExpiration)
}

func (c *LabelsCache) Delete() {
	if c == nil || c.cache == nil {
		return
	}

	c.cache.Delete(labelsCacheKey)
}

func cloneLabels(labels []usecase.Label) []usecase.Label {
	if labels == nil {
		return nil
	}

	cloned := make([]usecase.Label, len(labels))
	copy(cloned, labels)
	return cloned
}

var _ usecase.LabelsCache = (*LabelsCache)(nil)
