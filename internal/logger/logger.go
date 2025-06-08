// Package logger provides structured logging functionality with file rotation and multiple output formats.
// It wraps logrus to provide application-specific logging methods for the Mark Master Sheet Consolidator.
package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"mark-master-sheet/internal/config"
)

// Logger wraps logrus with additional functionality
type Logger struct {
	*logrus.Logger
	config *config.LoggingConfig
}

// NewLogger creates a new logger instance with the given configuration
func NewLogger(cfg *config.LoggingConfig, logDir string) (*Logger, error) {
	logger := logrus.New()

	// Set log level
	level, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		return nil, fmt.Errorf("invalid log level %s: %w", cfg.Level, err)
	}
	logger.SetLevel(level)

	// Set formatter
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// Configure output
	if cfg.FileOutput {
		// Ensure log directory exists
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create log directory: %w", err)
		}

		// Setup file rotation
		logFile := filepath.Join(logDir, fmt.Sprintf("mark-master-sheet-%s.log",
			time.Now().Format("2006-01-02")))

		fileWriter := &lumberjack.Logger{
			Filename:   logFile,
			MaxSize:    cfg.MaxFileSizeMB,
			MaxBackups: cfg.MaxBackupFiles,
			MaxAge:     cfg.MaxAgeDays,
			Compress:   true,
		}

		if cfg.ConsoleOutput {
			// Output to both file and console
			logger.SetOutput(io.MultiWriter(os.Stdout, fileWriter))
		} else {
			// Output to file only
			logger.SetOutput(fileWriter)
		}
	} else if cfg.ConsoleOutput {
		// Output to console only
		logger.SetOutput(os.Stdout)
	}

	return &Logger{
		Logger: logger,
		config: cfg,
	}, nil
}

// LogProcessingStart logs the start of processing
func (l *Logger) LogProcessingStart(totalFiles int) {
	l.WithFields(logrus.Fields{
		"total_files": totalFiles,
		"timestamp":   time.Now(),
	}).Info("Starting mark consolidation process")
}

// LogProcessingEnd logs the end of processing with summary
func (l *Logger) LogProcessingEnd(summary interface{}) {
	l.WithFields(logrus.Fields{
		"summary":   summary,
		"timestamp": time.Now(),
	}).Info("Mark consolidation process completed")
}

// LogFileProcessed logs successful file processing
func (l *Logger) LogFileProcessed(filePath, studentID string, markCount int, duration time.Duration) {
	l.WithFields(logrus.Fields{
		"file_path":  filePath,
		"student_id": studentID,
		"mark_count": markCount,
		"duration":   duration,
	}).Info("File processed successfully")
}

// LogFileError logs file processing errors
func (l *Logger) LogFileError(filePath string, err error, stage string) {
	l.WithFields(logrus.Fields{
		"file_path": filePath,
		"stage":     stage,
		"error":     err.Error(),
	}).Error("File processing failed")
}

// LogStudentNotFound logs when a student ID is not found in master sheet
func (l *Logger) LogStudentNotFound(studentID, filePath string, suggestions []string) {
	fields := logrus.Fields{
		"student_id": studentID,
		"file_path":  filePath,
	}

	if len(suggestions) > 0 {
		fields["suggestions"] = suggestions
	}

	l.WithFields(fields).Warn("Student ID not found in master sheet")
}

// LogBackupCreated logs successful backup creation
func (l *Logger) LogBackupCreated(originalPath, backupPath string) {
	l.WithFields(logrus.Fields{
		"original_path": originalPath,
		"backup_path":   backupPath,
	}).Info("Backup created successfully")
}

// LogValidationError logs validation errors
func (l *Logger) LogValidationError(filePath, field, value, message string) {
	l.WithFields(logrus.Fields{
		"file_path": filePath,
		"field":     field,
		"value":     value,
		"message":   message,
	}).Warn("Validation error")
}

// LogProgress logs processing progress
func (l *Logger) LogProgress(processed, total int, currentFile string) {
	percentage := float64(processed) / float64(total) * 100
	l.WithFields(logrus.Fields{
		"processed":    processed,
		"total":        total,
		"percentage":   fmt.Sprintf("%.1f%%", percentage),
		"current_file": currentFile,
	}).Info("Processing progress")
}

// LogRetry logs retry attempts
func (l *Logger) LogRetry(filePath string, attempt int, maxAttempts int, err error) {
	l.WithFields(logrus.Fields{
		"file_path":    filePath,
		"attempt":      attempt,
		"max_attempts": maxAttempts,
		"error":        err.Error(),
	}).Warn("Retrying file processing")
}

// LogSkippedFile logs when a file is skipped
func (l *Logger) LogSkippedFile(filePath, reason string) {
	l.WithFields(logrus.Fields{
		"file_path": filePath,
		"reason":    reason,
	}).Warn("File skipped")
}
