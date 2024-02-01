package postgres

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	appConfig "github.com/AntonyIS/notelify-articles-service/config"
	"github.com/AntonyIS/notelify-articles-service/internal/core/domain"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type postgresDBClient struct {
	db        *sql.DB
	tablename string
}

func NewPostgresClient(conf appConfig.Config) (*postgresDBClient, error) {
	dbname := conf.POSTGRES_DB
	tablename := conf.ARTICLE_TABLE
	user := conf.POSTGRES_USER
	password := conf.POSTGRES_PASSWORD
	port := conf.POSTGRES_PORT
	host := conf.POSTGRES_HOST

	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", host, port, user, dbname, password)

	db, err := sql.Open("postgres", dsn)

	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil {
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
			updated_date TIMESTAMP,
			author JSONB NOT NULL,
			author_id VARCHAR(255)
	)
	`, tablename)

	_, err = db.Exec(queryString)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	return &postgresDBClient{db: db, tablename: tablename}, nil
}

func (psql *postgresDBClient) CreateArticle(article *domain.Article) (*domain.Article, error) {
	// Convert Author struct to JSON string
	authorJSON, err := json.Marshal(article.Author)
	if err != nil {
		return nil, err
	}
	query := fmt.Sprintf(`
		INSERT INTO %s (
			article_id,
			title,
			subtitle,
			introduction,
			body,
			tags,
			publish_date,
			updated_date,
			author,
			author_id
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`, psql.tablename)
	_, err = psql.db.Exec(
		query,
		article.ArticleID,
		article.Title,
		article.Subtitle,
		article.Introduction,
		article.Body,
		pq.Array(article.Tags),
		article.PublishDate,
		article.UpdatedDate,
		string(authorJSON), // Convert Author struct to JSON string
		article.AuthorID,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return article, nil
}

func (psql *postgresDBClient) GetArticleByID(article_id string) (*domain.Article, error) {
	query := fmt.Sprintf(`
		SELECT 
			article_id,
			title,
			subtitle,
			introduction,
			body,
			tags,
			publish_date,
			updated_date,
			author,
			author_id
		FROM %s 
		WHERE article_id = $1`,
		psql.tablename,
	)
	article := &domain.Article{}
	row := psql.db.QueryRow(query, article_id)
	var authorJSON []byte
	err := row.Scan(
		&article.ArticleID,
		&article.Title,
		&article.Subtitle,
		&article.Introduction,
		&article.Body,
		pq.Array(&article.Tags),
		&article.PublishDate,
		&article.UpdatedDate,
		&authorJSON,
		&article.AuthorID,
	)

	if err != nil {
		return nil, err
	}

	// Unmarshal JSONB data into Author struct
	err = json.Unmarshal(authorJSON, &article.Author)
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
			updated_date,
			author,
			author_id
		FROM %s 
		WHERE author_id = $1`, psql.tablename)
	rows, err := psql.db.Query(query, author_id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	articles := []domain.Article{}

	for rows.Next() {
		var article domain.Article
		var authorJSON []byte
		err := rows.Scan(
			&article.ArticleID,
			&article.Title,
			&article.Subtitle,
			&article.Introduction,
			&article.Body,
			pq.Array(&article.Tags),
			&article.PublishDate,
			&article.UpdatedDate,
			&authorJSON,
			&article.AuthorID,
		)
		if err != nil {
			return nil, err
		}

		// Unmarshal JSONB data into Author struct
		err = json.Unmarshal(authorJSON, &article.Author)
		if err != nil {
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
			updated_date,
			author,
			author_id
		FROM %s 
		WHERE $1 = ANY(tags)`,
		psql.tablename,
	)
	rows, err := psql.db.Query(query, tag)

	if err != nil {
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
			&article.UpdatedDate,
			&article.Author,
			&article.AuthorID,
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
	SELECT 
		article_id,
		title,
		subtitle,
		introduction,
		body,
		tags,
		publish_date,
		updated_date,
		author,
		author_id
	FROM %s `, psql.tablename)

	rows, err := psql.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	articles := []domain.Article{}
	for rows.Next() {
		var article domain.Article
		var authorJSON []byte
		err := rows.Scan(
			&article.ArticleID,
			&article.Title,
			&article.Subtitle,
			&article.Introduction,
			&article.Body,
			pq.Array(&article.Tags),
			&article.PublishDate,
			&article.UpdatedDate,
			&authorJSON,
			&article.AuthorID,
		)
		if err != nil {
			return nil, err
		}

		// Unmarshal JSONB data into Author struct
		err = json.Unmarshal(authorJSON, &article.Author)
		if err != nil {
			return nil, err
		}
		articles = append(articles, article)
	}

	return &articles, nil
}

func (psql *postgresDBClient) UpdateArticle(article_id string, article *domain.Article) (*domain.Article, error) {
	res, err := psql.GetArticleByID(article_id)

	if err != nil {
		return nil, err
	}

	res.Title = article.Title
	res.Subtitle = article.Subtitle
	res.Introduction = article.Introduction
	res.Body = article.Body
	res.Tags = article.Tags
	res.PublishDate = article.PublishDate
	res.UpdatedDate = article.UpdatedDate
	res.AuthorID = article.AuthorID

	query := fmt.Sprintf(`
	UPDATE 
		%s 
	SET 
		title=$1,
		subtitle=$2,
		introduction=$3,
		body=$4,
		tags=$5,
		publish_date=$6,
		updated_date=$7,
		author=$8
		author_id=$9	
	WHERE 
		article_id=$8`,
		psql.tablename,
	)

	_, err = psql.db.Exec(
		query,
		res.Title,
		res.Subtitle,
		res.Introduction,
		res.Body,
		pq.Array(res.Tags),
		res.PublishDate,
		res.UpdatedDate,
		res.Author,
		res.ArticleID,
	)

	if err != nil {
		return nil, err
	}

	return psql.GetArticleByID(article_id)
}

func (psql *postgresDBClient) DeleteArticle(article_id string) error {
	_, err := psql.GetArticleByID(article_id)

	if err != nil {
		return err
	}
	query := fmt.Sprintf(`DELETE FROM %s WHERE article_id = $1`, psql.tablename)

	_, err = psql.db.Exec(query, article_id)
	if err != nil {
		return err
	}
	return nil
}

func (psql *postgresDBClient) DeleteArticleAll() error {
	articles, err := psql.GetArticles()
	if err != nil {
		return err
	}
	if len(*articles) == 0 {
		return errors.New("no Articles to delete")
	}
	query := fmt.Sprintf(`DELETE FROM %s `, psql.tablename)
	_, err = psql.db.Exec(query)

	if err != nil {
		return err
	}
	return nil
}
