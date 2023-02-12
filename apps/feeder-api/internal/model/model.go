package model

import "time"

type Post struct {
	User                string    `json:"user" bson:"user"`
	Message             string    `json:"message" bson:"message"`
	UtcCreatedAt        time.Time `json:"utcCreatedAt" bson:"utcCreatedAt"`
	UtcNotificationTime time.Time `json:"utcNotificationTime" bson:"utcNotificationTime"`
}
