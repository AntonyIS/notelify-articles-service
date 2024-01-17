package services

import (
	"sort"
	"strings"
	"time"

	"github.com/AntonyIS/notelify-articles-service/internal/core/domain"
	"github.com/AntonyIS/notelify-articles-service/internal/core/ports"
	"github.com/google/uuid"
)

type articleManagementService struct {
	repo   ports.ArticleRepository
	logger ports.LoggingService
}

func NewArticleManagementService(repo ports.ArticleRepository, logger ports.LoggingService) *articleManagementService {
	svc := articleManagementService{
		repo:   repo,
		logger: logger,
	}
	return &svc
}

func (svc *articleManagementService) CreateArticle(article *domain.Article) (*domain.Article, error) {
	article.ArticleID = uuid.New().String()
	article.PublishDate = time.Now()
	article, err := svc.repo.CreateArticle(article)
	if err != nil {
		logEntry := domain.LogMessage{
			LogLevel: "critical",
			Service:  "articles",
			Message:  err.Error(),
		}
		svc.logger.CreateLog(logEntry)
		return nil, err
	}
	return article, nil
}

func (svc *articleManagementService) GetArticleByID(article_id string) (*domain.Article, error) {
	article, err := svc.repo.GetArticleByID(article_id)
	if err != nil {
		logEntry := domain.LogMessage{
			LogLevel: "critical",
			Service:  "articles",
			Message:  err.Error(),
		}
		svc.logger.CreateLog(logEntry)
		return nil, err
	}
	return article, nil
}

func (svc *articleManagementService) GetArticlesByAuthor(author_id string) (*[]domain.Article, error) {
	articles, err := svc.repo.GetArticlesByAuthor(author_id)
	if err != nil {
		logEntry := domain.LogMessage{
			LogLevel: "critical",
			Service:  "articles",
			Message:  err.Error(),
		}
		svc.logger.CreateLog(logEntry)
		return nil, err
	}
	return articles, nil
}

func (svc *articleManagementService) GetArticlesByTag(tag string) (*[]domain.Article, error) {
	articles, err := svc.GetArticles()
	if err != nil {
		logEntry := domain.LogMessage{
			LogLevel: "critical",
			Service:  "articles",
			Message:  err.Error(),
		}
		svc.logger.CreateLog(logEntry)
		return nil, err
	}
	articleArray := []domain.Article{}
	for _, article := range *articles {
		arr := article.Tags
		sort.Strings(arr)
		// Perform a case-insensitive search using sort.SearchStrings
		index := sort.SearchStrings(arr, tag)
		if index < len(arr) && (arr[index] == tag || arr[index] == strings.ToLower(tag) || arr[index] == strings.ToUpper(tag)) {
			articleArray = append(articleArray, article)

		}
	}
	return &articleArray, nil
}

func (svc *articleManagementService) GetArticles() (*[]domain.Article, error) {
	artciles, err := svc.repo.GetArticles()
	if err != nil {
		logEntry := domain.LogMessage{
			LogLevel: "critical",
			Service:  "articles",
			Message:  err.Error(),
		}
		svc.logger.CreateLog(logEntry)
		return nil, err
	}

	return artciles, nil
}

func (svc *articleManagementService) UpdateArticle(article_id string, article *domain.Article) (*domain.Article, error) {
	article, err := svc.repo.UpdateArticle(article_id, article)
	if err != nil {
		logEntry := domain.LogMessage{
			LogLevel: "critical",
			Service:  "articles",
			Message:  err.Error(),
		}
		svc.logger.CreateLog(logEntry)
		return nil, err
	}
	return article, nil
}

func (svc *articleManagementService) DeleteArticle(article_id string) error {
	err := svc.repo.DeleteArticle(article_id)
	if err != nil {
		logEntry := domain.LogMessage{
			LogLevel: "critical",
			Service:  "articles",
			Message:  err.Error(),
		}
		svc.logger.CreateLog(logEntry)
		return err
	}
	return nil
}

func (svc *articleManagementService) DeleteArticleAll() error {
	err := svc.repo.DeleteArticleAll()
	if err != nil {
		logEntry := domain.LogMessage{
			LogLevel: "critical",
			Service:  "articles",
			Message:  err.Error(),
		}
		svc.logger.CreateLog(logEntry)
		return err
	}
	return nil
}
