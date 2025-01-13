package log

import (
	"log"
	"log/slog"
	"os"
)

// Declare a global logger variable
var logger *slog.Logger

// Initialize the logger (called once at the start of the program)
func InitLogger() func() error {
	// Open the log file for writing (create if not exists, append to the file)
	logFile, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}

	// Create a slog handler that writes logs to the file
	fileHandler := slog.NewJSONHandler(logFile, nil) // You can use NewTextHandler for plain text logs

	// Initialize the global logger with the file handler
	logger = slog.New(fileHandler)

	return logFile.Close
}

// Global function to access the logger
func GetLogger() *slog.Logger {
	return logger
}
