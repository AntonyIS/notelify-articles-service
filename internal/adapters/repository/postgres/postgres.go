package postgres

import (
	"database/sql"
	"fmt"

	"github.com/AntonyIS/notlify-content-svc/config"
	"github.com/AntonyIS/notlify-content-svc/internal/adapters/logger"
	"github.com/AntonyIS/notlify-content-svc/internal/core/domain"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/rds/rdsutils"
	_ "github.com/lib/pq"
)

type PostgresDBClient struct {
	db            *sql.DB
	tablename     string
	loggerService logger.LoggerType
}

func NewPostgresClient(config config.Config, logger logger.LoggerType) (*PostgresDBClient, error) {
	databaseName := config.DatabaseName
	databaseContentTable := config.ContentTable
	databaseUser := config.DatabaseUser
	databasePassword := config.DatabasePassword
	databasePort := config.DatabasePort
	databaseHost := config.DatabaseHost
	databaseRegion := config.AWS_DEFAULT_REGION
	var dsn string

	if config.Env == "dev" {
		dsn = fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
			databaseHost,
			databasePort,
			databaseUser,
			databaseName,
			databasePassword,
		)
	} else {
		dbEndpoint := fmt.Sprintf("%s:%s", databaseHost, databasePort)
		creds := credentials.NewEnvCredentials()
		authToken, err := rdsutils.BuildAuthToken(dbEndpoint, databaseRegion, databaseUser, creds)

		if err != nil {
			logger.PostLogMessage(err.Error())
			return nil, err
		}

		dsn = fmt.Sprintf("%s:%s@tcp(%s)/%s?tls=true&allowCleartextPasswords=true",
			databaseUser, authToken, dbEndpoint, databaseName,
		)
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

	// Create content table
	err = migrate(db, databaseContentTable)
	if err != nil {
		logger.PostLogMessage(err.Error())
		return nil, err

	}

	return &PostgresDBClient{db: db, tablename: databaseContentTable, loggerService: logger}, nil
}

func (psql *PostgresDBClient) CreateContent(content *domain.Content) (*domain.Content, error) {
	queryString := fmt.Sprintf(
		`INSERT INTO %s 
			(content_id,creator_id,title,body,publication_date) 
			VALUES 
			($1, $2, $3, $4, $5)`,
		psql.tablename)
	_, err := psql.db.Exec(queryString, content.ContentId, content.CreatorId, content.Title, content.Body, content.PublicationDate)

	if err != nil {
		psql.loggerService.PostLogMessage(err.Error())
		return nil, err
	}

	return content, nil
}

func (psql *PostgresDBClient) ReadContent(id string) (*domain.Content, error) {
	var content domain.Content
	queryString := fmt.Sprintf(`SELECT content_id,creator_id,title,body,publication_date FROM %s WHERE content_id=$1`, psql.tablename)
	err := psql.db.QueryRow(queryString, id).Scan(&content.ContentId, &content.CreatorId, &content.Title, &content.Body, &content.PublicationDate)
	if err != nil {
		psql.loggerService.PostLogMessage(err.Error())
		return nil, err
	}

	return &content, nil
}

func (psql *PostgresDBClient) ReadContents() ([]domain.Content, error) {
	rows, err := psql.db.Query(fmt.Sprintf("SELECT * FROM %s", psql.tablename))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	contents := []domain.Content{}
	for rows.Next() {
		var content domain.Content

		if err := rows.Scan(&content.ContentId, &content.CreatorId, &content.Title, &content.Body, &content.PublicationDate); err != nil {
			psql.loggerService.PostLogMessage(err.Error())
			return nil, err
		}

		contents = append(contents, content)

	}
	return contents, nil
}

func (psql *PostgresDBClient) UpdateContent(content *domain.Content) (*domain.Content, error) {
	queryString := fmt.Sprintf(`UPDATE %s SET 
		title = $2,
		body = $3
	`, psql.tablename)

	_, err := psql.db.Exec(queryString, content.Title, content.Body)
	if err != nil {
		psql.loggerService.PostLogMessage(err.Error())
		return nil, err
	}
	return content, nil
}

func (psql *PostgresDBClient) DeleteContent(id string) (string, error) {

	queryString := fmt.Sprintf(`DELETE FROM %s WHERE content_id = $1`, psql.tablename)
	_, err := psql.db.Exec(queryString, id)
	if err != nil {
		psql.loggerService.PostLogMessage(err.Error())
		return "", err
	}
	return "Entity deleted successfully", nil
}

func migrate(db *sql.DB, contentTable string) error {
	// Creates new contentTable if does not exists
	queryString := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			content_id VARCHAR(255) PRIMARY KEY UNIQUE,
			creator_id VARCHAR(255) NOT NULL,
			title VARCHAR(255) NOT NULL,
			body VARCHAR(255) NOT NULL,
			publication_date DATE NOT NULL
	)
	`, contentTable)

	_, err := db.Exec(queryString)
	if err != nil {

		return err
	}

	return nil

}
