package gui

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"fyne.io/fyne/v2/test"

	"mark-master-sheet/internal/config"
)

// TestValidatePaths tests path validation functionality
func TestValidatePaths(t *testing.T) {
	testApp := test.NewApp()
	defer testApp.Quit()

	app := NewApp()
	tempDir := t.TempDir()

	// Create test files
	masterFile := filepath.Join(tempDir, "master.xlsx")
	studentDir := filepath.Join(tempDir, "students")
	
	// Create master file
	if err := os.WriteFile(masterFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test master file: %v", err)
	}
	
	// Create student directory
	if err := os.MkdirAll(studentDir, 0755); err != nil {
		t.Fatalf("Failed to create test student directory: %v", err)
	}

	tests := []struct {
		name      string
		cfg       *config.Config
		wantError bool
	}{
		{
			name: "valid paths",
			cfg: &config.Config{
				Paths: config.PathsConfig{
					MasterSheetPath:    masterFile,
					StudentFilesFolder: studentDir,
				},
			},
			wantError: false,
		},
		{
			name: "non-existent master file",
			cfg: &config.Config{
				Paths: config.PathsConfig{
					MasterSheetPath:    filepath.Join(tempDir, "nonexistent.xlsx"),
					StudentFilesFolder: studentDir,
				},
			},
			wantError: true,
		},
		{
			name: "non-existent student folder",
			cfg: &config.Config{
				Paths: config.PathsConfig{
					MasterSheetPath:    masterFile,
					StudentFilesFolder: filepath.Join(tempDir, "nonexistent"),
				},
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := app.validatePaths(tt.cfg)
			
			if tt.wantError && err == nil {
				t.Errorf("validatePaths() expected error but got none")
			}
			if !tt.wantError && err != nil {
				t.Errorf("validatePaths() unexpected error: %v", err)
			}
		})
	}
}

// TestUpdateProcessingUI tests processing UI state updates
func TestUpdateProcessingUI(t *testing.T) {
	testApp := test.NewApp()
	defer testApp.Quit()

	app := NewApp()
	app.setupUI()

	tests := []struct {
		name       string
		processing bool
		dryRun     bool
	}{
		{
			name:       "start processing",
			processing: true,
			dryRun:     false,
		},
		{
			name:       "start dry run",
			processing: true,
			dryRun:     true,
		},
		{
			name:       "stop processing",
			processing: false,
			dryRun:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app.updateProcessingUI(tt.processing, tt.dryRun)
			
			if tt.processing {
				if !app.progressBar.Visible() {
					t.Error("Progress bar should be visible during processing")
				}
				
				statusText := app.statusLabel.Text
				if tt.dryRun {
					if !strings.Contains(statusText, "dry run") {
						t.Error("Status should indicate dry run mode")
					}
				} else {
					if !strings.Contains(statusText, "Processing") {
						t.Error("Status should indicate processing mode")
					}
				}
			} else {
				if app.progressBar.Visible() {
					t.Error("Progress bar should be hidden when not processing")
				}
				
				statusText := app.statusLabel.Text
				if !strings.Contains(statusText, "Ready") {
					t.Error("Status should indicate ready state")
				}
			}
		})
	}
}

// TestAppendLog tests log output functionality
func TestAppendLog(t *testing.T) {
	testApp := test.NewApp()
	defer testApp.Quit()

	app := NewApp()
	app.setupUI()

	// Test initial log state
	initialText := app.logOutput.Text
	if initialText == "" {
		t.Error("Log output should have initial text")
	}

	// Test appending log
	testMessage := "Test log message"
	app.appendLog(testMessage)

	updatedText := app.logOutput.Text
	if !strings.Contains(updatedText, testMessage) {
		t.Errorf("Log output should contain %q, got %q", testMessage, updatedText)
	}

	// Test multiple appends
	secondMessage := "Second log message"
	app.appendLog(secondMessage)

	finalText := app.logOutput.Text
	if !strings.Contains(finalText, testMessage) || !strings.Contains(finalText, secondMessage) {
		t.Error("Log output should contain both messages")
	}
}

// TestUpdateProgress tests progress update functionality
func TestUpdateProgress(t *testing.T) {
	testApp := test.NewApp()
	defer testApp.Quit()

	app := NewApp()
	app.setupUI()

	tests := []struct {
		name        string
		current     int
		total       int
		currentFile string
	}{
		{
			name:        "progress update with file",
			current:     5,
			total:       10,
			currentFile: "student1.xlsx",
		},
		{
			name:        "progress update without file",
			current:     8,
			total:       10,
			currentFile: "",
		},
		{
			name:        "zero total",
			current:     0,
			total:       0,
			currentFile: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app.updateProgress(tt.current, tt.total, tt.currentFile)
			
			if tt.total > 0 {
				expectedProgress := float64(tt.current) / float64(tt.total)
				actualProgress := app.progressBar.Value
				
				if actualProgress != expectedProgress {
					t.Errorf("Progress bar value = %v, want %v", actualProgress, expectedProgress)
				}
				
				statusText := app.statusLabel.Text
				if !strings.Contains(statusText, "Processing") {
					t.Error("Status should indicate processing")
				}
			}
		})
	}
}

// TestLogMethods tests different log level methods
func TestLogMethods(t *testing.T) {
	testApp := test.NewApp()
	defer testApp.Quit()

	app := NewApp()
	app.setupUI()

	tests := []struct {
		name     string
		logFunc  func(string)
		message  string
		expected string
	}{
		{
			name:     "log error",
			logFunc:  app.logError,
			message:  "Test error",
			expected: "ERROR: Test error",
		},
		{
			name:     "log warning",
			logFunc:  app.logWarning,
			message:  "Test warning",
			expected: "WARNING: Test warning",
		},
		{
			name:     "log info",
			logFunc:  app.logInfo,
			message:  "Test info",
			expected: "INFO: Test info",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initialLength := len(app.logOutput.Text)
			tt.logFunc(tt.message)
			
			updatedText := app.logOutput.Text
			if !strings.Contains(updatedText, tt.expected) {
				t.Errorf("Log should contain %q, got %q", tt.expected, updatedText)
			}
			
			if len(updatedText) <= initialLength {
				t.Error("Log text should have increased in length")
			}
		})
	}
}

// TestDisplayProcessingSummary tests processing summary display
func TestDisplayProcessingSummary(t *testing.T) {
	testApp := test.NewApp()
	defer testApp.Quit()

	app := NewApp()
	app.setupUI()

	tests := []struct {
		name     string
		summary  interface{}
		dryRun   bool
		duration time.Duration
	}{
		{
			name:     "dry run summary",
			summary:  "test summary",
			dryRun:   true,
			duration: 5 * time.Second,
		},
		{
			name:     "production summary",
			summary:  "test summary",
			dryRun:   false,
			duration: 10 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initialText := app.logOutput.Text
			app.displayProcessingSummary(tt.summary, tt.dryRun, tt.duration)
			
			updatedText := app.logOutput.Text
			if len(updatedText) <= len(initialText) {
				t.Error("Log should have been updated with summary")
			}
			
			if !strings.Contains(updatedText, "PROCESSING SUMMARY") {
				t.Error("Log should contain processing summary header")
			}
			
			if tt.dryRun {
				if !strings.Contains(updatedText, "DRY RUN") {
					t.Error("Summary should indicate dry run mode")
				}
			} else {
				if !strings.Contains(updatedText, "PRODUCTION") {
					t.Error("Summary should indicate production mode")
				}
			}
			
			if !strings.Contains(updatedText, "Duration:") {
				t.Error("Summary should contain duration information")
			}
		})
	}
}

// TestProcessingStateManagement tests processing state management
func TestProcessingStateManagement(t *testing.T) {
	testApp := test.NewApp()
	defer testApp.Quit()

	app := NewApp()
	app.setupUI()

	// Test initial state
	if app.isProcessing {
		t.Error("App should not be processing initially")
	}

	// Test state during processing simulation
	app.isProcessing = true
	if !app.isProcessing {
		t.Error("App should be in processing state")
	}

	// Test context creation
	ctx, cancel := context.WithCancel(context.Background())
	app.processingContext = ctx
	app.cancelProcessing = cancel

	if app.processingContext == nil {
		t.Error("Processing context should be set")
	}
	if app.cancelProcessing == nil {
		t.Error("Cancel function should be set")
	}

	// Test context cancellation
	app.cancelProcessing()
	
	select {
	case <-app.processingContext.Done():
		// Context was cancelled successfully
	case <-time.After(100 * time.Millisecond):
		t.Error("Context should have been cancelled")
	}
}

// TestStopProcessing tests processing cancellation
func TestStopProcessing(t *testing.T) {
	testApp := test.NewApp()
	defer testApp.Quit()

	app := NewApp()
	app.setupUI()

	// Test stopping when not processing
	app.stopProcessing()
	// Should not cause any errors

	// Test stopping during processing
	app.isProcessing = true
	ctx, cancel := context.WithCancel(context.Background())
	app.processingContext = ctx
	app.cancelProcessing = cancel

	app.stopProcessing()

	// Verify context was cancelled
	select {
	case <-ctx.Done():
		// Context was cancelled successfully
	case <-time.After(100 * time.Millisecond):
		t.Error("Context should have been cancelled")
	}
}

// BenchmarkAppendLog benchmarks log appending
func BenchmarkAppendLog(b *testing.B) {
	testApp := test.NewApp()
	defer testApp.Quit()

	app := NewApp()
	app.setupUI()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		app.appendLog("Benchmark log message")
	}
}

// BenchmarkUpdateProgress benchmarks progress updates
func BenchmarkUpdateProgress(b *testing.B) {
	testApp := test.NewApp()
	defer testApp.Quit()

	app := NewApp()
	app.setupUI()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		app.updateProgress(i, b.N, "test-file.xlsx")
	}
}

// BenchmarkUpdateProcessingUI benchmarks processing UI updates
func BenchmarkUpdateProcessingUI(b *testing.B) {
	testApp := test.NewApp()
	defer testApp.Quit()

	app := NewApp()
	app.setupUI()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		app.updateProcessingUI(i%2 == 0, i%4 == 0)
	}
}
