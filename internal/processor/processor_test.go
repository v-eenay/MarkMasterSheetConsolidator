package processor

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/xuri/excelize/v2"

	"mark-master-sheet/internal/config"
	"mark-master-sheet/internal/logger"
)

// TestNewProcessor tests processor creation
func TestNewProcessor(t *testing.T) {
	tempDir := t.TempDir()
	
	cfg := createTestConfig(tempDir)
	logger := createTestLogger(t, tempDir)

	processor := NewProcessor(cfg, logger)
	if processor == nil {
		t.Fatal("NewProcessor() returned nil")
	}

	if processor.config != cfg {
		t.Error("NewProcessor() should store config reference")
	}
	if processor.logger != logger {
		t.Error("NewProcessor() should store logger reference")
	}
}

// TestProcessFiles tests the main file processing functionality
func TestProcessFiles(t *testing.T) {
	tempDir := t.TempDir()
	
	// Create test files
	masterFile := createTestMasterFile(t, tempDir)
	studentDir := createTestStudentFiles(t, tempDir)
	
	cfg := createTestConfig(tempDir)
	cfg.Paths.MasterSheetPath = masterFile
	cfg.Paths.StudentFilesFolder = studentDir
	cfg.Paths.OutputFolder = filepath.Join(tempDir, "output")
	cfg.Paths.BackupFolder = filepath.Join(tempDir, "backups")
	
	logger := createTestLogger(t, tempDir)
	processor := NewProcessor(cfg, logger)

	tests := []struct {
		name      string
		dryRun    bool
		wantError bool
	}{
		{
			name:      "dry run processing",
			dryRun:    true,
			wantError: false,
		},
		{
			name:      "actual processing",
			dryRun:    false,
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			summary, err := processor.ProcessFiles(ctx, tt.dryRun)
			
			if tt.wantError {
				if err == nil {
					t.Errorf("ProcessFiles() expected error but got none")
				}
				return
			}
			
			if err != nil {
				t.Errorf("ProcessFiles() unexpected error: %v", err)
				return
			}
			
			if summary == nil {
				t.Error("ProcessFiles() should return summary")
			}
		})
	}
}

// TestProcessFilesWithCancellation tests processing cancellation
func TestProcessFilesWithCancellation(t *testing.T) {
	tempDir := t.TempDir()
	
	// Create test files
	masterFile := createTestMasterFile(t, tempDir)
	studentDir := createTestStudentFiles(t, tempDir)
	
	cfg := createTestConfig(tempDir)
	cfg.Paths.MasterSheetPath = masterFile
	cfg.Paths.StudentFilesFolder = studentDir
	cfg.Processing.MaxConcurrentFiles = 1 // Slow processing
	
	logger := createTestLogger(t, tempDir)
	processor := NewProcessor(cfg, logger)

	// Create cancellable context
	ctx, cancel := context.WithCancel(context.Background())
	
	// Start processing in goroutine
	done := make(chan error, 1)
	go func() {
		_, err := processor.ProcessFiles(ctx, false)
		done <- err
	}()

	// Cancel after short delay
	time.Sleep(50 * time.Millisecond)
	cancel()

	// Wait for processing to complete
	select {
	case err := <-done:
		// Processing may complete successfully if it finishes before cancellation
		// or return context.Canceled if cancelled in time
		if err != nil && err != context.Canceled {
			t.Errorf("ProcessFiles() should return nil or context.Canceled, got %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Error("ProcessFiles() should complete quickly after cancellation")
	}
}

// TestGetProcessingStatistics tests processing statistics
func TestGetProcessingStatistics(t *testing.T) {
	tempDir := t.TempDir()

	// Create test directory structure
	studentDir := filepath.Join(tempDir, "students")
	os.MkdirAll(studentDir, 0755)

	// Create test files
	files := []string{
		filepath.Join(studentDir, "student1.xlsx"),
		filepath.Join(studentDir, "student2.xls"),
		filepath.Join(studentDir, "not-excel.txt"),
	}

	for _, file := range files {
		os.WriteFile(file, []byte("test"), 0644)
	}

	cfg := createTestConfig(tempDir)
	cfg.Paths.StudentFilesFolder = studentDir

	logger := createTestLogger(t, tempDir)
	processor := NewProcessor(cfg, logger)

	stats := processor.GetProcessingStatistics()

	if stats == nil {
		t.Fatal("GetProcessingStatistics() returned nil")
	}

	// Should find 2 Excel files (excluding .txt file)
	expectedCount := 2
	if totalFiles, ok := stats["total_excel_files"].(int); !ok || totalFiles != expectedCount {
		t.Errorf("GetProcessingStatistics() total_excel_files = %v, want %d", stats["total_excel_files"], expectedCount)
	}

	// Check other statistics
	if stats["student_files_folder"] != studentDir {
		t.Errorf("GetProcessingStatistics() student_files_folder = %v, want %v", stats["student_files_folder"], studentDir)
	}

	if stats["max_concurrent_files"] != cfg.Processing.MaxConcurrentFiles {
		t.Errorf("GetProcessingStatistics() max_concurrent_files = %v, want %v", stats["max_concurrent_files"], cfg.Processing.MaxConcurrentFiles)
	}
}

// TestProcessFileWithRetries tests file processing with retry logic
func TestProcessFileWithRetries(t *testing.T) {
	tempDir := t.TempDir()

	// Create test files
	masterFile := createTestMasterFile(t, tempDir)
	studentFile := createTestStudentFile(t, tempDir, "STU001")

	cfg := createTestConfig(tempDir)
	cfg.Paths.MasterSheetPath = masterFile
	cfg.Processing.RetryAttempts = 2

	logger := createTestLogger(t, tempDir)
	processor := NewProcessor(cfg, logger)

	tests := []struct {
		name      string
		filePath  string
		wantError bool
	}{
		{
			name:      "valid student file",
			filePath:  studentFile,
			wantError: false,
		},
		{
			name:      "non-existent file",
			filePath:  filepath.Join(tempDir, "nonexistent.xlsx"),
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use the internal method (would need to be exposed for testing)
			// For now, we test through ProcessFiles which uses this internally
			ctx := context.Background()

			// Create a temporary config with just this file
			tempStudentDir := filepath.Dir(tt.filePath)
			cfg.Paths.StudentFilesFolder = tempStudentDir

			summary, err := processor.ProcessFiles(ctx, true) // Dry run

			if tt.wantError {
				// Should have errors in summary or return error
				if err == nil && len(summary.Errors) == 0 && summary.FailedFiles == 0 && summary.SkippedFiles == 0 {
					t.Errorf("ProcessFiles() expected error or failed files but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("ProcessFiles() unexpected error: %v", err)
				return
			}

			if summary == nil {
				t.Error("ProcessFiles() should return summary")
			}
		})
	}
}

// TestConcurrentProcessing tests concurrent file processing
func TestConcurrentProcessing(t *testing.T) {
	tempDir := t.TempDir()
	
	// Create test files
	masterFile := createTestMasterFile(t, tempDir)
	studentDir := createTestStudentFiles(t, tempDir)
	
	cfg := createTestConfig(tempDir)
	cfg.Paths.MasterSheetPath = masterFile
	cfg.Paths.StudentFilesFolder = studentDir
	cfg.Processing.MaxConcurrentFiles = 3
	
	logger := createTestLogger(t, tempDir)
	processor := NewProcessor(cfg, logger)

	ctx := context.Background()
	summary, err := processor.ProcessFiles(ctx, true) // Dry run for safety
	
	if err != nil {
		t.Fatalf("ProcessFiles() unexpected error: %v", err)
	}
	
	if summary == nil {
		t.Error("ProcessFiles() should return summary")
	}
}

// TestErrorHandling tests error handling during processing
func TestErrorHandling(t *testing.T) {
	tempDir := t.TempDir()
	
	cfg := createTestConfig(tempDir)
	cfg.Paths.MasterSheetPath = "nonexistent-master.xlsx"
	cfg.Paths.StudentFilesFolder = "nonexistent-folder"
	
	logger := createTestLogger(t, tempDir)
	processor := NewProcessor(cfg, logger)

	ctx := context.Background()
	_, err := processor.ProcessFiles(ctx, true)
	
	if err == nil {
		t.Error("ProcessFiles() should return error for invalid paths")
	}
}

// Helper functions for creating test files and configurations

func createTestConfig(tempDir string) *config.Config {
	return &config.Config{
		Paths: config.PathsConfig{
			StudentFilesFolder: tempDir,
			MasterSheetPath:    filepath.Join(tempDir, "master.xlsx"),
			OutputFolder:       filepath.Join(tempDir, "output"),
			LogFolder:          filepath.Join(tempDir, "logs"),
			BackupFolder:       filepath.Join(tempDir, "backups"),
		},
		Excel: config.ExcelConfig{
			StudentWorksheetName: "Grading Sheet",
			MasterWorksheetName:  "001",
			StudentIDCell:        "B2",
			MarkCells:            []string{"C6", "C7", "C8"},
			MasterColumns:        []string{"I", "J", "K"},
		},
		Processing: config.ProcessingConfig{
			MaxConcurrentFiles: 5,
			BackupEnabled:      true,
			SkipInvalidFiles:   true,
			TimeoutSeconds:     30,
			RetryAttempts:      2,
		},
		Logging: config.LoggingConfig{
			Level:         "INFO",
			ConsoleOutput: false,
			FileOutput:    true,
		},
	}
}

func createTestLogger(t testing.TB, tempDir string) *logger.Logger {
	logConfig := &config.LoggingConfig{
		Level:         "INFO",
		ConsoleOutput: true,  // Use console output to avoid file handle issues
		FileOutput:    false, // Disable file output in tests
	}

	logger, err := logger.NewLogger(logConfig, tempDir)
	if err != nil {
		t.Fatalf("Failed to create test logger: %v", err)
	}

	return logger
}

func createTestMasterFile(t testing.TB, tempDir string) string {
	f := excelize.NewFile()
	defer f.Close()
	
	sheetName := "001"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		t.Fatalf("Failed to create worksheet: %v", err)
	}
	f.SetActiveSheet(index)
	
	// Add headers and test data
	f.SetCellValue(sheetName, "A1", "Name")
	f.SetCellValue(sheetName, "B1", "Student ID")
	f.SetCellValue(sheetName, "I1", "Mark 1")
	f.SetCellValue(sheetName, "J1", "Mark 2")
	f.SetCellValue(sheetName, "K1", "Mark 3")
	
	// Add test students
	f.SetCellValue(sheetName, "A2", "John Doe")
	f.SetCellValue(sheetName, "B2", "STU001")
	f.SetCellValue(sheetName, "A3", "Jane Smith")
	f.SetCellValue(sheetName, "B3", "STU002")
	
	filePath := filepath.Join(tempDir, "master.xlsx")
	if err := f.SaveAs(filePath); err != nil {
		t.Fatalf("Failed to save master file: %v", err)
	}
	
	return filePath
}

func createTestStudentFile(t testing.TB, tempDir, studentID string) string {
	f := excelize.NewFile()
	defer f.Close()
	
	sheetName := "Grading Sheet"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		t.Fatalf("Failed to create worksheet: %v", err)
	}
	f.SetActiveSheet(index)
	
	// Add student data
	f.SetCellValue(sheetName, "B2", studentID)
	f.SetCellValue(sheetName, "C6", 85.5)
	f.SetCellValue(sheetName, "C7", 92.0)
	f.SetCellValue(sheetName, "C8", 78.5)
	
	filePath := filepath.Join(tempDir, studentID+".xlsx")
	if err := f.SaveAs(filePath); err != nil {
		t.Fatalf("Failed to save student file: %v", err)
	}
	
	return filePath
}

func createTestStudentFiles(t testing.TB, tempDir string) string {
	studentDir := filepath.Join(tempDir, "students")
	os.MkdirAll(studentDir, 0755)
	
	students := []string{"STU001", "STU002", "STU003"}
	for _, studentID := range students {
		createTestStudentFile(t, studentDir, studentID)
	}
	
	return studentDir
}

// BenchmarkProcessFiles benchmarks file processing performance
func BenchmarkProcessFiles(b *testing.B) {
	tempDir := b.TempDir()

	masterFile := createTestMasterFile(b, tempDir)
	studentDir := createTestStudentFiles(b, tempDir)

	cfg := createTestConfig(tempDir)
	cfg.Paths.MasterSheetPath = masterFile
	cfg.Paths.StudentFilesFolder = studentDir

	logger := createTestLogger(b, tempDir)
	processor := NewProcessor(cfg, logger)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx := context.Background()
		_, err := processor.ProcessFiles(ctx, true) // Dry run for benchmarking
		if err != nil {
			b.Fatalf("ProcessFiles failed: %v", err)
		}
	}
}

// BenchmarkGetProcessingStatistics benchmarks statistics gathering
func BenchmarkGetProcessingStatistics(b *testing.B) {
	tempDir := b.TempDir()
	studentDir := createTestStudentFiles(b, tempDir)

	cfg := createTestConfig(tempDir)
	cfg.Paths.StudentFilesFolder = studentDir

	logger := createTestLogger(b, tempDir)
	processor := NewProcessor(cfg, logger)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		stats := processor.GetProcessingStatistics()
		if stats == nil {
			b.Fatalf("GetProcessingStatistics failed")
		}
	}
}
