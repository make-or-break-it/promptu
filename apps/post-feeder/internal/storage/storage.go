package storage

import (
	"context"
	"promptu/apps/post-feeder/internal/model"
	"time"
)

type Store interface {
	GetFeed(ctx context.Context, date time.Time) ([]model.Post, error)
}
