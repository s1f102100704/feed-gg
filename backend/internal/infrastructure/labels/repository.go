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
	queries *db.Queries
}

func NewRepository(queries *db.Queries) *Repository {
	return &Repository{queries: queries}
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

func (r *Repository) ListPlayerLabels(
	ctx context.Context,
	input usecase.PlayerLabelsInput,
) (*usecase.PlayerLabelsResult, error) {
	if r == nil || r.queries == nil {
		return nil, ErrRepositoryNotConfigured
	}

	if _, err := r.queries.GetSavedPlayerByPuuid(ctx, input.PUUID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, usecase.ErrPlayerLabelsNotFound
		}
		return nil, fmt.Errorf("get player %q: %w", input.PUUID, err)
	}

	return r.playerLabelsResult(ctx, input.PUUID)
}

func (r *Repository) SavePlayerLabelVote(
	ctx context.Context,
	input usecase.PlayerLabelVoteInput,
) (*usecase.PlayerLabelVoteResult, error) {
	if r == nil || r.queries == nil {
		return nil, ErrRepositoryNotConfigured
	}

	label, err := r.queries.GetLabelByID(ctx, input.LabelID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, usecase.ErrPlayerLabelsNotFound
		}
		return nil, fmt.Errorf("get label %d: %w", input.LabelID, err)
	}

	vote, err := r.queries.UpsertPlayerLabelVoteByPuuid(ctx, db.UpsertPlayerLabelVoteByPuuidParams{
		LabelID:  input.LabelID,
		VoterKey: input.VoterKey,
		Puuid:    input.PUUID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, usecase.ErrPlayerLabelsNotFound
		}
		return nil, fmt.Errorf("upsert player label vote: %w", err)
	}

	result, err := r.playerLabelsResult(ctx, input.PUUID)
	if err != nil {
		return nil, err
	}

	selected := usecase.PlayerLabelSummary{
		ID:   vote.LabelID,
		Name: label.Name,
	}
	for _, summary := range result.Labels {
		if summary.ID == selected.ID {
			selected.VoteCount = summary.VoteCount
			break
		}
	}

	return &usecase.PlayerLabelVoteResult{
		SelectedLabel: selected,
		Labels:        result.Labels,
		TotalVotes:    result.TotalVotes,
	}, nil
}

func (r *Repository) playerLabelsResult(
	ctx context.Context,
	puuid string,
) (*usecase.PlayerLabelsResult, error) {
	rows, err := r.queries.ListPlayerLabelVoteSummariesByPuuid(ctx, puuid)
	if err != nil {
		return nil, fmt.Errorf("list player label summaries: %w", err)
	}

	labels := make([]usecase.PlayerLabelSummary, 0, len(rows))
	for _, row := range rows {
		labels = append(labels, usecase.PlayerLabelSummary{
			ID:        row.ID,
			Name:      row.Name,
			VoteCount: row.VoteCount,
		})
	}

	totalVotes, err := r.queries.CountPlayerLabelVotesByPuuid(ctx, puuid)
	if err != nil {
		return nil, fmt.Errorf("count player label votes: %w", err)
	}

	return &usecase.PlayerLabelsResult{
		Labels:     labels,
		TotalVotes: totalVotes,
	}, nil
}

var _ usecase.PlayerLabelsRepository = (*Repository)(nil)
