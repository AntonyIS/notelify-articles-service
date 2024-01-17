package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AntonyIS/notelify-articles-service/internal/core/domain"
	"github.com/AntonyIS/notelify-articles-service/internal/core/ports"
)

type loggerService struct {
	logger ports.Logger
}

func NewLoggerService(logger ports.Logger) *loggerService {
	return &loggerService{logger: logger}
}

func (l *loggerService) Info(message string) {
	l.logger.Info(message)
}

func (l *loggerService) Error(message string) {
	l.logger.Error(message)
}

func (l *loggerService) Close() {
	l.logger.Close()
}

type loggingService struct {
	loggerURL string
}

func NewLoggingService(loggerURL string) *loggingService {
	svc := loggingService{
		loggerURL: loggerURL,
	}
	return &svc
}

func (svc *loggingService) CreateLog(logEntry domain.LogMessage) {
	// Marshal the struct into JSON
	payloadBytes, err := json.Marshal(logEntry)

	if err != nil {
		fmt.Println("Error encoding JSON payload:", err)
		return
	}

	// Create a new POST request with the JSON payload
	resp, err := http.Post(svc.loggerURL, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		fmt.Println("Error making POST request:", err)
		return
	}
	defer resp.Body.Close()

}
