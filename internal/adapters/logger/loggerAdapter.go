package logger

import (
	"io"
	"log"
	"os"
)

type ConsoleFileLogger struct {
	logger *log.Logger
	logFile *os.File
}

func NewLogger() ConsoleFileLogger {
	path := "internal/adapters/logger/logs.log"
	logFile, err := os.Create(path)
	if err != nil {
		log.Fatal("Error creating a file: ", err)
	}

	// Create a multi-writer that writes to both os.Stdout and the log file
	multiWriter := io.MultiWriter(os.Stdout, logFile)

	// Create a custom logger with a custom log.Formatter
	logger := log.New(multiWriter, "", log.Ldate|log.Ltime)

	return ConsoleFileLogger{logger: logger}

}

func (l *ConsoleFileLogger) Info(message string) {
	l.logger.Println(message)
}

func (l *ConsoleFileLogger) Error(message string) {
	l.logger.Println(message)
}

func (cfl *ConsoleFileLogger) Close() {
    if cfl.logFile != nil {
        cfl.logFile.Close()
    }
}
