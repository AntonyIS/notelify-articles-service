package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	appConfig "github.com/AntonyIS/notlify-content-svc/config"
	"github.com/AntonyIS/notlify-content-svc/internal/adapters/logger"
	"github.com/AntonyIS/notlify-content-svc/internal/core/domain"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
	_ "github.com/lib/pq"
)

type PostgresDBClient struct {
	db            *sql.DB
	tablename     string
	loggerService logger.LoggerType
}

func NewPostgresClient(appConfig appConfig.Config, logger logger.LoggerType) (*PostgresDBClient, error) {
	dbname := appConfig.DatabaseName
	tablename := appConfig.ContentTable
	user := appConfig.DatabaseUser
	password := appConfig.DatabasePassword
	port := appConfig.DatabasePort
	host := appConfig.DatabaseHost
	region := appConfig.AWS_DEFAULT_REGION
	rdsInstanceIdentifier := appConfig.RDSInstanceIdentifier
	var dsn string

	if appConfig.Env == "dev" {
		dsn = fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable", host, port, user, dbname, password)
	} else {
		awsSession := session.Must(session.NewSession(&aws.Config{
			Region: aws.String(region),
		}))
		rdsClient := rds.New(awsSession)

		// Describe the DB instance to get its endpoint
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

		dsn = fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=require", host, port, dbname, user, password)
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

	err = migrateDB(db, tablename)
	if err != nil {
		logger.PostLogMessage(err.Error())
		return nil, err

	}

	return &PostgresDBClient{db: db, tablename: tablename, loggerService: logger}, nil
}

func (psql *PostgresDBClient) CreateContent(content *domain.Content) (*domain.Content, error) {
	queryString := fmt.Sprintf(
		`INSERT INTO %s 
			(content_id,creator_id,title,body,publication_date) 
			VALUES 
			($1, $2, $3, $4, $5)`,
		psql.tablename,
	)
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
		psql.loggerService.PostLogMessage(fmt.Sprintf("content with id [%s] not found: %s", id, err.Error()))
		return nil, errors.New(fmt.Sprintf("content with id [%s] not found", id))
	}

	if err != nil {
		psql.loggerService.PostLogMessage(err.Error())
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

		if err != nil {
			psql.loggerService.PostLogMessage(err.Error())
		}
		contents = append(contents, content)

	}
	return contents, nil
}

func (psql *PostgresDBClient) UpdateContent(content *domain.Content) (*domain.Content, error) {
	queryString := fmt.Sprintf(`
		UPDATE %s 
		SET title = $1, body = $2
		WHERE content_id = $3
	`, psql.tablename)

	_, err := psql.db.Exec(queryString, content.Title, content.Body, content.ContentId)
	if err != nil {
		psql.loggerService.PostLogMessage(err.Error())
		return nil, err
	}
	if err != nil {
		psql.loggerService.PostLogMessage(err.Error())
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
	// Delete
	return "Entity deleted successfully", nil
}

func (psql *PostgresDBClient) DeleteAllContent() (string, error) {
	queryString := fmt.Sprintf(`DELETE FROM %s`, psql.tablename)
	_, err := psql.db.Exec(queryString)
	if err != nil {
		return "", err
	}

	return "All items deletes successfully", nil

}

func migrateDB(db *sql.DB, contentTable string) error {
	queryString := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			content_id VARCHAR(255) PRIMARY KEY UNIQUE,
			creator_id VARCHAR(255) NOT NULL,
			title VARCHAR(255) NOT NULL,
			body VARCHAR(2000) NOT NULL,
			publication_date DATE NOT NULL
	)
	`, contentTable)

	_, err := db.Exec(queryString)
	if err != nil {

		return err
	}

	return nil

}
