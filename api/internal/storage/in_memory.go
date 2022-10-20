package storage

import (	
	"context"
	"time"
	"promptu/api/internal/model"
)

type InMemoryStore struct {
	feed model.Feed
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{feed: model.Feed{}}
}

func (s *InMemoryStore) GetFeed(ctx context.Context) (model.Feed, error) {
	return s.feed, nil
}

func (s *InMemoryStore) PostMessage(ctx context.Context, post model.Post, createdAt time.Time) error {
	s.feed.Posts = append(s.feed.Posts, model.Post{User: post.User, Message: post.Message})
	return nil 
}
