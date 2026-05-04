package usecase

import (
	"context"
	"errors"
	"net/http"
	"testing"
)

type fakePlayerLabelsRepository struct {
	listResult *PlayerLabelsResult
	voteResult *PlayerLabelVoteResult
	err        error
	listInput  PlayerLabelsInput
	voteInput  PlayerLabelVoteInput
}

func (f *fakePlayerLabelsRepository) ListPlayerLabels(
	ctx context.Context,
	input PlayerLabelsInput,
) (*PlayerLabelsResult, error) {
	f.listInput = input
	return f.listResult, f.err
}

func (f *fakePlayerLabelsRepository) SavePlayerLabelVote(
	ctx context.Context,
	input PlayerLabelVoteInput,
) (*PlayerLabelVoteResult, error) {
	f.voteInput = input
	return f.voteResult, f.err
}

func TestPlayerLabels_ListNormalizesInput(t *testing.T) {
	t.Parallel()

	repository := &fakePlayerLabelsRepository{listResult: &PlayerLabelsResult{}}
	usecase := NewPlayerLabels(repository)

	_, statusCode, err := usecase.List(context.Background(), PlayerLabelsInput{
		PUUID: " test-puuid ",
	})

	if err != nil {
		t.Fatalf("err = %v, want nil", err)
	}
	if statusCode != http.StatusOK {
		t.Fatalf("statusCode = %d, want %d", statusCode, http.StatusOK)
	}
	if repository.listInput.PUUID != "test-puuid" {
		t.Fatalf("PUUID = %q, want normalized value", repository.listInput.PUUID)
	}
}

func TestPlayerLabels_VoteMapsNotFound(t *testing.T) {
	t.Parallel()

	repository := &fakePlayerLabelsRepository{err: ErrPlayerLabelsNotFound}
	usecase := NewPlayerLabels(repository)

	_, statusCode, err := usecase.Vote(context.Background(), PlayerLabelVoteInput{
		PUUID:    "test-puuid",
		LabelID:  1,
		VoterKey: "test-voter",
	})

	if !errors.Is(err, ErrPlayerLabelsNotFound) {
		t.Fatalf("err = %v, want ErrPlayerLabelsNotFound", err)
	}
	if statusCode != http.StatusNotFound {
		t.Fatalf("statusCode = %d, want %d", statusCode, http.StatusNotFound)
	}
}

func TestPlayerLabels_VoteRejectsInvalidInput(t *testing.T) {
	t.Parallel()

	usecase := NewPlayerLabels(&fakePlayerLabelsRepository{})

	_, statusCode, err := usecase.Vote(context.Background(), PlayerLabelVoteInput{
		PUUID:    "test-puuid",
		LabelID:  0,
		VoterKey: "test-voter",
	})

	if !errors.Is(err, ErrInvalidPlayerLabelVoteInput) {
		t.Fatalf("err = %v, want ErrInvalidPlayerLabelVoteInput", err)
	}
	if statusCode != http.StatusBadRequest {
		t.Fatalf("statusCode = %d, want %d", statusCode, http.StatusBadRequest)
	}
}
