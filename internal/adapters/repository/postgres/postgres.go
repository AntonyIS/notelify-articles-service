package postgres

import (
	"database/sql"
	"fmt"

	appConfig "github.com/AntonyIS/notelify-articles-service/config"
	"github.com/AntonyIS/notelify-articles-service/internal/adapters/logger"
	"github.com/AntonyIS/notelify-articles-service/internal/core/domain"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type postgresDBClient struct {
	db            *sql.DB
	tablename     string
	loggerService logger.LoggerType
}

func NewPostgresClient(appConfig appConfig.Config, logger logger.LoggerType) (*postgresDBClient, error) {

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

	} else {
		awsSession := session.Must(session.NewSession(&aws.Config{
			Region: aws.String(region),
		}))
		rdsClient := rds.New(awsSession)

		describeInput := &rds.DescribeDBInstancesInput{
			DBInstanceIdentifier: &rdsInstanceIdentifier,
		}

		describeOutput, err := rdsClient.DescribeDBInstances(describeInput)
		if err != nil {
			logger.PostLogMessage(fmt.Sprintf("Failed to describe DB instance: %s", err.Error()))
		}

		if len(describeOutput.DBInstances) == 0 {
			logger.PostLogMessage("DB instance not found")
		}

		dsn = fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=require", host, port, dbname, user, password)
	}

	db, err := sql.Open("postgres", dsn)

	if err != nil {
		logger.PostLogMessage(err.Error())
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		logger.PostLogMessage(err.Error())
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
			author_id VARCHAR(255) NOT NULL,
			author_name VARCHAR(255) NOT NULL,
			author_bio TEXT,
			author_profile_pic VARCHAR(255),
			author_social_links TEXT[],
			author_followers INT,
			author_following INT
	)
	`, tablename)

	_, err = db.Exec(queryString)
	if err != nil {
		return nil, err
	}

	if err != nil {
		logger.PostLogMessage(err.Error())
		return nil, err

	}

	return &postgresDBClient{db: db, tablename: tablename, loggerService: logger}, nil
}

func (psql *postgresDBClient) CreateArticle(article *domain.Article) (*domain.Article, error) {
	query := fmt.Sprintf(`
		INSERT INTO %s (
			article_id,title,subtitle,introduction,body,tags,publish_date,author_id,author_name,author_bio,author_profile_pic,author_social_links,author_followers,author_following)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)`,
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
		article.Author.ID,
		article.Author.Name,
		article.Author.Bio,
		article.Author.ProfilePicture,
		pq.Array(article.Author.SocialLinks),
		article.Author.Followers,
		article.Author.Following,
	)

	if err != nil {
		psql.loggerService.PostLogMessage(err.Error())
		return nil, err
	}

	return article, nil
}

func (psql *postgresDBClient) GetArticleByID(article_id string) (*domain.Article, error) {
	query := fmt.Sprintf(`
		SELECT article_id, title, subtitle, introduction, body, tags, publish_date,author_id,author_name,author_bio, author_profile_pic, author_social_links,author_followers,author_following
		FROM %s WHERE article_id = $1`, psql.tablename)
	article := &domain.Article{}
	row := psql.db.QueryRow(query, article_id)
	err := row.Scan(
		&article.ArticleID,
		&article.Title,
		&article.Subtitle,
		&article.Introduction,
		&article.Body, pq.Array(&article.Tags),
		&article.PublishDate,
		&article.Author.ID,
		&article.Author.Name,
		&article.Author.Bio,
		&article.Author.ProfilePicture,
		pq.Array(&article.Author.SocialLinks),
		&article.Author.Followers,
		&article.Author.Following,
	)
	if err != nil {
		return nil, err
	}
	return article, nil
}

func (psql *postgresDBClient) GetArticlesByAuthor(author_id string) (*[]domain.Article, error) {
	query := fmt.Sprintf(`
		SELECT article_id, title, subtitle, introduction, body, tags, publish_date,author_id,author_name, author_bio, author_profile_pic, author_social_links,author_followers, author_following
		FROM %s WHERE author_id = $1`, psql.tablename)
	rows, err := psql.db.Query(query, author_id)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	articles := []domain.Article{}

	for rows.Next() {
		var article domain.Article
		err := rows.Scan(
			&article.ArticleID, &article.Title, &article.Subtitle, &article.Introduction,
			&article.Body, pq.Array(&article.Tags), &article.PublishDate, &article.Author.ID,
			&article.Author.Name, &article.Author.Bio, &article.Author.ProfilePicture,
			pq.Array(&article.Author.SocialLinks), &article.Author.Followers,
			&article.Author.Following,
		)
		if err != nil {
			return nil, err
		}
		articles = append(articles, article)
	}

	return &articles, nil
}

func (psql *postgresDBClient) GetArticlesByTag(tag string) (*[]domain.Article, error) {
	query := fmt.Sprintf(`
	SELECT article_id, title, subtitle, introduction, body, tags, publish_date,author_id,author_name, author_bio, author_profile_pic, author_social_links,author_followers, author_following
	FROM %s WHERE $1 = ANY(tags)`, psql.tablename)
	rows, err := psql.db.Query(query, tag)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []domain.Article

	for rows.Next() {
		var article domain.Article
		err := rows.Scan(
			&article.ArticleID, &article.Title, &article.Subtitle, &article.Introduction,
			&article.Body, pq.Array(&article.Tags), &article.PublishDate, &article.Author.ID,
			&article.Author.Name, &article.Author.Bio, &article.Author.ProfilePicture,
			pq.Array(&article.Author.SocialLinks), &article.Author.Followers,
			&article.Author.Following,
		)
		if err != nil {
			return nil, err
		}
		articles = append(articles, article)
	}

	return &articles, nil
}

func (psql *postgresDBClient) GetArticles() (*[]domain.Article, error) {
	query := fmt.Sprintf(`
	SELECT article_id, title, subtitle, introduction, body, tags, publish_date,author_id, author_name,author_bio, author_profile_pic, author_social_links,author_followers, author_following
	FROM %s `, psql.tablename)

	rows, err := psql.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	articles := []domain.Article{}
	for rows.Next() {
		var article domain.Article
		err := rows.Scan(
			&article.ArticleID, &article.Title, &article.Subtitle, &article.Introduction,
			&article.Body, pq.Array(&article.Tags), &article.PublishDate, &article.Author.ID,
			&article.Author.Name, &article.Author.Bio, &article.Author.ProfilePicture,
			pq.Array(&article.Author.SocialLinks), &article.Author.Followers,
			&article.Author.Following,
		)
		if err != nil {
			return nil, err
		}
		articles = append(articles, article)
	}

	return &articles, nil
}

func (psql *postgresDBClient) UpdateArticle(article *domain.Article) (*domain.Article, error) {
	_, err := psql.db.Exec(fmt.Sprintf(`
		UPDATE %s
		SET title=$1,
		subtitle=$2,
		introduction=$3,
		body=$4,
		tags=$5,
		publish_date=$6,
		author_id=$7,
		author_name=$8,
		author_bio=$9,
		author_profile_pic=$10,
		author_social_links=$11,
		author_followers=$12,
		author_following=$13
		WHERE article_id=$14`, psql.tablename),
		article.Title,
		article.Subtitle,
		article.Introduction,
		article.Body,
		pq.Array(article.Tags),
		article.PublishDate,
		article.Author.ID,
		article.Author.Name,
		article.Author.Bio,
		article.Author.ProfilePicture,
		pq.Array(article.Author.SocialLinks),
		article.Author.Followers,
		article.Author.Following,
		article.ArticleID,
	)

	if err != nil {
		return nil, err
	}

	return psql.GetArticleByID(article.ArticleID)
}

func (psql *postgresDBClient) DeleteArticle(article_id string) error {
	_, err := psql.db.Exec(fmt.Sprintf("DELETE FROM %s WHERE article_id = $1", psql.tablename), article_id)
	if err != nil {
		return err
	}
	return nil
}

func (psql *postgresDBClient) DeleteArticleAll() error {
	_, err := psql.db.Exec(fmt.Sprintf("DELETE FROM %s", psql.tablename))
	if err != nil {
		return err
	}
	return nil
}
