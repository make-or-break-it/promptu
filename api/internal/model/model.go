package model

type Feed struct {
	Posts []Post `json:"posts,omitempty"`
}

type Post struct {
	User      string `json:"user"`
	Message   string `json:"message"`
	CreatedAt string `json:"createdAt"`
}
