package postgres

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

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
	databaseUserTable := config.UserTable
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
	migrate(db, databaseUserTable)

	return &PostgresDBClient{db: db, tablename: databaseUserTable, loggerService: logger}, nil
}

func (psql *PostgresDBClient) CreateUser(user *domain.User) (*domain.User, error) {
	queryString := fmt.Sprintf(
		`INSERT INTO %s 
			(id,firstname,lastname,email,password,handle,about,profile_image,following,followers) 
			VALUES 
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		psql.tablename)

	_, err := psql.db.Exec(queryString, user.Id, user.Firstname, user.Lastname, user.Email, user.Password, user.Handle, user.About, user.ProfileImage, user.Following, user.Followers)

	if err != nil {
		psql.loggerService.PostLogMessage(err.Error())
		return nil, err
	}

	return user, nil
}

func (psql *PostgresDBClient) ReadUser(id string) (*domain.User, error) {
	var user domain.User
	queryString := fmt.Sprintf(`SELECT id,firstname, lastname,email, handle,about,profile_image,following, followers FROM %s WHERE id=$1`, psql.tablename)
	err := psql.db.QueryRow(queryString, id).Scan(&user.Id, &user.Firstname, &user.Lastname, &user.Email, &user.Handle, &user.About, &user.ProfileImage, &user.Following, &user.Followers)
	if err != nil {
		psql.loggerService.PostLogMessage(err.Error())
		return nil, err
	}
	contents, err := psql.readUserContent(id)

	if err != nil {
		psql.loggerService.PostLogMessage(err.Error())
		return nil, err
	}
	user.Contents = contents
	return &user, nil
}

func (psql *PostgresDBClient) ReadUserWithEmail(email string) (*domain.User, error) {
	var user domain.User
	queryString := fmt.Sprintf(`SELECT id,firstname, lastname,email, handle,about,profile_image,following, followers FROM %s WHERE email=$1`, psql.tablename)
	err := psql.db.QueryRow(queryString, email).Scan(&user.Id, &user.Firstname, &user.Lastname, &user.Email, &user.Handle, &user.About, &user.ProfileImage, &user.Following, &user.Followers)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (psql *PostgresDBClient) ReadUsers() ([]domain.User, error) {
	rows, err := psql.db.Query(fmt.Sprintf("SELECT * FROM %s", psql.tablename))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []domain.User{}
	for rows.Next() {
		var user domain.User

		if err := rows.Scan(&user.Id, &user.Firstname, &user.Lastname, &user.Email, &user.Password, &user.Handle, &user.About, &user.ProfileImage, &user.Following, &user.Followers); err != nil {
			psql.loggerService.PostLogMessage(err.Error())
			return nil, err
		}

		users = append(users, user)

	}
	return users, nil
}

func (psql *PostgresDBClient) UpdateUser(user *domain.User) (*domain.User, error) {
	queryString := fmt.Sprintf(`UPDATE %s SET 
		firstname = $2,
		lastname = $3,
		handle = $4,
		about = $5,
		profile_image = $6,
		following = $7,
		followers = $8
		WHERE id =$1

	`, psql.tablename)

	_, err := psql.db.Exec(queryString, user.Id, user.Firstname, user.Lastname, user.Handle, user.About, user.ProfileImage, user.Following, user.Followers)
	if err != nil {
		psql.loggerService.PostLogMessage(err.Error())
		return nil, err
	}
	return user, nil
}

func (psql *PostgresDBClient) DeleteUser(id string) (string, error) {

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
