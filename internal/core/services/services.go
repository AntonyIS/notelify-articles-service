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
	article.Author.ID = "2361a252-215b-4174-805c-6b67fcb428dc"
	article.Author.Name = "Antony Injila"
	article.Author.Bio = "Hello guys, if you want to become a professional Java developer or want to take your Java skill to next level but are not sure which technology, tools,"
	article.Author.ProfilePicture = " "
	article.Author.SocialLinks = []string{"https://medium.com/@antonyshikubu"}
	article.Author.Following = 1000
	article.Author.Followers = 100
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
