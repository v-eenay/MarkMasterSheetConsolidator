package gui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"fyne.io/fyne/v2/test"

	"mark-master-sheet/internal/config"
)

// TestLoadDefaultConfig tests default configuration loading
func TestLoadDefaultConfig(t *testing.T) {
	testApp := test.NewApp()
	defer testApp.Quit()

	app := NewApp()
	app.setupUI()
	app.loadDefaultConfig()

	// Test default values
	expectedDefaults := map[string]string{
		"outputFolder":         "./output",
		"backupFolder":         "./backups",
		"studentWorksheet":     "Grading Sheet",
		"masterWorksheet":      "001",
		"studentIDCell":        "B2",
		"studentIDColumn":      "B",
		"maxConcurrent":        "10",
	}

	tests := []struct {
		name     string
		actual   string
		expected string
	}{
		{"output folder", app.outputFolderEntry.Text, expectedDefaults["outputFolder"]},
		{"backup folder", app.backupFolderEntry.Text, expectedDefaults["backupFolder"]},
		{"student worksheet", app.studentWorksheetEntry.Text, expectedDefaults["studentWorksheet"]},
		{"master worksheet", app.masterWorksheetEntry.Text, expectedDefaults["masterWorksheet"]},
		{"student ID cell", app.studentIDCellEntry.Text, expectedDefaults["studentIDCell"]},
		{"student ID column", app.studentIDColumnEntry.Text, expectedDefaults["studentIDColumn"]},
		{"max concurrent", app.maxConcurrentEntry.Text, expectedDefaults["maxConcurrent"]},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.actual != tt.expected {
				t.Errorf("loadDefaultConfig() %s = %v, want %v", tt.name, tt.actual, tt.expected)
			}
		})
	}

	// Test checkboxes
	if !app.enableBackupCheck.Checked {
		t.Error("loadDefaultConfig() should enable backup by default")
	}
	if !app.skipInvalidCheck.Checked {
		t.Error("loadDefaultConfig() should enable skip invalid files by default")
	}
}

// TestBuildConfigFromUI tests building configuration from UI values
func TestBuildConfigFromUI(t *testing.T) {
	testApp := test.NewApp()
	defer testApp.Quit()

	app := NewApp()
	app.setupUI()

	// Set test values
	app.masterFileEntry.SetText("test-master.xlsx")
	app.studentFolderEntry.SetText("./students")
	app.outputFolderEntry.SetText("./output")
	app.backupFolderEntry.SetText("./backups")
	app.studentWorksheetEntry.SetText("Grading Sheet")
	app.masterWorksheetEntry.SetText("001")
	app.studentIDCellEntry.SetText("B2")
	app.studentIDColumnEntry.SetText("B")
	app.maxConcurrentEntry.SetText("5")
	app.enableBackupCheck.SetChecked(true)
	app.skipInvalidCheck.SetChecked(false)

	// Set some mark mappings
	app.markMappings = []MarkMapping{
		{StudentCell: "C6", MasterColumn: "I"},
		{StudentCell: "C7", MasterColumn: "J"},
	}

	cfg, err := app.buildConfigFromUI()
	if err != nil {
		t.Fatalf("buildConfigFromUI() unexpected error: %v", err)
	}

	// Test configuration values
	if cfg.Paths.MasterSheetPath != "test-master.xlsx" {
		t.Errorf("buildConfigFromUI() master sheet path = %v, want %v", cfg.Paths.MasterSheetPath, "test-master.xlsx")
	}
	if cfg.Paths.StudentFilesFolder != "./students" {
		t.Errorf("buildConfigFromUI() student folder = %v, want %v", cfg.Paths.StudentFilesFolder, "./students")
	}
	if cfg.Processing.MaxConcurrentFiles != 5 {
		t.Errorf("buildConfigFromUI() max concurrent = %v, want %v", cfg.Processing.MaxConcurrentFiles, 5)
	}
	if !cfg.Processing.BackupEnabled {
		t.Error("buildConfigFromUI() should enable backup")
	}
	if cfg.Processing.SkipInvalidFiles {
		t.Error("buildConfigFromUI() should disable skip invalid files")
	}

	// Test mark mappings
	if len(cfg.Excel.MarkCells) != 2 {
		t.Errorf("buildConfigFromUI() mark cells count = %v, want %v", len(cfg.Excel.MarkCells), 2)
	}
	if len(cfg.Excel.MasterColumns) != 2 {
		t.Errorf("buildConfigFromUI() master columns count = %v, want %v", len(cfg.Excel.MasterColumns), 2)
	}
}

// TestBuildConfigFromUIValidation tests configuration validation
func TestBuildConfigFromUIValidation(t *testing.T) {
	testApp := test.NewApp()
	defer testApp.Quit()

	tests := []struct {
		name           string
		setupFunc      func(*App)
		expectedError  string
	}{
		{
			name: "missing master file",
			setupFunc: func(app *App) {
				app.masterFileEntry.SetText("")
				app.studentFolderEntry.SetText("./students")
				app.outputFolderEntry.SetText("./output")
			},
			expectedError: "master file path is required",
		},
		{
			name: "missing student folder",
			setupFunc: func(app *App) {
				app.masterFileEntry.SetText("master.xlsx")
				app.studentFolderEntry.SetText("")
				app.outputFolderEntry.SetText("./output")
			},
			expectedError: "student files folder is required",
		},
		{
			name: "missing output folder",
			setupFunc: func(app *App) {
				app.masterFileEntry.SetText("master.xlsx")
				app.studentFolderEntry.SetText("./students")
				app.outputFolderEntry.SetText("")
			},
			expectedError: "output folder is required",
		},
		{
			name: "invalid max concurrent",
			setupFunc: func(app *App) {
				app.masterFileEntry.SetText("master.xlsx")
				app.studentFolderEntry.SetText("./students")
				app.outputFolderEntry.SetText("./output")
				app.maxConcurrentEntry.SetText("invalid")
			},
			expectedError: "max concurrent files must be between 1 and 20",
		},
		{
			name: "no mark mappings",
			setupFunc: func(app *App) {
				app.masterFileEntry.SetText("master.xlsx")
				app.studentFolderEntry.SetText("./students")
				app.outputFolderEntry.SetText("./output")
				app.maxConcurrentEntry.SetText("10")
				app.markMappings = []MarkMapping{} // Empty mappings
			},
			expectedError: "at least one mark mapping is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := NewApp()
			app.setupUI()
			tt.setupFunc(app)

			_, err := app.buildConfigFromUI()
			if err == nil {
				t.Errorf("buildConfigFromUI() expected error containing %q but got none", tt.expectedError)
				return
			}

			if !strings.Contains(err.Error(), tt.expectedError) {
				t.Errorf("buildConfigFromUI() error = %v, want error containing %q", err, tt.expectedError)
			}
		})
	}
}

// TestApplyConfigToUI tests applying configuration to UI elements
func TestApplyConfigToUI(t *testing.T) {
	testApp := test.NewApp()
	defer testApp.Quit()

	app := NewApp()
	app.setupUI()

	// Create test configuration
	cfg := &config.Config{
		Paths: config.PathsConfig{
			MasterSheetPath:    "test-master.xlsx",
			StudentFilesFolder: "./test-students",
			OutputFolder:       "./test-output",
			BackupFolder:       "./test-backups",
		},
		Excel: config.ExcelConfig{
			StudentWorksheetName: "Test Sheet",
			MasterWorksheetName:  "Test Master",
			StudentIDCell:        "A1",
			MarkCells:            []string{"B1", "B2"},
			MasterColumns:        []string{"X", "Y"},
		},
		Processing: config.ProcessingConfig{
			MaxConcurrentFiles: 15,
			BackupEnabled:      false,
			SkipInvalidFiles:   true,
		},
	}

	app.applyConfigToUI(cfg)

	// Test that UI values match configuration
	if app.masterFileEntry.Text != cfg.Paths.MasterSheetPath {
		t.Errorf("applyConfigToUI() master file = %v, want %v", app.masterFileEntry.Text, cfg.Paths.MasterSheetPath)
	}
	if app.studentFolderEntry.Text != cfg.Paths.StudentFilesFolder {
		t.Errorf("applyConfigToUI() student folder = %v, want %v", app.studentFolderEntry.Text, cfg.Paths.StudentFilesFolder)
	}
	if app.studentWorksheetEntry.Text != cfg.Excel.StudentWorksheetName {
		t.Errorf("applyConfigToUI() student worksheet = %v, want %v", app.studentWorksheetEntry.Text, cfg.Excel.StudentWorksheetName)
	}
	if app.maxConcurrentEntry.Text != "15" {
		t.Errorf("applyConfigToUI() max concurrent = %v, want %v", app.maxConcurrentEntry.Text, "15")
	}
	if app.enableBackupCheck.Checked != cfg.Processing.BackupEnabled {
		t.Errorf("applyConfigToUI() backup enabled = %v, want %v", app.enableBackupCheck.Checked, cfg.Processing.BackupEnabled)
	}

	// Test mark mappings
	if len(app.markMappings) != len(cfg.Excel.MarkCells) {
		t.Errorf("applyConfigToUI() mark mappings count = %v, want %v", len(app.markMappings), len(cfg.Excel.MarkCells))
	}
}

// TestSaveConfigToPath tests saving configuration to file
func TestSaveConfigToPath(t *testing.T) {
	testApp := test.NewApp()
	defer testApp.Quit()

	app := NewApp()
	app.setupUI()

	// Create test configuration
	cfg := &config.Config{
		Paths: config.PathsConfig{
			MasterSheetPath:    "test.xlsx",
			StudentFilesFolder: "./students",
			OutputFolder:       "./output",
			LogFolder:          "./logs",
			BackupFolder:       "./backups",
		},
		Excel: config.ExcelConfig{
			StudentWorksheetName: "Grading Sheet",
			MasterWorksheetName:  "001",
			StudentIDCell:        "B2",
			MarkCells:            []string{"C6", "C7"},
			MasterColumns:        []string{"I", "J"},
		},
		Processing: config.ProcessingConfig{
			MaxConcurrentFiles: 10,
			BackupEnabled:      true,
			SkipInvalidFiles:   true,
			TimeoutSeconds:     300,
			RetryAttempts:      3,
		},
		Logging: config.LoggingConfig{
			Level:          "INFO",
			ConsoleOutput:  true,
			FileOutput:     true,
			MaxFileSizeMB:  100,
			MaxBackupFiles: 5,
			MaxAgeDays:     30,
		},
	}

	// Save to temporary file
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "test-config.toml")

	err := app.saveConfigToPath(cfg, configPath)
	if err != nil {
		t.Fatalf("saveConfigToPath() unexpected error: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("saveConfigToPath() did not create config file")
	}

	// Verify file content
	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}

	contentStr := string(content)
	expectedStrings := []string{
		"[paths]",
		"student_files_folder",
		"master_sheet_path",
		"[excel_settings]",
		"student_worksheet_name",
		"[processing]",
		"max_concurrent_files",
		"[logging]",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(contentStr, expected) {
			t.Errorf("saveConfigToPath() config file should contain %q", expected)
		}
	}
}

// TestFormatStringArray tests string array formatting for TOML
func TestFormatStringArray(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected string
	}{
		{
			name:     "empty array",
			input:    []string{},
			expected: "",
		},
		{
			name:     "single item",
			input:    []string{"item1"},
			expected: `"item1"`,
		},
		{
			name:     "multiple items",
			input:    []string{"item1", "item2", "item3"},
			expected: `"item1", "item2", "item3"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatStringArray(tt.input)
			if result != tt.expected {
				t.Errorf("formatStringArray() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// BenchmarkBuildConfigFromUI benchmarks configuration building
func BenchmarkBuildConfigFromUI(b *testing.B) {
	testApp := test.NewApp()
	defer testApp.Quit()

	app := NewApp()
	app.setupUI()
	app.loadDefaultConfig()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := app.buildConfigFromUI()
		if err != nil {
			b.Fatalf("buildConfigFromUI failed: %v", err)
		}
	}
}

// BenchmarkApplyConfigToUI benchmarks configuration application
func BenchmarkApplyConfigToUI(b *testing.B) {
	testApp := test.NewApp()
	defer testApp.Quit()

	app := NewApp()
	app.setupUI()

	cfg := &config.Config{
		Paths: config.PathsConfig{
			MasterSheetPath:    "test.xlsx",
			StudentFilesFolder: "./students",
			OutputFolder:       "./output",
			BackupFolder:       "./backups",
		},
		Excel: config.ExcelConfig{
			StudentWorksheetName: "Grading Sheet",
			MasterWorksheetName:  "001",
			StudentIDCell:        "B2",
			MarkCells:            []string{"C6", "C7", "C8"},
			MasterColumns:        []string{"I", "J", "K"},
		},
		Processing: config.ProcessingConfig{
			MaxConcurrentFiles: 10,
			BackupEnabled:      true,
			SkipInvalidFiles:   true,
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		app.applyConfigToUI(cfg)
	}
}
