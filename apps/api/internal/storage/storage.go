package storage

import (
	"context"
	"promptu/api/internal/model"
	"time"
)

type Store interface {
	GetFeed(ctx context.Context, date time.Time) ([]model.Post, error)
	PostMessage(ctx context.Context, post model.Post) error
}
