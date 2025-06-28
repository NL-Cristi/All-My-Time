package logger

import (
	"log"
	"os"
)

type Level int

const (
	DEBUG Level = iota
	INFO
	ERROR
)

type Logger struct {
	debugLogger *log.Logger
	infoLogger  *log.Logger
	errorLogger *log.Logger
	level       Level
}

func NewLogger(level string) *Logger {
	var logLevel Level
	switch level {
	case "debug":
		logLevel = DEBUG
	case "info":
		logLevel = INFO
	case "error":
		logLevel = ERROR
	default:
		logLevel = INFO
	}

	return &Logger{
		debugLogger: log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile),
		infoLogger:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		errorLogger: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
		level:       logLevel,
	}
}

func (l *Logger) Debug(format string, v ...interface{}) {
	if l.level <= DEBUG {
		l.debugLogger.Printf(format, v...)
	}
}

func (l *Logger) Info(format string, v ...interface{}) {
	if l.level <= INFO {
		l.infoLogger.Printf(format, v...)
	}
}

func (l *Logger) Error(format string, v ...interface{}) {
	if l.level <= ERROR {
		l.errorLogger.Printf(format, v...)
	}
}
