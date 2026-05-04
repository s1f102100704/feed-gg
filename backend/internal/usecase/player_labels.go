package usecase

import (
	"context"
	"errors"
	"net/http"
	"strings"
)

const maxPlayerLabelVoterKeyLength = 100

var (
	ErrInvalidPlayerLabelsInput    = errors.New("puuid is required")
	ErrInvalidPlayerLabelVoteInput = errors.New("puuid, labelId, and voterKey are required")
	ErrPlayerLabelsNotFound        = errors.New("player or label not found")
	ErrPlayerLabelsLoadFailed      = errors.New("failed to load player labels")
	ErrPlayerLabelVoteFailed       = errors.New("failed to save player label vote")
)

type PlayerLabelSummary struct {
	ID        int16  `json:"id"`
	Name      string `json:"name"`
	VoteCount int64  `json:"voteCount"`
}

type PlayerLabelsInput struct {
	PUUID string
}

type PlayerLabelVoteInput struct {
	PUUID    string
	LabelID  int16
	VoterKey string
}

type PlayerLabelsResult struct {
	Labels     []PlayerLabelSummary `json:"labels"`
	TotalVotes int64                `json:"totalVotes"`
}

type PlayerLabelVoteResult struct {
	SelectedLabel PlayerLabelSummary   `json:"selectedLabel"`
	Labels        []PlayerLabelSummary `json:"labels"`
	TotalVotes    int64                `json:"totalVotes"`
}

type PlayerLabelsRepository interface {
	ListPlayerLabels(ctx context.Context, input PlayerLabelsInput) (*PlayerLabelsResult, error)
	SavePlayerLabelVote(ctx context.Context, input PlayerLabelVoteInput) (*PlayerLabelVoteResult, error)
}

type PlayerLabels struct {
	repository PlayerLabelsRepository
}

func NewPlayerLabels(repository PlayerLabelsRepository) *PlayerLabels {
	return &PlayerLabels{repository: repository}
}

func (u *PlayerLabels) List(
	ctx context.Context,
	input PlayerLabelsInput,
) (*PlayerLabelsResult, int, error) {
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return nil, http.StatusBadRequest, err
	}

	result, err := u.repository.ListPlayerLabels(ctx, input)
	if err != nil {
		if errors.Is(err, ErrPlayerLabelsNotFound) {
			return nil, http.StatusNotFound, err
		}
		return nil, http.StatusInternalServerError, ErrPlayerLabelsLoadFailed
	}

	return result, http.StatusOK, nil
}

func (u *PlayerLabels) Vote(
	ctx context.Context,
	input PlayerLabelVoteInput,
) (*PlayerLabelVoteResult, int, error) {
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return nil, http.StatusBadRequest, err
	}

	result, err := u.repository.SavePlayerLabelVote(ctx, input)
	if err != nil {
		if errors.Is(err, ErrPlayerLabelsNotFound) {
			return nil, http.StatusNotFound, err
		}
		return nil, http.StatusInternalServerError, ErrPlayerLabelVoteFailed
	}

	return result, http.StatusOK, nil
}

func (i PlayerLabelsInput) Normalize() PlayerLabelsInput {
	return PlayerLabelsInput{PUUID: strings.TrimSpace(i.PUUID)}
}

func (i PlayerLabelsInput) Validate() error {
	if i.PUUID == "" {
		return ErrInvalidPlayerLabelsInput
	}
	return nil
}

func (i PlayerLabelVoteInput) Normalize() PlayerLabelVoteInput {
	return PlayerLabelVoteInput{
		PUUID:    strings.TrimSpace(i.PUUID),
		LabelID:  i.LabelID,
		VoterKey: strings.TrimSpace(i.VoterKey),
	}
}

func (i PlayerLabelVoteInput) Validate() error {
	if i.PUUID == "" || i.LabelID <= 0 || i.VoterKey == "" || len(i.VoterKey) > maxPlayerLabelVoterKeyLength {
		return ErrInvalidPlayerLabelVoteInput
	}
	return nil
}
