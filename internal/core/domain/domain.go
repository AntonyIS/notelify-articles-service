package domain

import (
	"time"
)

type Article struct {
	ArticleID    string     `json:"article_id"`
	Title        string     `json:"title"`
	Subtitle     string     `json:"subtitle"`
	Introduction string     `json:"introduction"`
	Body         string     `json:"body"`
	Tags         []string   `json:"tags"`
	PublishDate  time.Time  `json:"publish_date"`
	Author       AuthorInfo `json:"author_info"`
}

type AuthorInfo struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	Bio            string   `json:"bio"`
	ProfilePicture string   `json:"profile_picture"`
	SocialLinks    []string `json:"social_links"`
	Following      int      `json:"following"`
	Followers      int      `json:"followers"`
}

type User struct {
	UserId       string    `json:"user_id"`
	Firstname    string    `json:"firstname"`
	Lastname     string    `json:"lastname"`
	Handle       string    `json:"handle"`
	About        string    `json:"about"`
	ProfileImage string    `json:"profile_image"`
	Following    int       `json:"following"`
	Followers    int       `json:"followers"`
	Contents     []Content `json:"contents"`
}

type Content struct {
	ContentId       string    `json:"content_id"`
	CreatorId       string    `json:"creator_id"`
	Title           string    `json:"title"`
	Body            string    `json:"body"`
	PublicationDate time.Time `json:"publication_date"`
}

type ContentResponse struct {
	Content
	User
}
