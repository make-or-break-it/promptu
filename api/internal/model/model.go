package model

import "time"

type Feed struct {
	Posts []Post `json:"posts,omitempty"`
}

type Post struct {
	User      string    `json:"user"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"createdAt"`
}
