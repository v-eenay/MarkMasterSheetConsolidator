// Package excel provides Excel file reading and writing operations for the Mark Master Sheet Consolidator.
// This file contains the writer functionality for updating master spreadsheets.
package excel

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/xuri/excelize/v2"
	"mark-master-sheet/internal/config"
	"mark-master-sheet/pkg/models"
)

// Writer handles writing to Excel files
type Writer struct {
	config *config.ExcelConfig
	reader *Reader
}

// NewWriter creates a new Excel writer
func NewWriter(cfg *config.ExcelConfig) *Writer {
	return &Writer{
		config: cfg,
		reader: NewReader(cfg),
	}
}

// CreateBackup creates a timestamped backup of the master sheet
func (w *Writer) CreateBackup(masterSheetPath, backupDir string) (string, error) {
	// Ensure backup directory exists
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create backup directory: %w", err)
	}

	// Generate backup filename with timestamp
	timestamp := time.Now().Format("20060102_150405")
	originalName := filepath.Base(masterSheetPath)
	ext := filepath.Ext(originalName)
	nameWithoutExt := originalName[:len(originalName)-len(ext)]
	backupName := fmt.Sprintf("%s_backup_%s%s", nameWithoutExt, timestamp, ext)
	backupPath := filepath.Join(backupDir, backupName)

	// Copy the file
	sourceFile, err := os.Open(masterSheetPath)
	if err != nil {
		return "", fmt.Errorf("failed to open source file: %w", err)
	}
	defer sourceFile.Close()

	destFile, err := os.Create(backupPath)
	if err != nil {
		return "", fmt.Errorf("failed to create backup file: %w", err)
	}
	defer destFile.Close()

	// Copy file contents
	buffer := make([]byte, 64*1024) // 64KB buffer
	for {
		n, err := sourceFile.Read(buffer)
		if n > 0 {
			if _, writeErr := destFile.Write(buffer[:n]); writeErr != nil {
				return "", fmt.Errorf("failed to write to backup file: %w", writeErr)
			}
		}
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return "", fmt.Errorf("failed to read from source file: %w", err)
		}
	}

	return backupPath, nil
}

// UpdateMasterSheet updates the master sheet with student data
func (w *Writer) UpdateMasterSheet(masterSheetPath string, studentData *models.StudentData) error {
	// Open the master sheet
	masterFile, err := excelize.OpenFile(masterSheetPath)
	if err != nil {
		return fmt.Errorf("failed to open master sheet: %w", err)
	}
	defer masterFile.Close()

	// Check if the master worksheet exists
	worksheets := masterFile.GetSheetList()
	worksheetExists := false
	for _, sheet := range worksheets {
		if sheet == w.config.MasterWorksheetName {
			worksheetExists = true
			break
		}
	}

	if !worksheetExists {
		return fmt.Errorf("master worksheet '%s' not found", w.config.MasterWorksheetName)
	}

	// Find the student in the master sheet
	rowNumber, err := w.reader.FindStudentInMasterSheet(masterFile, studentData.StudentID)
	if err != nil {
		return fmt.Errorf("student not found in master sheet: %w", err)
	}

	// Update marks in the corresponding columns
	for i, markCell := range w.config.MarkCells {
		if i >= len(w.config.MasterColumns) {
			break // Safety check
		}

		mark, exists := studentData.Marks[markCell]
		if !exists {
			continue // Skip if mark doesn't exist
		}

		// Skip empty marks (represented as -1)
		if mark < 0 {
			continue
		}

		// Calculate the target cell (column + row)
		targetCell := fmt.Sprintf("%s%d", w.config.MasterColumns[i], rowNumber)

		// Set the mark value
		if err := masterFile.SetCellFloat(w.config.MasterWorksheetName, targetCell, mark, 2, 64); err != nil {
			return fmt.Errorf("failed to set mark in cell %s: %w", targetCell, err)
		}
	}

	// Save the updated master sheet
	if err := masterFile.Save(); err != nil {
		return fmt.Errorf("failed to save master sheet: %w", err)
	}

	return nil
}

// SaveMasterSheetCopy saves a copy of the master sheet to the output directory
func (w *Writer) SaveMasterSheetCopy(masterSheetPath, outputDir string) (string, error) {
	// Ensure output directory exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create output directory: %w", err)
	}

	// Generate output filename with timestamp
	timestamp := time.Now().Format("20060102_150405")
	originalName := filepath.Base(masterSheetPath)
	ext := filepath.Ext(originalName)
	nameWithoutExt := originalName[:len(originalName)-len(ext)]
	outputName := fmt.Sprintf("%s_updated_%s%s", nameWithoutExt, timestamp, ext)
	outputPath := filepath.Join(outputDir, outputName)

	// Open the master sheet
	masterFile, err := excelize.OpenFile(masterSheetPath)
	if err != nil {
		return "", fmt.Errorf("failed to open master sheet: %w", err)
	}
	defer masterFile.Close()

	// Save as new file
	if err := masterFile.SaveAs(outputPath); err != nil {
		return "", fmt.Errorf("failed to save master sheet copy: %w", err)
	}

	return outputPath, nil
}

// BatchUpdateMasterSheet updates the master sheet with multiple student data entries
func (w *Writer) BatchUpdateMasterSheet(masterSheetPath string, studentDataList []*models.StudentData) (*models.ProcessingSummary, error) {
	summary := &models.ProcessingSummary{
		StartTime: time.Now(),
	}

	// Open the master sheet once for all updates
	masterFile, err := excelize.OpenFile(masterSheetPath)
	if err != nil {
		return summary, fmt.Errorf("failed to open master sheet: %w", err)
	}
	defer masterFile.Close()

	// Check if the master worksheet exists
	worksheets := masterFile.GetSheetList()
	worksheetExists := false
	for _, sheet := range worksheets {
		if sheet == w.config.MasterWorksheetName {
			worksheetExists = true
			break
		}
	}

	if !worksheetExists {
		return summary, fmt.Errorf("master worksheet '%s' not found", w.config.MasterWorksheetName)
	}

	// Process each student data
	for _, studentData := range studentDataList {
		// Find the student in the master sheet
		rowNumber, err := w.reader.FindStudentInMasterSheet(masterFile, studentData.StudentID)
		if err != nil {
			summary.StudentsNotFound++
			summary.Warnings = append(summary.Warnings,
				fmt.Sprintf("Student %s not found in master sheet", studentData.StudentID))
			continue
		}

		// Update marks in the corresponding columns
		markCount := 0
		for i, markCell := range w.config.MarkCells {
			if i >= len(w.config.MasterColumns) {
				break // Safety check
			}

			mark, exists := studentData.Marks[markCell]
			if !exists || mark < 0 {
				continue // Skip if mark doesn't exist or is empty
			}

			// Calculate the target cell (column + row)
			targetCell := fmt.Sprintf("%s%d", w.config.MasterColumns[i], rowNumber)

			// Set the mark value
			if err := masterFile.SetCellFloat(w.config.MasterWorksheetName, targetCell, mark, 2, 64); err != nil {
				summary.Errors = append(summary.Errors,
					fmt.Sprintf("Failed to set mark for student %s in cell %s: %v",
						studentData.StudentID, targetCell, err))
				continue
			}
			markCount++
		}

		if markCount > 0 {
			summary.StudentsUpdated++
		}
	}

	// Save the updated master sheet
	if err := masterFile.Save(); err != nil {
		return summary, fmt.Errorf("failed to save master sheet: %w", err)
	}

	summary.EndTime = time.Now()
	summary.TotalDuration = summary.EndTime.Sub(summary.StartTime)

	return summary, nil
}

// ValidateMasterSheet checks if the master sheet has the expected structure
func (w *Writer) ValidateMasterSheet(masterSheetPath string) error {
	masterFile, err := excelize.OpenFile(masterSheetPath)
	if err != nil {
		return fmt.Errorf("failed to open master sheet: %w", err)
	}
	defer masterFile.Close()

	// Check if the master worksheet exists
	worksheets := masterFile.GetSheetList()
	worksheetExists := false
	for _, sheet := range worksheets {
		if sheet == w.config.MasterWorksheetName {
			worksheetExists = true
			break
		}
	}

	if !worksheetExists {
		return fmt.Errorf("master worksheet '%s' not found", w.config.MasterWorksheetName)
	}

	// Check if there are any rows (at least header row)
	rows, err := masterFile.GetRows(w.config.MasterWorksheetName)
	if err != nil {
		return fmt.Errorf("failed to read master sheet rows: %w", err)
	}

	if len(rows) < 2 {
		return fmt.Errorf("master sheet appears to be empty or has no data rows")
	}

	return nil
}
