package domain

import (
	"time"
)

type User struct {
	Id           string    `json:"id"`
	Firstname    string    `json:"firstname"`
	Lastname     string    `json:"lastname"`
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	Handle       string    `json:"handle"`
	About        string    `json:"about"`
	Contents     []Content `json:"contents"`
	ProfileImage string    `json:"profile_image"`
	Following    int       `json:"following"`
	Followers    int       `json:"followers"`
}

type ContentUser struct {
	Id           string `json:"id"`
	Firstname    string `json:"firstname"`
	Lastname     string `json:"lastname"`
	Handle       string `json:"handle"`
	About        string `json:"about"`
	ProfileImage string `json:"profile_image"`
	Following    int    `json:"following"`
	Followers    int    `json:"followers"`
}

type Content struct {
	ContentId       string      `json:"content_id"`
	Title           string      `json:"title"`
	Body            string      `json:"body"`
	User            ContentUser `json:"user"`
	PublicationDate time.Time   `json:"publication_date"`
}
