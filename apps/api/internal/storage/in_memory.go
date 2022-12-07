package storage

import (
	"context"
	"promptu/api/internal/model"
	"sync"
	"time"
)

type InMemoryStore struct {
	mu                        sync.Mutex
	feed                      []model.Post
	notificationLocalTime     time.Time
	nextNotificationLocalTime time.Time
}

func NewInMemoryStore() Store {
	return &InMemoryStore{feed: []model.Post{}}
}

func (s *InMemoryStore) GetFeed(ctx context.Context, date time.Time) ([]model.Post, error) {
	filteredFeed := []model.Post{}
	dayAfter := date.Add(time.Hour * 24)

	var within24hrsOfDate bool

	for _, post := range s.feed {
		within24hrsOfDate = (post.UtcCreatedAt.Before(dayAfter) && post.UtcCreatedAt.After(date)) || post.UtcCreatedAt.Equal(date)

		if within24hrsOfDate {
			filteredFeed = append(filteredFeed, post)
		}
	}

	return filteredFeed, nil
}

func (s *InMemoryStore) PostMessage(ctx context.Context, post model.Post) error {
	s.feed = append(s.feed, post)
	return nil
}
