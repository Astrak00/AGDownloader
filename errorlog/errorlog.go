package errorlog

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// ErrorLogger handles logging of all errors during download operations
type ErrorLogger struct {
	file       *os.File
	logger     *log.Logger
	mu         sync.Mutex
	errorCount int
}

// ErrorType represents the type of error that occurred
type ErrorType string

const (
	ErrorTypeDownload        ErrorType = "DOWNLOAD"
	ErrorTypeCourseContent   ErrorType = "COURSE_CONTENT"
	ErrorTypeFileSystem      ErrorType = "FILE_SYSTEM"
	ErrorTypeNetwork         ErrorType = "NETWORK"
	ErrorTypeCourseRetrieval ErrorType = "COURSE_RETRIEVAL"
)

// New creates a new ErrorLogger with a timestamped log file
func New(dirPath string) (*ErrorLogger, error) {
	// Create error logs directory if it doesn't exist
	errorLogDir := filepath.Join(dirPath, "error_logs")
	if err := os.MkdirAll(errorLogDir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create error log directory: %v", err)
	}

	// Create log file with timestamp
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	logFileName := fmt.Sprintf("download_errors_%s.log", timestamp)
	logFilePath := filepath.Join(errorLogDir, logFileName)

	file, err := os.Create(logFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create error log file: %v", err)
	}

	logger := log.New(file, "", 0)

	// Write header
	logger.Printf("=======================================================\n")
	logger.Printf("AGDownloader Error Log\n")
	logger.Printf("Started: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	logger.Printf("=======================================================\n\n")

	return &ErrorLogger{
		file:   file,
		logger: logger,
	}, nil
}

// LogError logs an error with context information
func (el *ErrorLogger) LogError(errorType ErrorType, context string, err error) {
	if el == nil || el.logger == nil {
		return
	}

	el.mu.Lock()
	defer el.mu.Unlock()

	el.errorCount++
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	el.logger.Printf("[%s] [%s] %s\n", timestamp, errorType, context)
	el.logger.Printf("Error: %v\n", err)
	el.logger.Printf("-------------------------------------------------------\n\n")
}

// LogErrorWithDetails logs an error with additional details
func (el *ErrorLogger) LogErrorWithDetails(errorType ErrorType, context string, err error, details map[string]string) {
	if el == nil || el.logger == nil {
		return
	}

	el.mu.Lock()
	defer el.mu.Unlock()

	el.errorCount++
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	el.logger.Printf("[%s] [%s] %s\n", timestamp, errorType, context)
	el.logger.Printf("Error: %v\n", err)

	if len(details) > 0 {
		el.logger.Printf("Details:\n")
		for key, value := range details {
			el.logger.Printf("  %s: %s\n", key, value)
		}
	}

	el.logger.Printf("-------------------------------------------------------\n\n")
}

// GetErrorCount returns the total number of errors logged
func (el *ErrorLogger) GetErrorCount() int {
	if el == nil {
		return 0
	}

	el.mu.Lock()
	defer el.mu.Unlock()
	return el.errorCount
}

// Close closes the error log file and writes summary
func (el *ErrorLogger) Close() error {
	if el == nil || el.file == nil {
		return nil
	}

	el.mu.Lock()
	defer el.mu.Unlock()

	// Write summary
	el.logger.Printf("\n=======================================================\n")
	el.logger.Printf("Summary\n")
	el.logger.Printf("Total errors logged: %d\n", el.errorCount)
	el.logger.Printf("Completed: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	el.logger.Printf("=======================================================\n")

	return el.file.Close()
}

// GetLogFilePath returns the path to the log file
func (el *ErrorLogger) GetLogFilePath() string {
	if el == nil || el.file == nil {
		return ""
	}
	return el.file.Name()
}
