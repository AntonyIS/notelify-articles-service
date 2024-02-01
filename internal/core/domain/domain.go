package domain

import (
	"time"
)

type Article struct {
	ArticleID    string    `json:"article_id"`
	Title        string    `json:"title"`
	Subtitle     string    `json:"subtitle"`
	Introduction string    `json:"introduction"`
	Body         string    `json:"body"`
	Tags         []string  `json:"tags"`
	PublishDate  time.Time `json:"publish_date"`
	UpdatedDate  time.Time `json:"updated_date"`
	Author       Author    `json:"author"`
	AuthorID     string    `json:"author_id"`
}

type Author struct {
	AuthorID         string   `json:"author_id"`
	Firstname        string   `json:"firstname"`
	Lastname         string   `json:"lastname"`
	Handle           string   `json:"handle"`
	About            string   `json:"about"`
	ProfileImage     string   `json:"profile_image"`
	SocialMediaLinks []string `json:"social_media_links"`
	Following        int      `json:"following"`
	Followers        int      `json:"followers"`
}

type LogMessage struct {
	LogLevel string `json:"log_level"`
	Message  string `json:"message"`
	Service  string `json:"service"`
}
