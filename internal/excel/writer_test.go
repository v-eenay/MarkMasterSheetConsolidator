package excel

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/xuri/excelize/v2"

	"mark-master-sheet/internal/config"
	"mark-master-sheet/pkg/models"
)

// TestNewWriter tests Excel writer creation
func TestNewWriter(t *testing.T) {
	tests := []struct {
		name   string
		config *config.ExcelConfig
	}{
		{
			name: "valid config",
			config: &config.ExcelConfig{
				MasterWorksheetName: "001",
				MarkCells:           []string{"C6", "C7"},
				MasterColumns:       []string{"I", "J"},
			},
		},
		{
			name: "minimal config",
			config: &config.ExcelConfig{
				MasterWorksheetName: "Sheet1",
				MarkCells:           []string{"B1"},
				MasterColumns:       []string{"C"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := NewWriter(tt.config)

			if writer == nil {
				t.Error("NewWriter() returned nil writer")
			}

			if writer.config != tt.config {
				t.Error("NewWriter() should store config reference")
			}
		})
	}
}

// TestUpdateMasterSheet tests updating master sheet with student data
func TestUpdateMasterSheet(t *testing.T) {
	testFile := createTestMasterFileForWriter(t)
	defer os.Remove(testFile)

	config := &config.ExcelConfig{
		MasterWorksheetName: "001",
		MarkCells:           []string{"C6", "C7", "C8"},
		MasterColumns:       []string{"I", "J", "K"},
	}
	writer := NewWriter(config)

	tests := []struct {
		name        string
		studentData *models.StudentData
		wantError   bool
	}{
		{
			name: "valid student update",
			studentData: &models.StudentData{
				StudentID: "STU001",
				Marks: map[string]float64{
					"C6": 85.5,
					"C7": 92.0,
					"C8": 78.5,
				},
			},
			wantError: false,
		},
		{
			name: "student not found",
			studentData: &models.StudentData{
				StudentID: "STU999",
				Marks: map[string]float64{
					"C6": 85.5,
				},
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := writer.UpdateMasterSheet(testFile, tt.studentData)

			if tt.wantError {
				if err == nil {
					t.Errorf("UpdateMasterSheet() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("UpdateMasterSheet() unexpected error: %v", err)
				return
			}
		})
	}
}

// TestValidateMasterSheet tests master sheet validation
func TestValidateMasterSheet(t *testing.T) {
	testFile := createTestMasterFileForWriter(t)
	defer os.Remove(testFile)

	config := &config.ExcelConfig{
		MasterWorksheetName: "001",
	}
	writer := NewWriter(config)

	tests := []struct {
		name      string
		filePath  string
		wantError bool
	}{
		{
			name:      "valid master sheet",
			filePath:  testFile,
			wantError: false,
		},
		{
			name:      "non-existent file",
			filePath:  "nonexistent.xlsx",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := writer.ValidateMasterSheet(tt.filePath)

			if tt.wantError {
				if err == nil {
					t.Errorf("ValidateMasterSheet() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("ValidateMasterSheet() unexpected error: %v", err)
				return
			}
		})
	}
}

// TestWriterBasicFunctionality tests basic writer functionality
func TestWriterBasicFunctionality(t *testing.T) {
	config := &config.ExcelConfig{
		MasterWorksheetName: "001",
		MarkCells:           []string{"C6", "C7", "C8"},
		MasterColumns:       []string{"I", "J", "K"},
	}
	writer := NewWriter(config)

	if writer == nil {
		t.Error("NewWriter() should return valid writer")
	}

	if writer.config != config {
		t.Error("NewWriter() should store config reference")
	}

	// Test that writer has expected configuration
	if writer.config.MasterWorksheetName != "001" {
		t.Errorf("Writer config worksheet = %v, want %v", writer.config.MasterWorksheetName, "001")
	}

	if len(writer.config.MarkCells) != 3 {
		t.Errorf("Writer config mark cells count = %v, want %v", len(writer.config.MarkCells), 3)
	}

	if len(writer.config.MasterColumns) != 3 {
		t.Errorf("Writer config master columns count = %v, want %v", len(writer.config.MasterColumns), 3)
	}
}

// Helper function to create test master file for writer tests
func createTestMasterFileForWriter(t *testing.T) string {
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
	filePath := filepath.Join(tempDir, "master_writer.xlsx")

	if err := f.SaveAs(filePath); err != nil {
		t.Fatalf("Failed to create test master file: %v", err)
	}

	return filePath
}

// BenchmarkUpdateMasterSheet benchmarks master sheet updating
func BenchmarkUpdateMasterSheet(b *testing.B) {
	testFile := createTestMasterFileForWriterBench(b)
	defer os.Remove(testFile)

	config := &config.ExcelConfig{
		MasterWorksheetName: "001",
		MarkCells:           []string{"C6", "C7", "C8"},
		MasterColumns:       []string{"I", "J", "K"},
	}
	writer := NewWriter(config)

	studentData := &models.StudentData{
		StudentID: "STU001",
		Marks: map[string]float64{
			"C6": 85.5,
			"C7": 92.0,
			"C8": 78.5,
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := writer.UpdateMasterSheet(testFile, studentData)
		if err != nil {
			b.Fatalf("UpdateMasterSheet failed: %v", err)
		}
	}
}

// Helper function to create test master file for writer benchmarks
func createTestMasterFileForWriterBench(t testing.TB) string {
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
	filePath := filepath.Join(tempDir, "master_writer.xlsx")

	if err := f.SaveAs(filePath); err != nil {
		t.Fatalf("Failed to create test master file: %v", err)
	}

	return filePath
}

// BenchmarkWriterCreation benchmarks writer creation
func BenchmarkWriterCreation(b *testing.B) {
	config := &config.ExcelConfig{
		MasterWorksheetName: "001",
		MarkCells:           []string{"C6", "C7", "C8"},
		MasterColumns:       []string{"I", "J", "K"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		writer := NewWriter(config)
		if writer == nil {
			b.Fatalf("NewWriter failed")
		}
	}
}


