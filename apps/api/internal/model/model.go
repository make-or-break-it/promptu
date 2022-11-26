package model

import "time"

type Post struct {
	User      string    `json:"user"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"createdAt"`
}
