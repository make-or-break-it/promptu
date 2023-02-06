package storage

import (
	"context"
	"promptu/apps/post-db-updater/internal/model"
)

type Store interface {
	PostMessage(ctx context.Context, post model.Post) error
}
