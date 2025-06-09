package excel

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/xuri/excelize/v2"

	"mark-master-sheet/internal/config"
)

// TestNewReader tests Excel reader creation
func TestNewReader(t *testing.T) {
	tests := []struct {
		name   string
		config *config.ExcelConfig
	}{
		{
			name: "valid config",
			config: &config.ExcelConfig{
				StudentWorksheetName: "Grading Sheet",
				StudentIDCell:        "B2",
				MarkCells:            []string{"C6", "C7"},
			},
		},
		{
			name: "minimal config",
			config: &config.ExcelConfig{
				StudentWorksheetName: "Sheet1",
				StudentIDCell:        "A1",
				MarkCells:            []string{"B1"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := NewReader(tt.config)

			if reader == nil {
				t.Error("NewReader() returned nil reader")
			}

			if reader.config != tt.config {
				t.Error("NewReader() should store config reference")
			}
		})
	}
}

// TestReadStudentData tests reading student data from Excel files
func TestReadStudentData(t *testing.T) {
	// Create test Excel file with student data
	testFile := createTestStudentFile(t)
	defer os.Remove(testFile)

	tests := []struct {
		name          string
		config        *config.ExcelConfig
		wantError     bool
		expectedID    string
		expectedMarks int
	}{
		{
			name: "valid student data",
			config: &config.ExcelConfig{
				StudentWorksheetName: "Grading Sheet",
				StudentIDCell:        "B2",
				MarkCells:            []string{"C6", "C7", "C8"},
			},
			wantError:     false,
			expectedID:    "STU001",
			expectedMarks: 3,
		},
		{
			name: "invalid worksheet",
			config: &config.ExcelConfig{
				StudentWorksheetName: "NonExistent",
				StudentIDCell:        "B2",
				MarkCells:            []string{"C6"},
			},
			wantError: true,
		},
		{
			name: "empty student ID cell",
			config: &config.ExcelConfig{
				StudentWorksheetName: "Grading Sheet",
				StudentIDCell:        "Z99",
				MarkCells:            []string{"C6"},
			},
			wantError: true, // Should error on empty student ID
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := NewReader(tt.config)
			studentData, err := reader.ReadStudentData(testFile)

			if tt.wantError {
				if err == nil {
					t.Errorf("ReadStudentData() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("ReadStudentData() unexpected error: %v", err)
				return
			}

			if studentData.StudentID != tt.expectedID {
				t.Errorf("ReadStudentData() student ID = %v, want %v", studentData.StudentID, tt.expectedID)
			}

			if len(studentData.Marks) != tt.expectedMarks {
				t.Errorf("ReadStudentData() marks count = %v, want %v", len(studentData.Marks), tt.expectedMarks)
			}
		})
	}
}

// TestFindStudentInMasterSheet tests finding student in master sheet
func TestFindStudentInMasterSheet(t *testing.T) {
	// Create test master file
	masterFile := createTestMasterFile(t)
	defer os.Remove(masterFile)

	config := &config.ExcelConfig{
		MasterWorksheetName: "001",
	}
	reader := NewReader(config)

	// Open master file
	file, err := excelize.OpenFile(masterFile)
	if err != nil {
		t.Fatalf("Failed to open master file: %v", err)
	}
	defer file.Close()

	tests := []struct {
		name        string
		studentID   string
		expectedRow int
		wantError   bool
	}{
		{
			name:        "existing student",
			studentID:   "STU001",
			expectedRow: 2,
			wantError:   false,
		},
		{
			name:        "non-existent student",
			studentID:   "STU999",
			expectedRow: 0,
			wantError:   true,
		},
		{
			name:        "case insensitive match",
			studentID:   "stu001",
			expectedRow: 2,
			wantError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			row, err := reader.FindStudentInMasterSheet(file, tt.studentID)

			if tt.wantError {
				if err == nil {
					t.Errorf("FindStudentInMasterSheet() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("FindStudentInMasterSheet() unexpected error: %v", err)
				return
			}

			if row != tt.expectedRow {
				t.Errorf("FindStudentInMasterSheet() = %v, want %v", row, tt.expectedRow)
			}
		})
	}
}

// Helper functions for creating test files

func createTestMasterFile(t *testing.T) string {
	f := excelize.NewFile()
	defer f.Close()

	// Create "001" worksheet
	sheetName := "001"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		t.Fatalf("Failed to create worksheet: %v", err)
	}
	f.SetActiveSheet(index)

	// Add test data - headers and student IDs
	f.SetCellValue(sheetName, "A1", "Name")
	f.SetCellValue(sheetName, "B1", "Student ID")
	f.SetCellValue(sheetName, "I1", "Mark 1")
	f.SetCellValue(sheetName, "J1", "Mark 2")
	f.SetCellValue(sheetName, "K1", "Mark 3")

	// Add sample students
	f.SetCellValue(sheetName, "A2", "John Doe")
	f.SetCellValue(sheetName, "B2", "STU001")
	f.SetCellValue(sheetName, "A3", "Jane Smith")
	f.SetCellValue(sheetName, "B3", "STU002")

	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "master.xlsx")

	if err := f.SaveAs(filePath); err != nil {
		t.Fatalf("Failed to create test master file: %v", err)
	}

	return filePath
}

func createTestStudentFile(t *testing.T) string {
	f := excelize.NewFile()
	defer f.Close()
	
	// Create "Grading Sheet" worksheet
	sheetName := "Grading Sheet"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		t.Fatalf("Failed to create worksheet: %v", err)
	}
	f.SetActiveSheet(index)
	
	// Add test data
	f.SetCellValue(sheetName, "B2", "STU001")  // Student ID
	f.SetCellValue(sheetName, "C6", 85)        // Mark 1
	f.SetCellValue(sheetName, "C7", 92)        // Mark 2
	f.SetCellValue(sheetName, "C8", 78)        // Mark 3
	
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "student.xlsx")
	
	if err := f.SaveAs(filePath); err != nil {
		t.Fatalf("Failed to create test student file: %v", err)
	}
	
	return filePath
}

// BenchmarkReadStudentData benchmarks student data reading
func BenchmarkReadStudentData(b *testing.B) {
	// Create test file
	f := excelize.NewFile()
	defer f.Close()

	sheetName := "Grading Sheet"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		b.Fatalf("Failed to create worksheet: %v", err)
	}
	f.SetActiveSheet(index)

	// Add test data
	f.SetCellValue(sheetName, "B2", "STU001")
	f.SetCellValue(sheetName, "C6", 85)
	f.SetCellValue(sheetName, "C7", 92)
	f.SetCellValue(sheetName, "C8", 78)

	tempDir := b.TempDir()
	testFile := filepath.Join(tempDir, "student.xlsx")

	if err := f.SaveAs(testFile); err != nil {
		b.Fatalf("Failed to create test file: %v", err)
	}

	config := &config.ExcelConfig{
		StudentWorksheetName: "Grading Sheet",
		StudentIDCell:        "B2",
		MarkCells:            []string{"C6", "C7", "C8"},
	}
	reader := NewReader(config)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := reader.ReadStudentData(testFile)
		if err != nil {
			b.Fatalf("ReadStudentData failed: %v", err)
		}
	}
}

// BenchmarkFindStudentInMasterSheet benchmarks student finding
func BenchmarkFindStudentInMasterSheet(b *testing.B) {
	// Create test master file
	f := excelize.NewFile()
	defer f.Close()

	sheetName := "001"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		b.Fatalf("Failed to create worksheet: %v", err)
	}
	f.SetActiveSheet(index)

	// Add test data
	f.SetCellValue(sheetName, "B2", "STU001")
	f.SetCellValue(sheetName, "B3", "STU002")

	tempDir := b.TempDir()
	testFile := filepath.Join(tempDir, "master.xlsx")

	if err := f.SaveAs(testFile); err != nil {
		b.Fatalf("Failed to create test file: %v", err)
	}

	config := &config.ExcelConfig{
		MasterWorksheetName: "001",
	}
	reader := NewReader(config)

	// Open file for benchmarking
	file, err := excelize.OpenFile(testFile)
	if err != nil {
		b.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := reader.FindStudentInMasterSheet(file, "STU001")
		if err != nil {
			b.Fatalf("FindStudentInMasterSheet failed: %v", err)
		}
	}
}
