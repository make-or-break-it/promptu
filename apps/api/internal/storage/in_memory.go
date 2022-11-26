package storage

import (
	"context"
	"promptu/api/internal/model"
	"time"
)

type InMemoryStore struct {
	feed []model.Post
}

func NewInMemoryStore() Store {
	return &InMemoryStore{feed: []model.Post{}}
}

func (s *InMemoryStore) GetFeed(ctx context.Context) ([]model.Post, error) {
	return s.feed, nil
}

func (s *InMemoryStore) PostMessage(ctx context.Context, post model.Post, createdAt time.Time) error {
	s.feed = append(s.feed, model.Post{User: post.User, Message: post.Message, CreatedAt: createdAt})
	return nil
}
