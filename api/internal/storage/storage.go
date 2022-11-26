package storage

import (
	"context"
	"promptu/api/internal/model"
	"time"
)

type Store interface {
	GetFeed(ctx context.Context) ([]model.Post, error)
	PostMessage(ctx context.Context, post model.Post, createdAt time.Time) error
}
