package usecase

import (
	"context"
	"errors"
)

var ErrListLabelsFailed = errors.New("failed to list labels")

type Label struct {
	ID   int16  `json:"id"`
	Name string `json:"name"`
}

type LabelsInput struct {
	ForceRefresh bool
}

type LabelsCache interface {
	Get() ([]Label, bool)
	Set([]Label)
	Delete()
}

type LabelsRepository interface {
	ListLabels(ctx context.Context) ([]Label, error)
}

type Labels struct {
	cache      LabelsCache
	repository LabelsRepository
}

func NewLabels(
	cache LabelsCache,
	repository LabelsRepository,
) *Labels {
	return &Labels{
		cache:      cache,
		repository: repository,
	}
}

func (u *Labels) Execute(ctx context.Context, input LabelsInput) ([]Label, error) {
	if !input.ForceRefresh {
		if cached, found := u.cache.Get(); found {
			return cached, nil
		}
	}

	labels, err := u.repository.ListLabels(ctx)
	if err != nil {
		return nil, ErrListLabelsFailed
	}

	if input.ForceRefresh {
		u.cache.Delete()
	}
	u.cache.Set(labels)

	return labels, nil
}
