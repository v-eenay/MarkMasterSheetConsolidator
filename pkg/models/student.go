// Package models defines data structures and validation methods for the Mark Master Sheet Consolidator.
// It provides types for student data, processing results, and error handling.
package models

import (
	"fmt"
	"time"
)

// StudentData represents the extracted data from a student's Excel file
type StudentData struct {
	StudentID string             `json:"student_id"`
	FilePath  string             `json:"file_path"`
	Marks     map[string]float64 `json:"marks"`
	Timestamp time.Time          `json:"timestamp"`
}

// ProcessingResult represents the result of processing a single file
type ProcessingResult struct {
	StudentData *StudentData `json:"student_data,omitempty"`
	FilePath    string       `json:"file_path"`
	Success     bool         `json:"success"`
	Error       error        `json:"error,omitempty"`
	Duration    time.Duration `json:"duration"`
}

// ProcessingSummary contains overall processing statistics
type ProcessingSummary struct {
	TotalFiles       int           `json:"total_files"`
	SuccessfulFiles  int           `json:"successful_files"`
	FailedFiles      int           `json:"failed_files"`
	SkippedFiles     int           `json:"skipped_files"`
	StudentsUpdated  int           `json:"students_updated"`
	StudentsNotFound int           `json:"students_not_found"`
	TotalDuration    time.Duration `json:"total_duration"`
	StartTime        time.Time     `json:"start_time"`
	EndTime          time.Time     `json:"end_time"`
	Errors           []string      `json:"errors,omitempty"`
	Warnings         []string      `json:"warnings,omitempty"`
}

// ValidationError represents a validation error with context
type ValidationError struct {
	Field   string `json:"field"`
	Value   string `json:"value"`
	Message string `json:"message"`
	File    string `json:"file"`
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error in file %s, field %s (value: %s): %s", 
		e.File, e.Field, e.Value, e.Message)
}

// FileProcessingError represents an error during file processing
type FileProcessingError struct {
	FilePath string `json:"file_path"`
	Stage    string `json:"stage"`
	Message  string `json:"message"`
	Cause    error  `json:"cause,omitempty"`
}

func (e FileProcessingError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("file processing error at %s stage for %s: %s (caused by: %v)", 
			e.Stage, e.FilePath, e.Message, e.Cause)
	}
	return fmt.Sprintf("file processing error at %s stage for %s: %s", 
		e.Stage, e.FilePath, e.Message)
}

// IsValidStudentID checks if a student ID is valid (alphanumeric, not empty)
func (s *StudentData) IsValidStudentID() bool {
	if s.StudentID == "" {
		return false
	}
	
	// Check if it contains only alphanumeric characters
	for _, char := range s.StudentID {
		if !((char >= 'a' && char <= 'z') || 
			 (char >= 'A' && char <= 'Z') || 
			 (char >= '0' && char <= '9')) {
			return false
		}
	}
	return true
}

// GetMarkCount returns the number of valid marks
func (s *StudentData) GetMarkCount() int {
	count := 0
	for _, mark := range s.Marks {
		if mark >= 0 { // Assuming negative marks are invalid
			count++
		}
	}
	return count
}

// String returns a string representation of the student data
func (s *StudentData) String() string {
	return fmt.Sprintf("Student{ID: %s, File: %s, Marks: %d}", 
		s.StudentID, s.FilePath, len(s.Marks))
}
