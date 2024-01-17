package ports

import "github.com/AntonyIS/notelify-articles-service/internal/core/domain"

type Logger interface {
	Info(message string)
	Error(message string)
	Close()
}

type LoggingService interface {
	CreateLog(LogEntry domain.LogMessage)
}