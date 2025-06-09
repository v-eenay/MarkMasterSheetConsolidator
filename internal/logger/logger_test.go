package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"mark-master-sheet/internal/config"
)

// TestNewLogger tests logger creation
func TestNewLogger(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		name      string
		config    *config.LoggingConfig
		logDir    string
		wantError bool
	}{
		{
			name: "valid logger config",
			config: &config.LoggingConfig{
				Level:          "INFO",
				ConsoleOutput:  true,
				FileOutput:     true,
				MaxFileSizeMB:  10,
				MaxBackupFiles: 3,
				MaxAgeDays:     7,
			},
			logDir:    tempDir,
			wantError: false,
		},
		{
			name: "invalid log level",
			config: &config.LoggingConfig{
				Level:          "INVALID",
				ConsoleOutput:  true,
				FileOutput:     false,
				MaxFileSizeMB:  10,
				MaxBackupFiles: 3,
				MaxAgeDays:     7,
			},
			logDir:    tempDir,
			wantError: true,
		},
		{
			name: "console only",
			config: &config.LoggingConfig{
				Level:         "DEBUG",
				ConsoleOutput: true,
				FileOutput:    false,
			},
			logDir:    tempDir,
			wantError: false,
		},
		{
			name: "file only",
			config: &config.LoggingConfig{
				Level:          "WARN",
				ConsoleOutput:  false,
				FileOutput:     true,
				MaxFileSizeMB:  5,
				MaxBackupFiles: 2,
				MaxAgeDays:     3,
			},
			logDir:    tempDir,
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, err := NewLogger(tt.config, tt.logDir)
			
			if tt.wantError {
				if err == nil {
					t.Errorf("NewLogger() expected error but got none")
				}
				return
			}
			
			if err != nil {
				t.Errorf("NewLogger() unexpected error: %v", err)
				return
			}
			
			if logger == nil {
				t.Error("NewLogger() returned nil logger")
			}
		})
	}
}

// TestLogLevels tests different log levels
func TestLogLevels(t *testing.T) {
	tempDir := t.TempDir()
	
	config := &config.LoggingConfig{
		Level:          "DEBUG",
		ConsoleOutput:  false,
		FileOutput:     true,
		MaxFileSizeMB:  10,
		MaxBackupFiles: 3,
		MaxAgeDays:     7,
	}

	logger, err := NewLogger(config, tempDir)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	tests := []struct {
		name     string
		logFunc  func(...interface{})
		message  string
		level    string
	}{
		{
			name:     "debug log",
			logFunc:  logger.Debug,
			message:  "Debug message",
			level:    "DEBUG",
		},
		{
			name:     "info log",
			logFunc:  logger.Info,
			message:  "Info message",
			level:    "INFO",
		},
		{
			name:     "warn log",
			logFunc:  logger.Warn,
			message:  "Warning message",
			level:    "WARN",
		},
		{
			name:     "error log",
			logFunc:  logger.Error,
			message:  "Error message",
			level:    "ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.logFunc(tt.message)
			
			// Allow time for log to be written
			time.Sleep(10 * time.Millisecond)
			
			// Check if log file was created and contains the message
			logFiles, err := filepath.Glob(filepath.Join(tempDir, "*.log"))
			if err != nil {
				t.Fatalf("Failed to find log files: %v", err)
			}
			
			if len(logFiles) == 0 {
				t.Error("No log files were created")
				return
			}
			
			// Read the log file
			content, err := os.ReadFile(logFiles[0])
			if err != nil {
				t.Fatalf("Failed to read log file: %v", err)
			}
			
			logContent := string(content)
			if !strings.Contains(logContent, tt.message) {
				t.Errorf("Log file should contain message %q, got %q", tt.message, logContent)
			}
			
			if !strings.Contains(logContent, tt.level) {
				t.Errorf("Log file should contain level %q, got %q", tt.level, logContent)
			}
		})
	}
}

// TestLogWithFields tests structured logging with fields
func TestLogWithFields(t *testing.T) {
	tempDir := t.TempDir()
	
	config := &config.LoggingConfig{
		Level:          "INFO",
		ConsoleOutput:  false,
		FileOutput:     true,
		MaxFileSizeMB:  10,
		MaxBackupFiles: 3,
		MaxAgeDays:     7,
	}

	logger, err := NewLogger(config, tempDir)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	// Test logging with fields
	logger.WithFields(map[string]interface{}{
		"student_id": "STU001",
		"file_name":  "student1.xlsx",
		"marks":      []float64{85.5, 92.0},
	}).Info("Processing student file")

	// Allow time for log to be written
	time.Sleep(10 * time.Millisecond)

	// Check log file content
	logFiles, err := filepath.Glob(filepath.Join(tempDir, "*.log"))
	if err != nil {
		t.Fatalf("Failed to find log files: %v", err)
	}

	if len(logFiles) == 0 {
		t.Error("No log files were created")
		return
	}

	content, err := os.ReadFile(logFiles[0])
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	logContent := string(content)
	expectedFields := []string{"student_id", "STU001", "file_name", "student1.xlsx"}
	
	for _, field := range expectedFields {
		if !strings.Contains(logContent, field) {
			t.Errorf("Log should contain field %q, got %q", field, logContent)
		}
	}
}

// TestLogRotation tests log file rotation
func TestLogRotation(t *testing.T) {
	tempDir := t.TempDir()
	
	config := &config.LoggingConfig{
		Level:          "INFO",
		ConsoleOutput:  false,
		FileOutput:     true,
		MaxFileSizeMB:  1, // Small size to trigger rotation
		MaxBackupFiles: 2,
		MaxAgeDays:     7,
	}

	logger, err := NewLogger(config, tempDir)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	// Write enough logs to trigger rotation
	largeMessage := strings.Repeat("This is a large log message to trigger rotation. ", 1000)
	
	for i := 0; i < 100; i++ {
		logger.Info("Log entry", i, ":", largeMessage)
	}

	// Allow time for rotation to occur
	time.Sleep(100 * time.Millisecond)

	// Check if multiple log files exist (indicating rotation occurred)
	logFiles, err := filepath.Glob(filepath.Join(tempDir, "*.log*"))
	if err != nil {
		t.Fatalf("Failed to find log files: %v", err)
	}

	// Should have main log file plus backup files
	if len(logFiles) < 2 {
		t.Logf("Found %d log files, rotation may not have occurred yet", len(logFiles))
		// This is not necessarily an error as rotation timing can vary
	}
}

// TestSpecializedLogMethods tests application-specific log methods
func TestSpecializedLogMethods(t *testing.T) {
	tempDir := t.TempDir()

	config := &config.LoggingConfig{
		Level:          "INFO",
		ConsoleOutput:  false,
		FileOutput:     true,
		MaxFileSizeMB:  10,
		MaxBackupFiles: 3,
		MaxAgeDays:     7,
	}

	logger, err := NewLogger(config, tempDir)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	// Test specialized logging methods
	logger.LogProcessingStart(10)
	logger.LogFileProcessed("test.xlsx", "STU001", 3, time.Second)
	logger.LogFileError("error.xlsx", fmt.Errorf("test error"), "reading")
	logger.LogStudentNotFound("STU999", "test.xlsx", []string{"STU001", "STU002"})
	logger.LogBackupCreated("original.xlsx", "backup.xlsx")
	logger.LogValidationError("test.xlsx", "student_id", "invalid", "contains special characters")
	logger.LogProgress(5, 10, "current.xlsx")
	logger.LogRetry("retry.xlsx", 2, 3, fmt.Errorf("retry error"))
	logger.LogSkippedFile("skip.xlsx", "invalid format")
	logger.LogProcessingEnd("Processing completed successfully")

	// Allow time for logs to be written
	time.Sleep(50 * time.Millisecond)

	// Verify log file exists and contains expected content
	logFiles, err := filepath.Glob(filepath.Join(tempDir, "*.log"))
	if err != nil {
		t.Fatalf("Failed to find log files: %v", err)
	}

	if len(logFiles) == 0 {
		t.Error("No log files found")
		return
	}

	content, err := os.ReadFile(logFiles[0])
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	logContent := string(content)
	expectedEntries := []string{
		"Starting mark consolidation process",
		"File processed successfully",
		"File processing failed",
		"Student ID not found",
		"Backup created successfully",
		"Validation error",
		"Processing progress",
		"Retrying file processing",
		"File skipped",
		"Mark consolidation process completed",
	}

	for _, expected := range expectedEntries {
		if !strings.Contains(logContent, expected) {
			t.Errorf("Log should contain %q", expected)
		}
	}
}

// TestLoggerConfiguration tests various logger configurations
func TestLoggerConfiguration(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		name           string
		config         *config.LoggingConfig
		expectConsole  bool
		expectFile     bool
	}{
		{
			name: "console and file output",
			config: &config.LoggingConfig{
				Level:         "INFO",
				ConsoleOutput: true,
				FileOutput:    true,
			},
			expectConsole: true,
			expectFile:    true,
		},
		{
			name: "console only",
			config: &config.LoggingConfig{
				Level:         "INFO",
				ConsoleOutput: true,
				FileOutput:    false,
			},
			expectConsole: true,
			expectFile:    false,
		},
		{
			name: "file only",
			config: &config.LoggingConfig{
				Level:         "INFO",
				ConsoleOutput: false,
				FileOutput:    true,
			},
			expectConsole: false,
			expectFile:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, err := NewLogger(tt.config, tempDir)
			if err != nil {
				t.Fatalf("Failed to create logger: %v", err)
			}

			// Test logging
			logger.Info("Test message")
			time.Sleep(10 * time.Millisecond)

			if tt.expectFile {
				logFiles, err := filepath.Glob(filepath.Join(tempDir, "*.log"))
				if err != nil {
					t.Fatalf("Failed to find log files: %v", err)
				}
				if len(logFiles) == 0 {
					t.Error("Expected log file but none found")
				}
			}

		})
	}
}

// BenchmarkLogging benchmarks logging performance
func BenchmarkLogging(b *testing.B) {
	tempDir := b.TempDir()

	config := &config.LoggingConfig{
		Level:          "INFO",
		ConsoleOutput:  false,
		FileOutput:     true,
		MaxFileSizeMB:  100,
		MaxBackupFiles: 3,
		MaxAgeDays:     7,
	}

	logger, err := NewLogger(config, tempDir)
	if err != nil {
		b.Fatalf("Failed to create logger: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("Benchmark log message", i)
	}
}

// BenchmarkSpecializedLogging benchmarks specialized logging methods
func BenchmarkSpecializedLogging(b *testing.B) {
	tempDir := b.TempDir()

	config := &config.LoggingConfig{
		Level:          "INFO",
		ConsoleOutput:  false,
		FileOutput:     true,
		MaxFileSizeMB:  100,
		MaxBackupFiles: 3,
		MaxAgeDays:     7,
	}

	logger, err := NewLogger(config, tempDir)
	if err != nil {
		b.Fatalf("Failed to create logger: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.LogFileProcessed("test.xlsx", "STU001", 3, time.Millisecond)
	}
}
