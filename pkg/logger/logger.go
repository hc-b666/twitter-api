package logger

import (
	"fmt"
	"log"
	"os"
)

type Logger struct {
	*os.File
	*log.Logger
}

func NewLogger(path string) (*Logger, error) {
	var file *os.File
	const mode = 0o644
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, mode)
	if err != nil {
		return nil, fmt.Errorf("log error: %w", err)
	}

	return &Logger{file, log.New(file, "[LOG] ", log.LstdFlags)}, nil
}

func (l *Logger) Info(v ...any) {
	l.Println("[INFO]", fmt.Sprint(v...))
}

func (l *Logger) Error(v ...any) {
	l.Println("[ERROR]", fmt.Sprint(v...))
}

func (l *Logger) Done() {
	if err := l.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "cannot close logger: %v", err)
	}
}
