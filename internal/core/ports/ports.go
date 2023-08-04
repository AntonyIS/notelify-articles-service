package ports

import "github.com/AntonyIS/notlify-content-svc/internal/core/domain"

type ContentService interface {
	CreateContent(Content *domain.Content) (*domain.Content, error)
	ReadContent(id string) (*domain.Content, error)
	ReadContents() ([]domain.Content, error)
	UpdateContent(Content *domain.Content) (*domain.Content, error)
	DeleteContent(id string) (string, error)
	ReadCreatorContents(creator_id string) ([]domain.Content, error)
}

type ContentRepository interface {
	CreateContent(Content *domain.Content) (*domain.Content, error)
	ReadContent(id string) (*domain.Content, error)
	ReadContents() ([]domain.Content, error)
	UpdateContent(Content *domain.Content) (*domain.Content, error)
	DeleteContent(id string) (string, error)
}