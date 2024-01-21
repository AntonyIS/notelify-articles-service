package services

/*
	Filename : services_test.go
	Description: This file contains tests for the application services
	Author: Antony Injila
	Date : September 13, 2023
*/

import (
	"reflect"
	"testing"

	"github.com/AntonyIS/notelify-articles-service/config"
	"github.com/AntonyIS/notelify-articles-service/internal/adapters/app"
	"github.com/AntonyIS/notelify-articles-service/internal/adapters/repository/postgres"
	"github.com/AntonyIS/notelify-articles-service/internal/core/domain"
)

func TestApplicationService(t *testing.T) {
	// env := "prod"
	// Read application environment and load configurations
	conf, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	// Initialize the logging service
	loggerSvc := NewLoggingManagementService(conf.LOGGER_URL)
	// // Postgres Clien
	databaseRepo, err := postgres.NewPostgresClient(*conf)
	if err != nil {
		logEntry := domain.LogMessage{
			LogLevel: "ERROR",
			Service:  "articles",
			Message:  err.Error(),
		}
		loggerSvc.LogError(logEntry)
		panic(err)
	}
	articleService := NewArticleManagementService(databaseRepo, loggerSvc)
	app.InitGinRoutes(articleService, loggerSvc, *conf)
	author := domain.Author{
		AuthorID:         "b967127d-7535-420c-96a7-1d01b437a619",
		Firstname:        "Antony",
		Lastname:         "Injila",
		Handle:           "About Antony Injila",
		About:            "Are you managing a complex software project? In that case, you probably already know how difficult it is to predict and plan all the potential scenarios several months ahead. After all, many things can suddenly go wrong without any warning.",
		ProfileImage:     "ProfileImage",
		SocialMediaLinks: []string{"https://github.com/AntonyIS", "https://medium.com/@antonyshikubu"},
		Following:        100,
		Followers:        100,
	}

	err = articleService.DeleteArticleAll()
	if err != nil {
		t.Error(err)
	}
	t.Run("Test create new article", func(t *testing.T) {
		err := articleService.DeleteArticleAll()
		if err != nil {
			t.Error(err)
		}
		title := "Article - Create article"
		body := "Article body"
		tags := []string{"Golang"}
		article := &domain.Article{
			Title:    title,
			Body:     body,
			Tags:     tags,
			AuthorID: author.AuthorID,
		}

		article, err = articleService.CreateArticle(article)

		if err != nil {
			t.Error(err)
		}

		if article.ArticleID == "" {
			t.Error("Test article title != ", title)
		}

		if article.Title != title {
			t.Error("Test article title != ", title)
		}
		if article.Body != body {
			t.Error("Test article body != ", body)
		}
	})

	t.Run("Test read article by id", func(t *testing.T) {
		// Create article
		title := "Article - Read article"
		body := "Article body"
		tags := []string{"Golang"}
		article := &domain.Article{
			Title:    title,
			Body:     body,
			Tags:     tags,
			AuthorID: author.AuthorID,
		}

		article, err := articleService.CreateArticle(article)
		if err != nil {
			t.Error(err)
		}

		// Read article with articleID from the database
		res, err := articleService.repo.GetArticleByID(article.ArticleID)
		if err != nil {
			t.Error(err)
		}

		if res.ArticleID == "" {
			t.Error("Test article title != ", title)
		}

		if res.Title != title {
			t.Error("Test article title != ", title)
		}
		if res.Body != body {
			t.Error("Test article body != ", body)
		}
	})

	t.Run("Test Get articles by author", func(t *testing.T) {
		articles, err := articleService.GetArticlesByAuthor(author.AuthorID)
		if err != nil {
			t.Error(err)
		}
		results := reflect.TypeOf(articles) == reflect.TypeOf([]domain.Article{})

		if !results {
			if err != nil {
				t.Error("Expected array of articles")
			}
		}
	})

	t.Run("Test get articles by tags", func(t *testing.T) {
		articles, err := articleService.GetArticlesByTag("Golang")
		if err != nil {
			t.Error(err)
		}

		results := reflect.TypeOf(articles) == reflect.TypeOf([]domain.Article{})
		if !results {
			if err != nil {
				t.Error("Expected array of articles")
			}
		}
		// At this point we have 2 articles , test number of articles returned
		size := len(*articles)
		if size != 2 {
			t.Error("Expected 2 articles, got ", size)
		}

	})

	t.Run("Test get all articles", func(t *testing.T) {
		articles, err := articleService.GetArticles()
		if err != nil {
			t.Error(err)
		}
		results := reflect.TypeOf(articles) == reflect.TypeOf([]domain.Article{})
		if !results {
			if err != nil {
				t.Error("Expected array of articles")
			}
		}
	})

	t.Run("Test update article", func(t *testing.T) {

		title := "Article - Create article"
		body := "Article body"
		tags := []string{"Golang"}
		article := &domain.Article{
			Title:    title,
			Body:     body,
			Tags:     tags,
			AuthorID: author.AuthorID,
		}

		article, err := articleService.CreateArticle(article)

		if err != nil {
			t.Error(err)
		}
		newTitle := "Article - Update article title"
		newBody := "Article - Update article body"
		article.Title = newTitle
		article.Body = newBody

		res, err := articleService.UpdateArticle(article.ArticleID, article)

		if err != nil {
			t.Error(err)
		}

		if res.Title != newTitle {
			t.Errorf("Expected title '%s' got '%s", newTitle, res.Title)
		}
		if res.Body != newBody {
			t.Errorf("Expected title '%s' got '%s", newBody, res.Body)
		}

	})

	t.Run("Test Delete article", func(t *testing.T) {
		title := "Article - Create article"
		body := "Article body"
		tags := []string{"Golang"}
		article := &domain.Article{
			Title:    title,
			Body:     body,
			Tags:     tags,
			AuthorID: author.AuthorID,
		}

		article, err := articleService.CreateArticle(article)

		if err != nil {
			t.Error(err)
		}

		err = articleService.DeleteArticle(article.ArticleID)

		if err != nil {
			t.Error(err)
		}

		article, err = articleService.GetArticleByID(article.ArticleID)

		if article != nil {
			if err != nil {
				t.Error("Expected not to found an article, but article is found!")
			}
		}

	})

	t.Run("Test all articles", func(t *testing.T) {
		err := articleService.DeleteArticleAll()
		if err != nil {
			t.Error("Expected to delete all articles: ", err)
		}

		articles, err := articleService.GetArticles()
		if err != nil {
			t.Error("Expected to delete all articles: ", err)
		}

		if len(*articles) > 1 {
			t.Error("Expected 0 articles, found", len(*articles))
		}
	})

}
