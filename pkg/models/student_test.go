package models

import (
	"testing"
	"time"
)

func TestStudentData_IsValidStudentID(t *testing.T) {
	tests := []struct {
		name      string
		studentID string
		want      bool
	}{
		{
			name:      "valid alphanumeric ID",
			studentID: "23049191",
			want:      true,
		},
		{
			name:      "valid mixed alphanumeric ID",
			studentID: "ABC123",
			want:      true,
		},
		{
			name:      "empty ID",
			studentID: "",
			want:      false,
		},
		{
			name:      "ID with spaces",
			studentID: "230 491",
			want:      false,
		},
		{
			name:      "ID with special characters",
			studentID: "230-491",
			want:      false,
		},
		{
			name:      "ID with lowercase letters",
			studentID: "abc123",
			want:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StudentData{
				StudentID: tt.studentID,
			}
			if got := s.IsValidStudentID(); got != tt.want {
				t.Errorf("StudentData.IsValidStudentID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStudentData_GetMarkCount(t *testing.T) {
	tests := []struct {
		name  string
		marks map[string]float64
		want  int
	}{
		{
			name: "all valid marks",
			marks: map[string]float64{
				"C6": 85.5,
				"C7": 90.0,
				"C8": 78.5,
			},
			want: 3,
		},
		{
			name: "some invalid marks",
			marks: map[string]float64{
				"C6": 85.5,
				"C7": -1.0, // Invalid (negative)
				"C8": 78.5,
			},
			want: 2,
		},
		{
			name:  "empty marks",
			marks: map[string]float64{},
			want:  0,
		},
		{
			name: "all invalid marks",
			marks: map[string]float64{
				"C6": -1.0,
				"C7": -2.0,
			},
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StudentData{
				Marks: tt.marks,
			}
			if got := s.GetMarkCount(); got != tt.want {
				t.Errorf("StudentData.GetMarkCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStudentData_String(t *testing.T) {
	s := &StudentData{
		StudentID: "23049191",
		FilePath:  "/path/to/file.xlsx",
		Marks: map[string]float64{
			"C6": 85.5,
			"C7": 90.0,
		},
	}

	result := s.String()
	expected := "Student{ID: 23049191, File: /path/to/file.xlsx, Marks: 2}"
	
	if result != expected {
		t.Errorf("StudentData.String() = %v, want %v", result, expected)
	}
}

func TestValidationError_Error(t *testing.T) {
	err := ValidationError{
		Field:   "student_id",
		Value:   "invalid-id",
		Message: "contains invalid characters",
		File:    "/path/to/file.xlsx",
	}

	expected := "validation error in file /path/to/file.xlsx, field student_id (value: invalid-id): contains invalid characters"
	if got := err.Error(); got != expected {
		t.Errorf("ValidationError.Error() = %v, want %v", got, expected)
	}
}

func TestFileProcessingError_Error(t *testing.T) {
	tests := []struct {
		name string
		err  FileProcessingError
		want string
	}{
		{
			name: "error with cause",
			err: FileProcessingError{
				FilePath: "/path/to/file.xlsx",
				Stage:    "reading",
				Message:  "failed to open file",
				Cause:    &ValidationError{Message: "file not found"},
			},
			want: "file processing error at reading stage for /path/to/file.xlsx: failed to open file (caused by: validation error in file , field  (value: ): file not found)",
		},
		{
			name: "error without cause",
			err: FileProcessingError{
				FilePath: "/path/to/file.xlsx",
				Stage:    "validation",
				Message:  "invalid format",
			},
			want: "file processing error at validation stage for /path/to/file.xlsx: invalid format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("FileProcessingError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProcessingSummary(t *testing.T) {
	summary := &ProcessingSummary{
		TotalFiles:       100,
		SuccessfulFiles:  95,
		FailedFiles:      3,
		SkippedFiles:     2,
		StudentsUpdated:  90,
		StudentsNotFound: 5,
		StartTime:        time.Now().Add(-5 * time.Minute),
		EndTime:          time.Now(),
		Errors:           []string{"error1", "error2"},
		Warnings:         []string{"warning1"},
	}

	// Test that all fields are properly set
	if summary.TotalFiles != 100 {
		t.Errorf("Expected TotalFiles to be 100, got %d", summary.TotalFiles)
	}

	if summary.SuccessfulFiles != 95 {
		t.Errorf("Expected SuccessfulFiles to be 95, got %d", summary.SuccessfulFiles)
	}

	if len(summary.Errors) != 2 {
		t.Errorf("Expected 2 errors, got %d", len(summary.Errors))
	}

	if len(summary.Warnings) != 1 {
		t.Errorf("Expected 1 warning, got %d", len(summary.Warnings))
	}
}
