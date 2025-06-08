// Package processor provides the main processing logic for the Mark Master Sheet Consolidator.
// It orchestrates file discovery, concurrent processing, and master sheet updates.
package processor

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/schollz/progressbar/v3"
	"mark-master-sheet/internal/config"
	"mark-master-sheet/internal/excel"
	"mark-master-sheet/internal/logger"
	"mark-master-sheet/pkg/models"
)

// Processor handles the main processing logic
type Processor struct {
	config *config.Config
	logger *logger.Logger
	reader *excel.Reader
	writer *excel.Writer
}

// NewProcessor creates a new processor instance
func NewProcessor(cfg *config.Config, log *logger.Logger) *Processor {
	return &Processor{
		config: cfg,
		logger: log,
		reader: excel.NewReader(&cfg.Excel),
		writer: excel.NewWriter(&cfg.Excel),
	}
}

// ProcessFiles processes all Excel files in the student files directory
func (p *Processor) ProcessFiles(ctx context.Context, dryRun bool) (*models.ProcessingSummary, error) {
	summary := &models.ProcessingSummary{
		StartTime: time.Now(),
	}

	// Validate master sheet first
	if err := p.writer.ValidateMasterSheet(p.config.Paths.MasterSheetPath); err != nil {
		return summary, fmt.Errorf("master sheet validation failed: %w", err)
	}

	// Find all Excel files
	excelFiles, err := p.findExcelFiles(p.config.Paths.StudentFilesFolder)
	if err != nil {
		return summary, fmt.Errorf("failed to find Excel files: %w", err)
	}

	summary.TotalFiles = len(excelFiles)
	p.logger.LogProcessingStart(summary.TotalFiles)

	if summary.TotalFiles == 0 {
		p.logger.Info("No Excel files found to process")
		return summary, nil
	}

	// Create backup if enabled and not in dry run mode
	var backupPath string
	if p.config.Processing.BackupEnabled && !dryRun {
		backupPath, err = p.writer.CreateBackup(
			p.config.Paths.MasterSheetPath,
			p.config.Paths.BackupFolder,
		)
		if err != nil {
			return summary, fmt.Errorf("failed to create backup: %w", err)
		}
		p.logger.LogBackupCreated(p.config.Paths.MasterSheetPath, backupPath)
	}

	// Process files concurrently
	studentDataList, processingSummary := p.processFilesConcurrently(ctx, excelFiles)
	
	// Merge processing summary
	summary.SuccessfulFiles = processingSummary.SuccessfulFiles
	summary.FailedFiles = processingSummary.FailedFiles
	summary.SkippedFiles = processingSummary.SkippedFiles
	summary.Errors = processingSummary.Errors
	summary.Warnings = processingSummary.Warnings

	// Update master sheet if not in dry run mode
	if !dryRun && len(studentDataList) > 0 {
		updateSummary, err := p.writer.BatchUpdateMasterSheet(
			p.config.Paths.MasterSheetPath,
			studentDataList,
		)
		if err != nil {
			return summary, fmt.Errorf("failed to update master sheet: %w", err)
		}

		summary.StudentsUpdated = updateSummary.StudentsUpdated
		summary.StudentsNotFound = updateSummary.StudentsNotFound
		summary.Errors = append(summary.Errors, updateSummary.Errors...)
		summary.Warnings = append(summary.Warnings, updateSummary.Warnings...)

		// Save updated master sheet to output directory
		outputPath, err := p.writer.SaveMasterSheetCopy(
			p.config.Paths.MasterSheetPath,
			p.config.Paths.OutputFolder,
		)
		if err != nil {
			p.logger.Error("Failed to save master sheet copy: ", err)
		} else {
			p.logger.Info("Updated master sheet saved to: ", outputPath)
		}
	}

	summary.EndTime = time.Now()
	summary.TotalDuration = summary.EndTime.Sub(summary.StartTime)

	p.logger.LogProcessingEnd(summary)
	return summary, nil
}

// findExcelFiles recursively finds all Excel files in the given directory
func (p *Processor) findExcelFiles(rootDir string) ([]string, error) {
	var excelFiles []string

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			p.logger.LogFileError(path, err, "directory_walk")
			return nil // Continue walking despite errors
		}

		if info.IsDir() {
			return nil
		}

		// Check if it's an Excel file
		ext := strings.ToLower(filepath.Ext(path))
		if ext == ".xlsx" || ext == ".xls" {
			excelFiles = append(excelFiles, path)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking directory %s: %w", rootDir, err)
	}

	return excelFiles, nil
}

// processFilesConcurrently processes files using goroutines with rate limiting
func (p *Processor) processFilesConcurrently(ctx context.Context, files []string) ([]*models.StudentData, *models.ProcessingSummary) {
	summary := &models.ProcessingSummary{}
	var studentDataList []*models.StudentData
	var mu sync.Mutex

	// Create progress bar
	bar := progressbar.NewOptions(len(files),
		progressbar.OptionSetDescription("Processing files..."),
		progressbar.OptionShowCount(),
		progressbar.OptionShowIts(),
		progressbar.OptionSetPredictTime(true),
	)

	// Create semaphore for rate limiting
	semaphore := make(chan struct{}, p.config.Processing.MaxConcurrentFiles)
	var wg sync.WaitGroup

	// Process each file
	for _, filePath := range files {
		select {
		case <-ctx.Done():
			p.logger.Warn("Processing cancelled by context")
			break
		default:
		}

		wg.Add(1)
		go func(path string) {
			defer wg.Done()
			defer bar.Add(1)

			// Acquire semaphore
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			// Process file with retries
			result := p.processFileWithRetries(path)

			// Update summary and collect data
			mu.Lock()
			if result.Success {
				summary.SuccessfulFiles++
				if result.StudentData != nil {
					studentDataList = append(studentDataList, result.StudentData)
				}
			} else {
				if p.config.Processing.SkipInvalidFiles {
					summary.SkippedFiles++
					p.logger.LogSkippedFile(path, result.Error.Error())
				} else {
					summary.FailedFiles++
					summary.Errors = append(summary.Errors, 
						fmt.Sprintf("File %s: %v", path, result.Error))
				}
			}
			mu.Unlock()

			// Log progress
			mu.Lock()
			processed := summary.SuccessfulFiles + summary.FailedFiles + summary.SkippedFiles
			mu.Unlock()
			
			if processed%10 == 0 { // Log every 10 files
				p.logger.LogProgress(processed, len(files), path)
			}
		}(filePath)
	}

	wg.Wait()
	bar.Finish()

	return studentDataList, summary
}

// processFileWithRetries processes a single file with retry logic
func (p *Processor) processFileWithRetries(filePath string) *models.ProcessingResult {
	result := &models.ProcessingResult{
		FilePath: filePath,
	}

	startTime := time.Now()
	defer func() {
		result.Duration = time.Since(startTime)
	}()

	var lastErr error
	for attempt := 1; attempt <= p.config.Processing.RetryAttempts; attempt++ {
		studentData, err := p.reader.ReadStudentData(filePath)
		if err == nil {
			result.Success = true
			result.StudentData = studentData
			
			p.logger.LogFileProcessed(
				filePath,
				studentData.StudentID,
				studentData.GetMarkCount(),
				result.Duration,
			)
			return result
		}

		lastErr = err
		if attempt < p.config.Processing.RetryAttempts {
			p.logger.LogRetry(filePath, attempt, p.config.Processing.RetryAttempts, err)
			time.Sleep(time.Duration(attempt) * time.Second) // Exponential backoff
		}
	}

	result.Success = false
	result.Error = lastErr
	p.logger.LogFileError(filePath, lastErr, "processing")

	return result
}

// GetProcessingStatistics returns current processing statistics
func (p *Processor) GetProcessingStatistics() map[string]interface{} {
	stats := make(map[string]interface{})
	
	// Count total files
	excelFiles, err := p.findExcelFiles(p.config.Paths.StudentFilesFolder)
	if err != nil {
		stats["error"] = err.Error()
		return stats
	}

	stats["total_excel_files"] = len(excelFiles)
	stats["student_files_folder"] = p.config.Paths.StudentFilesFolder
	stats["master_sheet_path"] = p.config.Paths.MasterSheetPath
	stats["max_concurrent_files"] = p.config.Processing.MaxConcurrentFiles
	stats["backup_enabled"] = p.config.Processing.BackupEnabled

	return stats
}
