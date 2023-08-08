package domain

import (
	"time"
)

type User struct {
	UserId       string    `json:"user_id"`
	Firstname    string    `json:"firstname"`
	Lastname     string    `json:"lastname"`
	Handle       string    `json:"handle"`
	About        string    `json:"about"`
	ProfileImage string    `json:"profile_image"`
	Following    int       `json:"following"`
	Followers    int       `json:"followers"`
}

type Content struct {
	ContentId       string    `json:"content_id"`
	CreatorId       string    `json:"creator_id"`
	Title           string    `json:"title"`
	Body            string    `json:"body"`
	PublicationDate time.Time `json:"publication_date"`
}
