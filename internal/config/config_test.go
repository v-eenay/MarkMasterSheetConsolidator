package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: Config{
				Paths: PathsConfig{
					StudentFilesFolder: "./students",
					MasterSheetPath:    "./master.xlsx",
					OutputFolder:       "./output",
				},
				Excel: ExcelConfig{
					MarkCells:     []string{"C6", "C7"},
					MasterColumns: []string{"I", "J"},
				},
				Processing: ProcessingConfig{
					MaxConcurrentFiles: 5,
					TimeoutSeconds:     300,
				},
			},
			wantErr: false,
		},
		{
			name: "empty student files folder",
			config: Config{
				Paths: PathsConfig{
					StudentFilesFolder: "",
					MasterSheetPath:    "./master.xlsx",
					OutputFolder:       "./output",
				},
				Excel: ExcelConfig{
					MarkCells:     []string{"C6", "C7"},
					MasterColumns: []string{"I", "J"},
				},
				Processing: ProcessingConfig{
					MaxConcurrentFiles: 5,
					TimeoutSeconds:     300,
				},
			},
			wantErr: true,
		},
		{
			name: "mismatched mark cells and columns",
			config: Config{
				Paths: PathsConfig{
					StudentFilesFolder: "./students",
					MasterSheetPath:    "./master.xlsx",
					OutputFolder:       "./output",
				},
				Excel: ExcelConfig{
					MarkCells:     []string{"C6", "C7", "C8"},
					MasterColumns: []string{"I", "J"}, // One less than mark cells
				},
				Processing: ProcessingConfig{
					MaxConcurrentFiles: 5,
					TimeoutSeconds:     300,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid max concurrent files",
			config: Config{
				Paths: PathsConfig{
					StudentFilesFolder: "./students",
					MasterSheetPath:    "./master.xlsx",
					OutputFolder:       "./output",
				},
				Excel: ExcelConfig{
					MarkCells:     []string{"C6", "C7"},
					MasterColumns: []string{"I", "J"},
				},
				Processing: ProcessingConfig{
					MaxConcurrentFiles: 0, // Invalid
					TimeoutSeconds:     300,
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLoadConfig(t *testing.T) {
	// Create a temporary config file
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "test_config.toml")

	configContent := `
[paths]
student_files_folder = "./students"
master_sheet_path = "./master.xlsx"
output_folder = "./output"
log_folder = "./logs"
backup_folder = "./backups"

[excel_settings]
student_worksheet_name = "Grading Sheet"
master_worksheet_name = "001"
student_id_cell = "B2"
mark_cells = ["C6", "C7"]
master_columns = ["I", "J"]

[processing]
max_concurrent_files = 5
backup_enabled = true
skip_invalid_files = true
timeout_seconds = 300
retry_attempts = 3

[logging]
level = "INFO"
console_output = true
file_output = true
max_file_size_mb = 100
max_backup_files = 5
max_age_days = 30
`

	err := os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	// Test loading the config
	config, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	// Verify some key values
	if config.Paths.StudentFilesFolder == "" {
		t.Error("StudentFilesFolder should not be empty")
	}

	if config.Excel.StudentWorksheetName != "Grading Sheet" {
		t.Errorf("Expected StudentWorksheetName to be 'Grading Sheet', got %s", config.Excel.StudentWorksheetName)
	}

	if config.Processing.MaxConcurrentFiles != 5 {
		t.Errorf("Expected MaxConcurrentFiles to be 5, got %d", config.Processing.MaxConcurrentFiles)
	}

	if len(config.Excel.MarkCells) != 2 {
		t.Errorf("Expected 2 mark cells, got %d", len(config.Excel.MarkCells))
	}
}

func TestLoadConfig_FileNotFound(t *testing.T) {
	_, err := LoadConfig("nonexistent.toml")
	if err == nil {
		t.Error("Expected error for nonexistent config file")
	}
}

func TestConfig_EnsureDirectories(t *testing.T) {
	tempDir := t.TempDir()

	config := &Config{
		Paths: PathsConfig{
			OutputFolder: filepath.Join(tempDir, "output"),
			LogFolder:    filepath.Join(tempDir, "logs"),
			BackupFolder: filepath.Join(tempDir, "backups"),
		},
	}

	err := config.EnsureDirectories()
	if err != nil {
		t.Fatalf("EnsureDirectories() error = %v", err)
	}

	// Check if directories were created
	dirs := []string{
		config.Paths.OutputFolder,
		config.Paths.LogFolder,
		config.Paths.BackupFolder,
	}

	for _, dir := range dirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			t.Errorf("Directory %s was not created", dir)
		}
	}
}

func TestConfig_ResolvePaths(t *testing.T) {
	config := &Config{
		Paths: PathsConfig{
			StudentFilesFolder: "./students",
			MasterSheetPath:    "./master.xlsx",
			OutputFolder:       "./output",
			LogFolder:          "./logs",
			BackupFolder:       "./backups",
		},
	}

	err := config.ResolvePaths()
	if err != nil {
		t.Fatalf("ResolvePaths() error = %v", err)
	}

	// Check that paths are now absolute
	if !filepath.IsAbs(config.Paths.StudentFilesFolder) {
		t.Error("StudentFilesFolder should be absolute after ResolvePaths()")
	}

	if !filepath.IsAbs(config.Paths.MasterSheetPath) {
		t.Error("MasterSheetPath should be absolute after ResolvePaths()")
	}

	if !filepath.IsAbs(config.Paths.OutputFolder) {
		t.Error("OutputFolder should be absolute after ResolvePaths()")
	}
}
