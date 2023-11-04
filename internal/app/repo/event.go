package repo

import (
    "github.com/ozonmp/omp-demo-api/internal/model"
)

type EventRepo interface {
    Lock(n uint64) ([]model.KeywordEvent, error)
    Unlock(eventIDs []uint64) error

    Add(event []model.KeywordEvent) error
    Remove(eventIDs []uint64) error
}
