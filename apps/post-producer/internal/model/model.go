package model

import (
	"encoding/json"
	"time"
)

type Post struct {
	User                string    `json:"user" bson:"user"`
	Message             string    `json:"message" bson:"message"`
	UtcCreatedAt        time.Time `json:"utcCreatedAt" bson:"utcCreatedAt"`
	UtcNotificationTime time.Time `json:"utcNotificationTime" bson:"utcNotificationTime"`

	encoded []byte
}

func (p *Post) Encode() ([]byte, error) {
	encoded, err := json.Marshal(p)
	if err != nil {
		return []byte{}, err
	}
	p.encoded = encoded
	return encoded, nil
}

func (p *Post) Length() int {
	return len(p.encoded)
}
