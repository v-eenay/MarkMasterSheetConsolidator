// Package excel provides Excel file reading and writing operations for the Mark Master Sheet Consolidator.
// It handles data extraction from student files and updates to master spreadsheets.
package excel

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
	"mark-master-sheet/internal/config"
	"mark-master-sheet/pkg/models"
)

// Reader handles reading Excel files
type Reader struct {
	config *config.ExcelConfig
}

// NewReader creates a new Excel reader
func NewReader(cfg *config.ExcelConfig) *Reader {
	return &Reader{
		config: cfg,
	}
}

// ReadStudentData reads student data from an Excel file
func (r *Reader) ReadStudentData(filePath string) (*models.StudentData, error) {
	// Check file extension
	ext := strings.ToLower(filepath.Ext(filePath))
	if ext != ".xlsx" && ext != ".xls" {
		return nil, &models.FileProcessingError{
			FilePath: filePath,
			Stage:    "validation",
			Message:  "unsupported file format",
		}
	}

	// Open the Excel file
	file, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, &models.FileProcessingError{
			FilePath: filePath,
			Stage:    "opening",
			Message:  "failed to open Excel file",
			Cause:    err,
		}
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			// Log the error but don't override the main error
		}
	}()

	// Check if the required worksheet exists
	worksheets := file.GetSheetList()
	worksheetExists := false
	for _, sheet := range worksheets {
		if sheet == r.config.StudentWorksheetName {
			worksheetExists = true
			break
		}
	}

	if !worksheetExists {
		return nil, &models.FileProcessingError{
			FilePath: filePath,
			Stage:    "worksheet_validation",
			Message:  fmt.Sprintf("worksheet '%s' not found", r.config.StudentWorksheetName),
		}
	}

	// Read student ID
	studentID, err := file.GetCellValue(r.config.StudentWorksheetName, r.config.StudentIDCell)
	if err != nil {
		return nil, &models.FileProcessingError{
			FilePath: filePath,
			Stage:    "student_id_reading",
			Message:  fmt.Sprintf("failed to read student ID from cell %s", r.config.StudentIDCell),
			Cause:    err,
		}
	}

	// Clean and validate student ID
	studentID = strings.TrimSpace(studentID)
	if studentID == "" {
		return nil, &models.ValidationError{
			Field:   "student_id",
			Value:   studentID,
			Message: "student ID is empty",
			File:    filePath,
		}
	}

	// Create student data structure
	studentData := &models.StudentData{
		StudentID: studentID,
		FilePath:  filePath,
		Marks:     make(map[string]float64),
		Timestamp: time.Now(),
	}

	// Validate student ID format
	if !studentData.IsValidStudentID() {
		return nil, &models.ValidationError{
			Field:   "student_id",
			Value:   studentID,
			Message: "student ID contains invalid characters (only alphanumeric allowed)",
			File:    filePath,
		}
	}

	// Read marks from specified cells
	for _, cell := range r.config.MarkCells {
		markValue, err := file.GetCellValue(r.config.StudentWorksheetName, cell)
		if err != nil {
			return nil, &models.FileProcessingError{
				FilePath: filePath,
				Stage:    "mark_reading",
				Message:  fmt.Sprintf("failed to read mark from cell %s", cell),
				Cause:    err,
			}
		}

		// Handle empty cells
		markValue = strings.TrimSpace(markValue)
		if markValue == "" {
			// Store as -1 to indicate empty/missing mark
			studentData.Marks[cell] = -1
			continue
		}

		// Parse numeric value
		mark, err := strconv.ParseFloat(markValue, 64)
		if err != nil {
			return nil, &models.ValidationError{
				Field:   fmt.Sprintf("mark_%s", cell),
				Value:   markValue,
				Message: "mark is not a valid number",
				File:    filePath,
			}
		}

		// Validate mark range (assuming 0-100 is valid range)
		if mark < 0 || mark > 100 {
			return nil, &models.ValidationError{
				Field:   fmt.Sprintf("mark_%s", cell),
				Value:   markValue,
				Message: "mark is outside valid range (0-100)",
				File:    filePath,
			}
		}

		studentData.Marks[cell] = mark
	}

	return studentData, nil
}

// FindStudentInMasterSheet finds a student ID in the master sheet and returns the row number
func (r *Reader) FindStudentInMasterSheet(masterFile *excelize.File, studentID string) (int, error) {
	// Get all rows from column B (student ID column)
	rows, err := masterFile.GetRows(r.config.MasterWorksheetName)
	if err != nil {
		return 0, fmt.Errorf("failed to read master sheet rows: %w", err)
	}

	// Search for student ID (case-insensitive)
	studentIDLower := strings.ToLower(strings.TrimSpace(studentID))

	for rowIndex, row := range rows {
		if len(row) > 1 { // Ensure column B exists
			cellValue := strings.ToLower(strings.TrimSpace(row[1])) // Column B is index 1
			if cellValue == studentIDLower {
				return rowIndex + 1, nil // Excel rows are 1-based
			}
		}
	}

	return 0, fmt.Errorf("student ID %s not found in master sheet", studentID)
}

// GetSimilarStudentIDs returns student IDs that are similar to the given ID
func (r *Reader) GetSimilarStudentIDs(masterFile *excelize.File, targetID string, maxSuggestions int) []string {
	rows, err := masterFile.GetRows(r.config.MasterWorksheetName)
	if err != nil {
		return nil
	}

	var suggestions []string
	targetIDLower := strings.ToLower(strings.TrimSpace(targetID))

	for _, row := range rows {
		if len(row) > 1 && len(suggestions) < maxSuggestions {
			cellValue := strings.TrimSpace(row[1])
			if cellValue != "" {
				cellValueLower := strings.ToLower(cellValue)

				// Simple similarity check: contains substring or similar length
				if strings.Contains(cellValueLower, targetIDLower) ||
					strings.Contains(targetIDLower, cellValueLower) ||
					levenshteinDistance(targetIDLower, cellValueLower) <= 2 {
					suggestions = append(suggestions, cellValue)
				}
			}
		}
	}

	return suggestions
}

// levenshteinDistance calculates the Levenshtein distance between two strings
func levenshteinDistance(s1, s2 string) int {
	if len(s1) == 0 {
		return len(s2)
	}
	if len(s2) == 0 {
		return len(s1)
	}

	matrix := make([][]int, len(s1)+1)
	for i := range matrix {
		matrix[i] = make([]int, len(s2)+1)
		matrix[i][0] = i
	}
	for j := range matrix[0] {
		matrix[0][j] = j
	}

	for i := 1; i <= len(s1); i++ {
		for j := 1; j <= len(s2); j++ {
			cost := 0
			if s1[i-1] != s2[j-1] {
				cost = 1
			}

			matrix[i][j] = min(
				matrix[i-1][j]+1,      // deletion
				matrix[i][j-1]+1,      // insertion
				matrix[i-1][j-1]+cost, // substitution
			)
		}
	}

	return matrix[len(s1)][len(s2)]
}

// min returns the minimum of three integers
func min(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}
