package ports

import "github.com/AntonyIS/notelify-articles-service/internal/core/domain"

type ArticleService interface {
	CreateArticle(article *domain.Article) (*domain.Article, error)
	GetArticleByID(article_id string) (*domain.Article, error)
	GetArticles() (*[]domain.Article, error)
	GetArticlesByAuthor(author_id string) (*[]domain.Article, error)
	GetArticlesByTag(tag string) (*[]domain.Article, error)
	UpdateArticle(article_id string, article *domain.Article) (*domain.Article, error)
	DeleteArticle(article_id string) error
	DeleteArticleAll() error
}

type ArticleRepository interface {
	CreateArticle(article *domain.Article) (*domain.Article, error)
	GetArticleByID(article_id string) (*domain.Article, error)
	GetArticles() (*[]domain.Article, error)
	GetArticlesByAuthor(author_id string) (*[]domain.Article, error)
	GetArticlesByTag(tag string) (*[]domain.Article, error)
	UpdateArticle(article_id string, article *domain.Article) (*domain.Article, error)
	DeleteArticle(article_id string) error
	DeleteArticleAll() error
}

type LoggingService interface {
	SendLog(LogEntry domain.LogMessage)
	LogDebug(LogEntry domain.LogMessage)
	LogInfo(LogEntry domain.LogMessage)
	LogWarning(LogEntry domain.LogMessage)
	LogError(LogEntry domain.LogMessage)
}
