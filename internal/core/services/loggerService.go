package services

import (
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
