package labels

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	db "feed-gg/backend/internal/infrastructure/db/sqlc"
	"feed-gg/backend/internal/usecase"
)

var ErrRepositoryNotConfigured = errors.New("labels repository is not configured")

type Repository struct {
	db      *sql.DB
	queries *db.Queries
}

func NewRepository(sqlDB *sql.DB, queries *db.Queries) *Repository {
	return &Repository{
		db:      sqlDB,
		queries: queries,
	}
}

func (r *Repository) ListLabels(ctx context.Context) ([]usecase.Label, error) {
	if r == nil || r.queries == nil {
		return nil, ErrRepositoryNotConfigured
	}

	rows, err := r.queries.ListLabels(ctx)
	if err != nil {
		return nil, fmt.Errorf("list labels: %w", err)
	}

	labels := make([]usecase.Label, 0, len(rows))
	for _, row := range rows {
		labels = append(labels, usecase.Label{
			ID:   row.ID,
			Name: row.Name,
		})
	}

	return labels, nil
}

var _ usecase.LabelsRepository = (*Repository)(nil)
