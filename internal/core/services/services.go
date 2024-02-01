package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
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
func NewLoggingManagementService(loggerURL string) *loggingManagementService {
	svc := loggingManagementService{
		loggerURL: loggerURL,
	}
	return &svc
}

func (svc *articleManagementService) CreateArticle(article *domain.Article) (*domain.Article, error) {
	article.ArticleID = uuid.New().String()
	article.PublishDate = time.Now()
	article.UpdatedDate = time.Now()

	article, err := svc.repo.CreateArticle(article)
	if err != nil {
		logEntry := domain.LogMessage{
			LogLevel: "ERROR",
			Service:  "articles",
			Message:  err.Error(),
		}
		svc.logger.LogError(logEntry)
		return nil, err
	}
	logEntry := domain.LogMessage{
		LogLevel: "INFO",
		Service:  "articles",
		Message:  "Article created successufly",
	}
	svc.logger.LogInfo(logEntry)
	return article, nil
}

func (svc *articleManagementService) GetArticleByID(article_id string) (*domain.Article, error) {
	article, err := svc.repo.GetArticleByID(article_id)
	if err != nil {
		logEntry := domain.LogMessage{
			LogLevel: "ERROR",
			Service:  "articles",
			Message:  err.Error(),
		}
		svc.logger.LogError(logEntry)
		return nil, err
	}
	logEntry := domain.LogMessage{
		LogLevel: "INFO",
		Service:  "articles",
		Message:  "Article with ID [%s] found successufly",
	}
	svc.logger.LogInfo(logEntry)
	return article, nil
}

func (svc *articleManagementService) GetArticlesByAuthor(author_id string) (*[]domain.Article, error) {
	articles, err := svc.repo.GetArticlesByAuthor(author_id)
	if err != nil {
		logEntry := domain.LogMessage{
			LogLevel: "ERROR",
			Service:  "articles",
			Message:  err.Error(),
		}
		svc.logger.LogError(logEntry)
		return nil, err
	}
	logEntry := domain.LogMessage{
		LogLevel: "INFO",
		Service:  "articles",
		Message:  "Articles by author found successufly",
	}
	svc.logger.LogInfo(logEntry)
	return articles, nil
}

func (svc *articleManagementService) GetArticlesByTag(tag string) (*[]domain.Article, error) {
	articles, err := svc.GetArticles()
	if err != nil {
		logEntry := domain.LogMessage{
			LogLevel: "ERROR",
			Service:  "articles",
			Message:  err.Error(),
		}
		svc.logger.LogError(logEntry)
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
	logEntry := domain.LogMessage{
		LogLevel: "INFO",
		Service:  "articles",
		Message:  "Articles by tag found successufly",
	}
	svc.logger.LogInfo(logEntry)
	return &articleArray, nil
}

func (svc *articleManagementService) GetArticles() (*[]domain.Article, error) {
	artciles, err := svc.repo.GetArticles()
	if err != nil {
		logEntry := domain.LogMessage{
			LogLevel: "ERROR",
			Service:  "articles",
			Message:  err.Error(),
		}
		svc.logger.LogError(logEntry)
		return nil, err
	}
	logEntry := domain.LogMessage{
		LogLevel: "INFO",
		Service:  "articles",
		Message:  "Articles found successufly",
	}
	svc.logger.LogInfo(logEntry)
	return artciles, nil
}

func (svc *articleManagementService) UpdateArticle(article_id string, article *domain.Article) (*domain.Article, error) {
	article, err := svc.repo.UpdateArticle(article_id, article)
	if err != nil {
		logEntry := domain.LogMessage{
			LogLevel: "ERROR",
			Service:  "articles",
			Message:  err.Error(),
		}
		svc.logger.LogError(logEntry)
		return nil, err
	}
	logEntry := domain.LogMessage{
		LogLevel: "INFO",
		Service:  "articles",
		Message:  "Article with ID [%s] updated successufly",
	}
	svc.logger.LogInfo(logEntry)
	return article, nil
}

func (svc *articleManagementService) DeleteArticle(article_id string) error {
	err := svc.repo.DeleteArticle(article_id)
	if err != nil {
		logEntry := domain.LogMessage{
			LogLevel: "ERROR",
			Service:  "articles",
			Message:  err.Error(),
		}
		svc.logger.LogError(logEntry)
		return err
	}
	logEntry := domain.LogMessage{
		LogLevel: "INFO",
		Service:  "articles",
		Message:  "Article with ID [%s] deleted successufly",
	}
	svc.logger.LogInfo(logEntry)
	return nil
}

func (svc *articleManagementService) DeleteArticleAll() error {
	err := svc.repo.DeleteArticleAll()
	if err != nil {
		logEntry := domain.LogMessage{
			LogLevel: "ERROR",
			Service:  "articles",
			Message:  err.Error(),
		}
		svc.logger.LogError(logEntry)
		return err
	}
	logEntry := domain.LogMessage{
		LogLevel: "INFO",
		Service:  "articles",
		Message:  "Articles deleted successufly",
	}
	svc.logger.LogInfo(logEntry)
	return nil
}

type loggingManagementService struct {
	loggerURL string
}

func (svc *loggingManagementService) SendLog(logEntry domain.LogMessage) {
	// Marshal the struct into JSON
	payloadBytes, err := json.Marshal(logEntry)
	if err != nil {
		fmt.Println("Error encoding JSON payload:", err)
		return
	}

	resp, err := http.Post(svc.loggerURL, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		fmt.Println("Error making POST request:", err)
		return
	}
	defer resp.Body.Close()
}

func (svc *loggingManagementService) LogDebug(logEntry domain.LogMessage) {
	message := fmt.Sprintf("[%s] [DEBUG] %s %s", logEntry.Service, getCurrentDateTime(), logEntry.Message)
	logEntry.Message = message
	svc.SendLog(logEntry)
}

func (svc *loggingManagementService) LogInfo(logEntry domain.LogMessage) {
	message := fmt.Sprintf("[%s] [INFO] %s %s", logEntry.Service, getCurrentDateTime(), logEntry.Message)
	logEntry.Message = message
	svc.SendLog(logEntry)
}

func (svc *loggingManagementService) LogWarning(logEntry domain.LogMessage) {
	message := fmt.Sprintf("[%s] [WARNING] %s %s", logEntry.Service, getCurrentDateTime(), logEntry.Message)
	logEntry.Message = message
	svc.SendLog(logEntry)
}

func (svc *loggingManagementService) LogError(logEntry domain.LogMessage) {
	message := fmt.Sprintf("[%s] [ERROR] %s %s", logEntry.Service, getCurrentDateTime(), logEntry.Message)
	logEntry.Message = message
	svc.SendLog(logEntry)
}

func getCurrentDateTime() string {
	currentTime := time.Now()
	return currentTime.Format("2006/01/02 15:04:05")
}
