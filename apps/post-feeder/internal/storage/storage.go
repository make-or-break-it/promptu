package storage

import (
	"context"
	"promptu/apps/feeder-api/internal/model"
	"time"
)

type Store interface {
	GetFeed(ctx context.Context, date time.Time) ([]model.Post, error)
}
