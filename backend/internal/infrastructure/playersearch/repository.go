package playersearch

import (
	"context"

	"feed-gg/backend/internal/infrastructure/riot"
	"feed-gg/backend/internal/usecase"
)

type NoopRepository struct{}

func NewNoopRepository() *NoopRepository {
	return &NoopRepository{}
}

func (r *NoopRepository) FindSavedPlayer(
	ctx context.Context,
	input usecase.PlayerSearchInput,
) (*usecase.PlayerSearchResult, error) {
	return nil, nil
}

func (r *NoopRepository) SaveFetchedPlayer(ctx context.Context, fetched *riot.PlayerProfile) error {
	return nil
}

var _ usecase.PlayerSearchRepository = (*NoopRepository)(nil)
