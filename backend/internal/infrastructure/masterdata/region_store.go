package masterdata

import (
	"context"
	"database/sql"
	"strings"
)

type RegionStore struct {
	db *sql.DB
}

func NewRegionStore(db *sql.DB) *RegionStore {
	return &RegionStore{db: db}
}

func (s *RegionStore) RegionExists(ctx context.Context, name string) (bool, error) {
	const query = `
		SELECT EXISTS(
			SELECT 1
			FROM region
			WHERE name = $1
		)
	`

	var exists bool
	err := s.db.QueryRowContext(ctx, query, normalizeRegionName(name)).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (s *RegionStore) ListRegionNames(ctx context.Context) ([]string, error) {
	const query = `
		SELECT name
		FROM region
		ORDER BY id
	`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	regions := make([]string, 0)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		regions = append(regions, name)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return regions, nil
}

func normalizeRegionName(name string) string {
	return strings.ToUpper(strings.TrimSpace(name))
}
