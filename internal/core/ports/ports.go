package ports

import "github.com/AntonyIS/notlify-content-svc/internal/core/domain"

type ContentService interface {
	CreateContent(Content *domain.Content) (*domain.Content, error)
	ReadContent(id string) (*domain.ContentResponse, error)
	ReadContents() ([]domain.Content, error)
	UpdateContent(Content *domain.Content) (*domain.Content, error)
	DeleteContent(id string) (string, error)
	ReadCreatorContents(creator_id string) ([]domain.Content, error)
	DeleteAllContent() (string, error)
}

type ContentRepository interface {
	CreateContent(Content *domain.Content) (*domain.Content, error)
	ReadContent(id string) (*domain.ContentResponse, error)
	ReadContents() ([]domain.Content, error)
	UpdateContent(Content *domain.Content) (*domain.Content, error)
	DeleteContent(id string) (string, error)
	DeleteAllContent() (string, error)
}

type ArticleService interface {
	CreateArticle(article *domain.Article) (*domain.Article, error)
	GetArticleByID(article_id string) (*domain.Article, error)
	GetArticles() (*[]domain.Article, error)
	GetArticlesByAuthor(author_id string) (*[]domain.Article, error)
	GetArticlesByTag(tag string) (*[]domain.Article, error)
	UpdateArticle(article *domain.Article) (*domain.Article, error)
	DeleteArticle(article_id string) error
	DeleteArticleAll() error
}

type ArticleRepository interface {
	CreateArticle(article *domain.Article) (*domain.Article, error)
	GetArticleByID(article_id string) (*domain.Article, error)
	GetArticles() (*[]domain.Article, error)
	GetArticlesByAuthor(author_id string) (*[]domain.Article, error)
	GetArticlesByTag(tag string) (*[]domain.Article, error)
	UpdateArticle(article *domain.Article) (*domain.Article, error)
	DeleteArticle(article_id string) error
	DeleteArticleAll() error
}
