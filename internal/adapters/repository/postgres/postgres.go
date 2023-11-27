package postgres

import (
	"database/sql"
	"fmt"

	appConfig "github.com/AntonyIS/notelify-articles-service/config"
	"github.com/AntonyIS/notelify-articles-service/internal/core/domain"
	"github.com/AntonyIS/notelify-articles-service/internal/core/ports"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type postgresDBClient struct {
	db            *sql.DB
	tablename     string
	loggerService ports.Logger
}

func NewPostgresClient(appConfig appConfig.Config, logger ports.Logger) (*postgresDBClient, error) {
	dbname := appConfig.DatabaseName
	tablename := appConfig.ContentTable
	user := appConfig.DatabaseUser
	password := appConfig.DatabasePassword
	port := appConfig.DatabasePort
	host := appConfig.DatabaseHost
	region := appConfig.AWS_DEFAULT_REGION
	rdsInstanceIdentifier := appConfig.RDSInstanceIdentifier
	var dsn string

	if appConfig.Env == "dev" || appConfig.Env == "test" {
		dsn = fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", host, port, user, dbname, password)

	}

	if appConfig.Env == "prod" || appConfig.Env == "test_prod" {
		awsSession := session.Must(session.NewSession(&aws.Config{
			Region: aws.String(region),
		}))
		rdsClient := rds.New(awsSession)

		describeInput := &rds.DescribeDBInstancesInput{
			DBInstanceIdentifier: &rdsInstanceIdentifier,
		}

		describeOutput, err := rdsClient.DescribeDBInstances(describeInput)
		if err != nil {
			logger.Error(fmt.Sprintf("Failed to describe DB instance: %s", err.Error()))
		}

		if len(describeOutput.DBInstances) == 0 {
			logger.Error("DB instance not found")
		}

		dsn = fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=require", host, port, dbname, user, password)
	}

	db, err := sql.Open("postgres", dsn)

	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	queryString := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			article_id VARCHAR(255) PRIMARY KEY UNIQUE,
			title VARCHAR(255) NOT NULL,
			subtitle VARCHAR(255),
			introduction TEXT,
			body TEXT,
			tags TEXT[],
			publish_date TIMESTAMP,
			author_id VARCHAR(255) NOT NULL
	)
	`, tablename)

	_, err = db.Exec(queryString)
	if err != nil {
		return nil, err
	}

	if err != nil {
		logger.Error(err.Error())
		return nil, err

	}

	return &postgresDBClient{db: db, tablename: tablename, loggerService: logger}, nil
}

func (psql *postgresDBClient) CreateArticle(article *domain.Article) (*domain.Article, error) {
	query := fmt.Sprintf(`
		INSERT INTO %s (
			article_id,
			title,
			subtitle,
			introduction,
			body,
			tags,
			publish_date,
			author_id
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		psql.tablename)
	_, err := psql.db.Exec(
		query,
		article.ArticleID,
		article.Title,
		article.Subtitle,
		article.Introduction,
		article.Body,
		pq.Array(article.Tags),
		article.PublishDate,
		article.AuthorID,
	)

	if err != nil {
		psql.loggerService.Error(err.Error())
		return nil, err
	}

	return article, nil
}

func (psql *postgresDBClient) GetArticleByID(article_id string) (*domain.Article, error) {
	query := fmt.Sprintf(`SELECT article_id,title,subtitle,introduction,body,tags,publish_date,author_id FROM %s WHERE article_id = $1`, psql.tablename)
	article := &domain.Article{}
	row := psql.db.QueryRow(query, article_id)
	err := row.Scan(&article.ArticleID, &article.Title, &article.Subtitle, &article.Introduction, &article.Body, pq.Array(&article.Tags), &article.PublishDate, &article.AuthorID)
	if err != nil {
		return nil, err
	}
	return article, nil
}

func (psql *postgresDBClient) GetArticlesByAuthor(author_id string) (*[]domain.Article, error) {
	query := fmt.Sprintf(`
		SELECT 
		article_id, 
		title, 
		subtitle, 
		introduction, 
		body, 
		tags, 
		publish_date,
		author_id
		FROM %s WHERE author_id = $1`, psql.tablename)
	rows, err := psql.db.Query(query, author_id)

	if err != nil {
		psql.loggerService.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	articles := []domain.Article{}

	for rows.Next() {
		var article domain.Article
		err := rows.Scan(
			&article.ArticleID,
			&article.Title,
			&article.Subtitle,
			&article.Introduction,
			&article.Body,
			pq.Array(&article.Tags),
			&article.PublishDate,
			&article.ArticleID,
		)
		if err != nil {
			psql.loggerService.Error(err.Error())
			return nil, err
		}
		articles = append(articles, article)
	}

	return &articles, nil
}

func (psql *postgresDBClient) GetArticlesByTag(tag string) (*[]domain.Article, error) {
	query := fmt.Sprintf(`
	SELECT 
	article_id, 
	title, 
	subtitle, 
	introduction, 
	body, 
	tags, 
	publish_date,
	author_id
	FROM %s WHERE $1 = ANY(tags)`, psql.tablename)
	rows, err := psql.db.Query(query, tag)

	if err != nil {
		psql.loggerService.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	var articles []domain.Article

	for rows.Next() {
		var article domain.Article
		err := rows.Scan(
			&article.ArticleID,
			&article.Title,
			&article.Subtitle,
			&article.Introduction,
			&article.Body,
			pq.Array(&article.Tags),
			&article.PublishDate,
			&article.AuthorID,
		)
		if err != nil {
			psql.loggerService.Error(err.Error())
			return nil, err
		}
		articles = append(articles, article)
	}

	return &articles, nil
}

func (psql *postgresDBClient) GetArticles() (*[]domain.Article, error) {
	query := fmt.Sprintf(`
	SELECT 
	article_id, 
	title, 
	subtitle, 
	introduction, 
	body, 
	tags, 
	publish_date,
	author_id
	FROM %s `, psql.tablename)

	rows, err := psql.db.Query(query)
	if err != nil {
		psql.loggerService.Error(err.Error())
		return nil, err
	}
	defer rows.Close()
	articles := []domain.Article{}
	for rows.Next() {
		var article domain.Article
		err := rows.Scan(
			&article.ArticleID,
			&article.Title,
			&article.Subtitle,
			&article.Introduction,
			&article.Body,
			pq.Array(&article.Tags),
			&article.PublishDate,
			&article.AuthorID,
		)
		if err != nil {
			psql.loggerService.Error(err.Error())
			return nil, err
		}
		articles = append(articles, article)
	}

	return &articles, nil
}

func (psql *postgresDBClient) UpdateArticle(article_id string, article *domain.Article) (*domain.Article, error) {
	DBArticle, err := psql.GetArticleByID(article_id)

	if err != nil {
		psql.loggerService.Error(err.Error())
		return nil, err
	}

	DBArticle.Title = article.Title
	DBArticle.Subtitle = article.Subtitle
	DBArticle.Introduction = article.Introduction
	DBArticle.Body = article.Body
	DBArticle.Tags = article.Tags
	DBArticle.PublishDate = article.PublishDate
	DBArticle.AuthorID = article.AuthorID

	query := fmt.Sprintf(`UPDATE %s SET title=$1,subtitle=$2,introduction=$3,body=$4,tags=$5,publish_date=$6,author_id=$7	WHERE article_id=$8`, psql.tablename)

	_, err = psql.db.Exec(query, DBArticle.Title, DBArticle.Subtitle, DBArticle.Introduction, DBArticle.Body, pq.Array(DBArticle.Tags), DBArticle.PublishDate, DBArticle.AuthorID, DBArticle.ArticleID)

	if err != nil {
		psql.loggerService.Error(err.Error())
		return nil, err
	}

	return psql.GetArticleByID(article_id)
}

func (psql *postgresDBClient) DeleteArticle(article_id string) error {
	_, err := psql.GetArticleByID(article_id)

	if err != nil {
		psql.loggerService.Error(err.Error())
		return err
	}
	query := fmt.Sprintf(`DELETE FROM %s WHERE article_id = $1`, psql.tablename)

	_, err = psql.db.Exec(query, article_id)
	if err != nil {
		psql.loggerService.Error(err.Error())
		return err
	}
	return nil
}

func (psql *postgresDBClient) DeleteArticleAll() error {
	query := fmt.Sprintf(`DELETE FROM %s `, psql.tablename)
	_, err := psql.db.Exec(query)
	if err != nil {
		psql.loggerService.Error(err.Error())
		return err
	}
	return nil
}
