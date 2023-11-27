package ports

type Logger interface {
	Info(message string)
	Error(message string)
	Close()
}