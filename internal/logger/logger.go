// Package logger provides a centralized logging system for the application.
// It handles file-based logging with rotation and different log levels.
package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

var (
	// InfoLogger handles informational messages
	InfoLogger *log.Logger
	// ErrorLogger handles error messages
	ErrorLogger *log.Logger
	// DebugLogger handles debug messages (can be disabled in production)
	DebugLogger *log.Logger

	logFile      *os.File
	debugEnabled bool
)

// Init initializes the logging system with file output
// Creates log files in the specified directory
func Init(logDir string, enableDebug bool) error {
	debugEnabled = enableDebug

	// Create logs directory if it doesn't exist
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %v", err)
	}

	// Create log file with timestamp
	timestamp := time.Now().Format("2006-01-02")
	logFileName := fmt.Sprintf("app_%s.log", timestamp)
	logPath := filepath.Join(logDir, logFileName)

	var err error
	logFile, err = os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("failed to open log file: %v", err)
	}

	// Create multi-writer to write to both file and stdout (for important messages)
	multiWriter := io.MultiWriter(logFile, os.Stdout)

	// Initialize loggers
	InfoLogger = log.New(multiWriter, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(multiWriter, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	if debugEnabled {
		// Debug logs only to file, not to stdout to reduce noise
		DebugLogger = log.New(logFile, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
	} else {
		// If debug is disabled, create a no-op logger
		DebugLogger = log.New(io.Discard, "", 0)
	}

	InfoLogger.Printf("Logger initialized. Log file: %s, Debug enabled: %v", logPath, debugEnabled)
	return nil
}

// Close closes the log file
func Close() {
	if logFile != nil {
		logFile.Close()
	}
}

// Info logs informational messages
func Info(format string, v ...interface{}) {
	if InfoLogger != nil {
		InfoLogger.Output(2, fmt.Sprintf(format, v...))
	}
}

// Error logs error messages
func Error(format string, v ...interface{}) {
	if ErrorLogger != nil {
		ErrorLogger.Output(2, fmt.Sprintf(format, v...))
	}
}

// Debug logs debug messages (only if debug is enabled)
func Debug(format string, v ...interface{}) {
	if debugEnabled && DebugLogger != nil {
		DebugLogger.Output(2, fmt.Sprintf(format, v...))
	}
}

// CleanOldLogs removes log files older than the specified number of days
func CleanOldLogs(logDir string, daysToKeep int) error {
	entries, err := os.ReadDir(logDir)
	if err != nil {
		return fmt.Errorf("failed to read log directory: %v", err)
	}

	cutoffTime := time.Now().AddDate(0, 0, -daysToKeep)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		// Only process log files
		if filepath.Ext(entry.Name()) != ".log" {
			continue
		}

		filePath := filepath.Join(logDir, entry.Name())
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			continue
		}

		if fileInfo.ModTime().Before(cutoffTime) {
			if err := os.Remove(filePath); err != nil {
				Error("Failed to remove old log file %s: %v", filePath, err)
			} else {
				Info("Removed old log file: %s", entry.Name())
			}
		}
	}

	return nil
}
