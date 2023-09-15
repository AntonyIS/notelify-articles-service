/*
	Filename : services_test.go
	Description: This file contains tests for the application services
	Author: Antony Injila
	Date : September 13, 2023
*/
package services

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/AntonyIS/notelify-articles-service/config"
	"github.com/AntonyIS/notelify-articles-service/internal/adapters/logger"
	"github.com/AntonyIS/notelify-articles-service/internal/adapters/repository/postgres"
	"github.com/AntonyIS/notelify-articles-service/internal/core/domain"
)

func TestApplicationService(t *testing.T) {
	env := "prod"
	conf, err := config.NewConfig(env)
	if err != nil {
		panic(err)
	}

	// Logger service
	logger := logger.NewLoggerService(conf.LoggerURL)
	// // Postgres Client
	postgresDBRepo, err := postgres.NewPostgresClient(*conf, logger)
	articleService := NewArticleManagementService(postgresDBRepo)

	author := domain.AuthorInfo{
		ID:             "b967127d-7535-420c-96a7-1d01b437a619",
		Name:           "Antony Injila",
		Bio:            "About Antony Injila",
		ProfilePicture: "",
		SocialLinks:    []string{"https://github.com/AntonyIS", "https://medium.com/@antonyshikubu"},
		Following:      100,
		Followers:      100,
	}

	err = articleService.DeleteArticleAll()
	t.Run("Test create new article", func(t *testing.T) {
		err := articleService.DeleteArticleAll()
		title := "Article - Create article"
		body := "Article body"
		tags := []string{"Golang"}
		var article = &domain.Article{
			Title: title,
			Body:  body,
			Tags:  tags,
		}
		article.Author = author

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
		var article = &domain.Article{
			Title: title,
			Body:  body,
			Tags:  tags,
		}
		article.Author = author

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
		articles, err := articleService.GetArticlesByAuthor(author.ID)
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
		if len(*articles) != 2 {
			t.Error("Expected 2 articles, got ", len(*articles))
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
		var article = &domain.Article{
			Title: title,
			Body:  body,
			Tags:  tags,
		}
		article.Author = author

		article, err := articleService.CreateArticle(article)

		if err != nil {
			t.Error(err)
		}
		newTitle := "Article - Update article title"
		newBody := "Article - Update article body"
		article.Title = newTitle
		article.Body = newBody

		res, err := articleService.UpdateArticle(article)

		if res.Title != newTitle {
			t.Error(fmt.Sprintf("Expected title '%s' got '%s", newTitle, res.Title))
		}
		if res.Body != newBody {
			t.Error(fmt.Sprintf("Expected title '%s' got '%s", newBody, res.Body))
		}

	})

	t.Run("Test Delete article", func(t *testing.T) {
		title := "Article - Create article"
		body := "Article body"
		tags := []string{"Golang"}
		var article = &domain.Article{
			Title: title,
			Body:  body,
			Tags:  tags,
		}
		article.Author = author

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
