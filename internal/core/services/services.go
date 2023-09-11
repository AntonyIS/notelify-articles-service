package services

import (
	"time"

	"github.com/AntonyIS/notelify-articles-service/internal/core/domain"
	"github.com/AntonyIS/notelify-articles-service/internal/core/ports"
	"github.com/google/uuid"
)

type articleManagementService struct {
	repo ports.ArticleRepository
}

func NewArticleManagementService(repo ports.ArticleRepository) *articleManagementService {
	svc := articleManagementService{
		repo: repo,
	}
	return &svc
}

func (svc *articleManagementService) CreateArticle(article *domain.Article) (*domain.Article, error) {
	article.ArticleID = uuid.New().String()
	article.PublishDate = time.Now()
	return svc.repo.CreateArticle(article)
}

func (svc *articleManagementService) GetArticleByID(article_id string) (*domain.Article, error) {
	return svc.repo.GetArticleByID(article_id)
}

func (svc *articleManagementService) GetArticlesByAuthor(author_id string) (*[]domain.Article, error) {
	return svc.repo.GetArticlesByAuthor(author_id)
}

func (svc *articleManagementService) GetArticlesByTag(tag string) (*[]domain.Article, error) {
	return svc.repo.GetArticlesByTag(tag)
}

func (svc *articleManagementService) GetArticles() (*[]domain.Article, error) {
	return svc.repo.GetArticles()
}

func (svc *articleManagementService) UpdateArticle(article *domain.Article) (*domain.Article, error) {
	return svc.repo.UpdateArticle(article)
}

func (svc *articleManagementService) DeleteArticle(article_id string) error {
	return svc.repo.DeleteArticle(article_id)
}

func (svc *articleManagementService) DeleteArticleAll() error {
	return svc.repo.DeleteArticleAll()
}
