package labels

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	db "feed-gg/backend/internal/infrastructure/db/sqlc"
	"feed-gg/backend/internal/usecase"
)

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

	return r.playerLabelsResult(ctx, r.queries, input.PUUID)
}

func (r *Repository) SavePlayerLabelVote(
	ctx context.Context,
	input usecase.PlayerLabelVoteInput,
) (result *usecase.PlayerLabelVoteResult, err error) {
	if r == nil || r.db == nil || r.queries == nil {
		return nil, ErrRepositoryNotConfigured
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin player label vote transaction: %w", err)
	}
	defer func() {
		if err == nil {
			return
		}
		_ = tx.Rollback()
	}()

	q := r.queries.WithTx(tx)
	label, err := q.GetLabelByID(ctx, input.LabelID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, usecase.ErrPlayerLabelsNotFound
		}
		return nil, fmt.Errorf("get label %d: %w", input.LabelID, err)
	}

	vote, err := q.UpsertPlayerLabelVoteByPuuid(ctx, db.UpsertPlayerLabelVoteByPuuidParams{
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

	labelsResult, err := r.playerLabelsResult(ctx, q, input.PUUID)
	if err != nil {
		return nil, err
	}

	selected := usecase.PlayerLabelSummary{
		ID:   vote.LabelID,
		Name: label.Name,
	}
	for _, summary := range labelsResult.Labels {
		if summary.ID == selected.ID {
			selected.VoteCount = summary.VoteCount
			break
		}
	}

	result = &usecase.PlayerLabelVoteResult{
		SelectedLabel: selected,
		Labels:        labelsResult.Labels,
		TotalVotes:    labelsResult.TotalVotes,
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit player label vote transaction: %w", err)
	}

	return result, nil
}

func (r *Repository) playerLabelsResult(
	ctx context.Context,
	q *db.Queries,
	puuid string,
) (*usecase.PlayerLabelsResult, error) {
	rows, err := q.ListPlayerLabelVoteSummariesByPuuid(ctx, puuid)
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

	totalVotes, err := q.CountPlayerLabelVotesByPuuid(ctx, puuid)
	if err != nil {
		return nil, fmt.Errorf("count player label votes: %w", err)
	}

	return &usecase.PlayerLabelsResult{
		Labels:     labels,
		TotalVotes: totalVotes,
	}, nil
}

var _ usecase.PlayerLabelsRepository = (*Repository)(nil)
