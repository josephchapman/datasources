package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

// Custom log entry structure
type LogEntry struct {
	Application string `json:"application"`
	Timestamp   string `json:"timestamp"`
	Stream      string `json:"stream"`
	Message     string `json:"message"`
}

// JSONLogger is a custom logger that outputs JSON-formatted log entries
type JSONLogger struct {
	logger *log.Logger
	stream string
}

// NewJSONLogger creates a new JSONLogger
func NewJSONLoggerOut() *JSONLogger {
	return &JSONLogger{
		logger: log.New(os.Stdout, "", 0), // Disable default log flags
		stream: "stdout",
	}
}

// NewJSONLogger creates a new JSONLogger
func NewJSONLoggerErr() *JSONLogger {
	return &JSONLogger{
		logger: log.New(os.Stderr, "", 0), // Disable default log flags
		stream: "stderr",
	}
}

// Println logs a message in JSON format
func (jl *JSONLogger) Println(v ...interface{}) {
	entry := LogEntry{
		Application: "prometheus-exporter-weather",
		Timestamp:   time.Now().Format(time.RFC3339),
		Stream:      jl.stream,
		Message:     fmt.Sprint(v...),
	}
	jsonData, err := json.Marshal(entry)
	if err != nil {
		log.Println("Error marshaling log entry:", err)
		return
	}
	jl.logger.Println(string(jsonData))
}

// WrapError logs the error in JSON format and returns the original error
func WrapError(err error) error {
	if err != nil {
		jsonLogger := NewJSONLoggerErr()
		jsonLogger.Println(err)
	}
	return err
}

// WrapOut logs the string in JSON format and returns the original string
func WrapOut(str string) string {
	jsonLogger := NewJSONLoggerOut()
	jsonLogger.Println(str)
	return str
}
