package masterdata

import (
	"context"

	db "feed-gg/backend/internal/infrastructure/db/sqlc"
)

type RegionStore struct {
	queries *db.Queries
}

func NewRegionStore(queries *db.Queries) *RegionStore {
	return &RegionStore{queries: queries}
}

func (s *RegionStore) RegionExists(ctx context.Context, name string) (bool, error) {
	return s.queries.RegionExists(ctx, name)
}

func (s *RegionStore) ListRegionNames(ctx context.Context) ([]string, error) {
	regions, err := s.queries.ListRegions(ctx)
	if err != nil {
		return nil, err
	}

	names := make([]string, 0, len(regions))
	for _, region := range regions {
		names = append(names, region.Name)
	}

	return names, nil
}
