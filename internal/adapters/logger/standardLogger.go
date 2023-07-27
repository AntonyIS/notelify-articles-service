package logger

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type LoggerService interface {
	PostLogMessage(msg string)
}

type LoggerType struct {
	logServiceURL string
}

func (l LoggerType) PostLogMessage(message string) error {

	currentTime := time.Now()
	filename := fmt.Sprintf("%s.log", currentTime.Format("2006-01-02"))
	filepath := fmt.Sprintf("sysLogs/%s", filename)
	file, err := os.Create(filepath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	log.SetOutput(file)
	log.Println(message)

	// Create a new HTTP client
	client := &http.Client{}

	// Create a new POST request with the request body
	req, err := http.NewRequest("POST", l.logServiceURL, bytes.NewBuffer([]byte(message)))
	if err != nil {
		return err
	}

	// Set any headers you need for the request (e.g., Content-Type)
	req.Header.Set("Content-Type", "application/json")

	// Perform the request
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Read the response body
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func NewLoggerService(URL string) LoggerType {
	return LoggerType{
		logServiceURL: URL,
	}
}
