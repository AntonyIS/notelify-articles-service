package postgres

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

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
	databaseContentTable := config.UserTable
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
			return nil, err
		}

		dsn = fmt.Sprintf("%s:%s@tcp(%s)/%s?tls=true&allowCleartextPasswords=true",
			databaseUser, authToken, dbEndpoint, databaseName,
		)
	}

	db, err := sql.Open("postgres", dsn)

	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// Create users table
	migrate(db, databaseContentTable)

	return &PostgresDBClient{db: db, tablename: databaseContentTable, loggerService: logger}, nil
}

func (psql *PostgresDBClient) CreateUser(content *domain.Content) (*domain.Content, error) {
	queryString := fmt.Sprintf(
		`INSERT INTO %s 
			(content_id,creator_id,title,body,images,vidoes,publication_date) 
			VALUES 
			($1, $2, $3, $4, $5, $6, $7)`,
		psql.tablename)
	_, err := psql.db.Exec(queryString, content.ContentId, content.CreatorId, content.Title, content.Body, content.Images, content.Videos, content.PublicationDate)

	if err != nil {
		psql.loggerService.PostLogMessage(err.Error())
		return nil, err
	}

	return content, nil
}

func (psql *PostgresDBClient) ReadContent(id string) (*domain.Content, error) {
	var content domain.Content
	queryString := fmt.Sprintf(`SELECT content_id,creator_id,title,body,images,vidoes,publication_date FROM %s WHERE content_id=$1`, psql.tablename)
	err := psql.db.QueryRow(queryString, id).Scan(&content.ContentId, &content.CreatorId, &content.Title, &content.Body, &content.Images, &content.Videos, &content.PublicationDate)
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

		if err := rows.Scan(&content.ContentId, &content.CreatorId, &content.Title, &content.Body, &content.Images, &content.Videos, &content.PublicationDate); err != nil {
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
		body = $3,
		images = $4,
		videos = $5,
		publication_date = $6
	`, psql.tablename)

	_, err := psql.db.Exec(queryString, content.Title, content.Body, content.Images, content.Videos, content.PublicationDate)
	if err != nil {
		psql.loggerService.PostLogMessage(err.Error())
		return nil, err
	}
	return content, nil
}

func (psql *PostgresDBClient) DeleteContent(id string) (string, error) {

	queryString := fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, psql.tablename)
	_, err := psql.db.Exec(queryString, id)
	if err != nil {
		psql.loggerService.PostLogMessage(err.Error())
		return "", err
	}
	return "Entity deleted successfully", nil
}
func (psql *PostgresDBClient) readUserContent(userId string) ([]domain.Content, error) {
	// URL of the API or website you want to request data from
	url := "http://127.0.0.1:5000"

	// Send GET request
	response, err := http.Get(url)
	if err != nil {
		fmt.Println(url)
		psql.loggerService.PostLogMessage(err.Error())
		return nil, err
	}
	defer response.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {

		psql.loggerService.PostLogMessage(err.Error())
		return nil, err
	}

	// Convert the response body to a string and print it
	var content []domain.Content
	err = json.Unmarshal(body, &content)
	if err != nil {
		psql.loggerService.PostLogMessage(err.Error())
		return nil, err
	}
	return content, nil
}

func migrate(db *sql.DB, userTable string) error {
	// Creates new usertable if does not exists
	queryString := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			id VARCHAR(255) PRIMARY KEY UNIQUE,
			firstname VARCHAR(255) NOT NULL,
			lastname VARCHAR(255) NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL,
			password VARCHAR(255) UNIQUE NOT NULL,
			handle VARCHAR(255),
			about TEXT,
			profile_image varchar(255),
			Following int,
			Followers int
	)
	`, userTable)

	_, err := db.Exec(queryString)
	if err != nil {
		return err
	}

	return nil

}
