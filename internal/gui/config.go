package gui

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"

	"mark-master-sheet/internal/config"
)

// loadDefaultConfig loads default configuration values
func (a *App) loadDefaultConfig() {
	// Set default values
	a.outputFolderEntry.SetText("./output")
	a.backupFolderEntry.SetText("./backups")
	
	a.studentWorksheetEntry.SetText("Grading Sheet")
	a.masterWorksheetEntry.SetText("001")
	a.studentIDCellEntry.SetText("B2")
	a.studentIDColumnEntry.SetText("B")
	
	a.enableBackupCheck.SetChecked(true)
	a.skipInvalidCheck.SetChecked(true)
	a.maxConcurrentEntry.SetText("10")
	
	a.updateStatus("Default configuration loaded")
}

// loadConfigFromFile loads configuration from a TOML file
func (a *App) loadConfigFromFile() {
	dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil || reader == nil {
			return
		}
		defer reader.Close()

		configPath := reader.URI().Path()

		cfg, err := config.LoadConfig(configPath)
		if err != nil {
			a.showError(fmt.Sprintf("Failed to load configuration: %v", err))
			return
		}

		a.applyConfigToUI(cfg)
		a.updateStatus(fmt.Sprintf("Configuration loaded from %s", filepath.Base(configPath)))

	}, a.window)
}

// saveConfigToFile saves current configuration to a TOML file
func (a *App) saveConfigToFile() {
	dialog.ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err != nil || writer == nil {
			return
		}
		defer writer.Close()

		cfg, err := a.buildConfigFromUI()
		if err != nil {
			a.showError(fmt.Sprintf("Configuration validation failed: %v", err))
			return
		}

		// Convert config to TOML and save
		configPath := writer.URI().Path()
		err = a.saveConfigToPath(cfg, configPath)
		if err != nil {
			a.showError(fmt.Sprintf("Failed to save configuration: %v", err))
			return
		}

		a.updateStatus(fmt.Sprintf("Configuration saved to %s", filepath.Base(configPath)))

	}, a.window)
}

// applyConfigToUI applies loaded configuration to UI elements
func (a *App) applyConfigToUI(cfg *config.Config) {
	// File paths
	a.masterFileEntry.SetText(cfg.Paths.MasterSheetPath)
	a.studentFolderEntry.SetText(cfg.Paths.StudentFilesFolder)
	a.outputFolderEntry.SetText(cfg.Paths.OutputFolder)
	a.backupFolderEntry.SetText(cfg.Paths.BackupFolder)
	
	// Excel settings
	a.studentWorksheetEntry.SetText(cfg.Excel.StudentWorksheetName)
	a.masterWorksheetEntry.SetText(cfg.Excel.MasterWorksheetName)
	a.studentIDCellEntry.SetText(cfg.Excel.StudentIDCell)
	
	// Extract column from master columns (assuming first column is the student ID column)
	if len(cfg.Excel.MasterColumns) > 0 {
		// For simplicity, we'll use "B" as default since that's where student IDs typically are
		a.studentIDColumnEntry.SetText("B")
	}
	
	// Processing settings
	a.enableBackupCheck.SetChecked(cfg.Processing.BackupEnabled)
	a.skipInvalidCheck.SetChecked(cfg.Processing.SkipInvalidFiles)
	a.maxConcurrentEntry.SetText(fmt.Sprintf("%d", cfg.Processing.MaxConcurrentFiles))
	
	// Mark mappings
	if len(cfg.Excel.MarkCells) == len(cfg.Excel.MasterColumns) {
		a.markMappings = make([]MarkMapping, len(cfg.Excel.MarkCells))
		for i, cell := range cfg.Excel.MarkCells {
			a.markMappings[i] = MarkMapping{
				StudentCell:  cell,
				MasterColumn: cfg.Excel.MasterColumns[i],
			}
		}
		a.refreshMarkMappingsDisplay()
	}
}

// buildConfigFromUI builds a configuration object from current UI values
func (a *App) buildConfigFromUI() (*config.Config, error) {
	// Validate required fields
	if a.masterFileEntry.Text == "" {
		return nil, fmt.Errorf("master file path is required")
	}
	if a.studentFolderEntry.Text == "" {
		return nil, fmt.Errorf("student files folder is required")
	}
	if a.outputFolderEntry.Text == "" {
		return nil, fmt.Errorf("output folder is required")
	}
	
	// Parse max concurrent files
	maxConcurrent, err := strconv.Atoi(a.maxConcurrentEntry.Text)
	if err != nil || maxConcurrent < 1 || maxConcurrent > 20 {
		return nil, fmt.Errorf("max concurrent files must be between 1 and 20")
	}
	
	// Build mark cells and columns from mappings
	var markCells []string
	var masterColumns []string
	for _, mapping := range a.markMappings {
		if mapping.StudentCell != "" && mapping.MasterColumn != "" {
			markCells = append(markCells, mapping.StudentCell)
			masterColumns = append(masterColumns, mapping.MasterColumn)
		}
	}
	
	if len(markCells) == 0 {
		return nil, fmt.Errorf("at least one mark mapping is required")
	}
	
	// Create configuration
	cfg := &config.Config{
		Paths: config.PathsConfig{
			StudentFilesFolder: a.studentFolderEntry.Text,
			MasterSheetPath:    a.masterFileEntry.Text,
			OutputFolder:       a.outputFolderEntry.Text,
			LogFolder:          "./logs",
			BackupFolder:       a.backupFolderEntry.Text,
		},
		Excel: config.ExcelConfig{
			StudentWorksheetName: a.studentWorksheetEntry.Text,
			MasterWorksheetName:  a.masterWorksheetEntry.Text,
			StudentIDCell:        a.studentIDCellEntry.Text,
			MarkCells:            markCells,
			MasterColumns:        masterColumns,
		},
		Processing: config.ProcessingConfig{
			MaxConcurrentFiles: maxConcurrent,
			BackupEnabled:      a.enableBackupCheck.Checked,
			SkipInvalidFiles:   a.skipInvalidCheck.Checked,
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
	
	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	
	return cfg, nil
}

// saveConfigToPath saves configuration to the specified file path
func (a *App) saveConfigToPath(cfg *config.Config, configPath string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	
	// For now, we'll create a simple TOML representation
	// In a full implementation, you'd use a TOML encoder
	content := fmt.Sprintf(`[paths]
student_files_folder = "%s"
master_sheet_path = "%s"
output_folder = "%s"
log_folder = "%s"
backup_folder = "%s"

[excel_settings]
student_worksheet_name = "%s"
master_worksheet_name = "%s"
student_id_cell = "%s"
mark_cells = [%s]
master_columns = [%s]

[processing]
max_concurrent_files = %d
backup_enabled = %t
skip_invalid_files = %t
timeout_seconds = %d
retry_attempts = %d

[logging]
level = "%s"
console_output = %t
file_output = %t
max_file_size_mb = %d
max_backup_files = %d
max_age_days = %d
`,
		cfg.Paths.StudentFilesFolder,
		cfg.Paths.MasterSheetPath,
		cfg.Paths.OutputFolder,
		cfg.Paths.LogFolder,
		cfg.Paths.BackupFolder,
		cfg.Excel.StudentWorksheetName,
		cfg.Excel.MasterWorksheetName,
		cfg.Excel.StudentIDCell,
		formatStringArray(cfg.Excel.MarkCells),
		formatStringArray(cfg.Excel.MasterColumns),
		cfg.Processing.MaxConcurrentFiles,
		cfg.Processing.BackupEnabled,
		cfg.Processing.SkipInvalidFiles,
		cfg.Processing.TimeoutSeconds,
		cfg.Processing.RetryAttempts,
		cfg.Logging.Level,
		cfg.Logging.ConsoleOutput,
		cfg.Logging.FileOutput,
		cfg.Logging.MaxFileSizeMB,
		cfg.Logging.MaxBackupFiles,
		cfg.Logging.MaxAgeDays,
	)
	
	return os.WriteFile(configPath, []byte(content), 0644)
}

// formatStringArray formats a string array for TOML
func formatStringArray(arr []string) string {
	if len(arr) == 0 {
		return ""
	}
	
	result := ""
	for i, s := range arr {
		if i > 0 {
			result += ", "
		}
		result += fmt.Sprintf(`"%s"`, s)
	}
	return result
}
