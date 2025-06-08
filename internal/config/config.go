// Package config provides configuration management for the Mark Master Sheet Consolidator.
// It handles loading, validation, and path resolution for TOML-based configuration files.
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Config represents the application configuration
type Config struct {
	Paths      PathsConfig      `toml:"paths"`
	Excel      ExcelConfig      `toml:"excel_settings"`
	Processing ProcessingConfig `toml:"processing"`
	Logging    LoggingConfig    `toml:"logging"`
}

// PathsConfig contains file and directory paths
type PathsConfig struct {
	StudentFilesFolder string `toml:"student_files_folder"`
	MasterSheetPath    string `toml:"master_sheet_path"`
	OutputFolder       string `toml:"output_folder"`
	LogFolder          string `toml:"log_folder"`
	BackupFolder       string `toml:"backup_folder"`
}

// ExcelConfig contains Excel-specific settings
type ExcelConfig struct {
	StudentWorksheetName string   `toml:"student_worksheet_name"`
	MasterWorksheetName  string   `toml:"master_worksheet_name"`
	StudentIDCell        string   `toml:"student_id_cell"`
	MarkCells            []string `toml:"mark_cells"`
	MasterColumns        []string `toml:"master_columns"`
}

// ProcessingConfig contains processing-related settings
type ProcessingConfig struct {
	MaxConcurrentFiles int  `toml:"max_concurrent_files"`
	BackupEnabled      bool `toml:"backup_enabled"`
	SkipInvalidFiles   bool `toml:"skip_invalid_files"`
	TimeoutSeconds     int  `toml:"timeout_seconds"`
	RetryAttempts      int  `toml:"retry_attempts"`
}

// LoggingConfig contains logging settings
type LoggingConfig struct {
	Level          string `toml:"level"`
	ConsoleOutput  bool   `toml:"console_output"`
	FileOutput     bool   `toml:"file_output"`
	MaxFileSizeMB  int    `toml:"max_file_size_mb"`
	MaxBackupFiles int    `toml:"max_backup_files"`
	MaxAgeDays     int    `toml:"max_age_days"`
}

// LoadConfig loads configuration from the specified file
func LoadConfig(configPath string) (*Config, error) {
	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("configuration file not found: %s", configPath)
	}

	var config Config
	if _, err := toml.DecodeFile(configPath, &config); err != nil {
		return nil, fmt.Errorf("failed to decode configuration file: %w", err)
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	// Convert relative paths to absolute paths
	if err := config.ResolvePaths(); err != nil {
		return nil, fmt.Errorf("failed to resolve paths: %w", err)
	}

	return &config, nil
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	// Validate paths
	if c.Paths.StudentFilesFolder == "" {
		return fmt.Errorf("student_files_folder cannot be empty")
	}
	if c.Paths.MasterSheetPath == "" {
		return fmt.Errorf("master_sheet_path cannot be empty")
	}
	if c.Paths.OutputFolder == "" {
		return fmt.Errorf("output_folder cannot be empty")
	}

	// Validate Excel settings
	if len(c.Excel.MarkCells) != len(c.Excel.MasterColumns) {
		return fmt.Errorf("mark_cells and master_columns must have the same length")
	}
	if len(c.Excel.MarkCells) == 0 {
		return fmt.Errorf("mark_cells cannot be empty")
	}

	// Validate processing settings
	if c.Processing.MaxConcurrentFiles <= 0 {
		return fmt.Errorf("max_concurrent_files must be greater than 0")
	}
	if c.Processing.TimeoutSeconds <= 0 {
		return fmt.Errorf("timeout_seconds must be greater than 0")
	}

	return nil
}

// ResolvePaths converts relative paths to absolute paths
func (c *Config) ResolvePaths() error {
	var err error

	// Resolve student files folder
	if c.Paths.StudentFilesFolder, err = filepath.Abs(c.Paths.StudentFilesFolder); err != nil {
		return fmt.Errorf("failed to resolve student_files_folder: %w", err)
	}

	// Resolve master sheet path
	if c.Paths.MasterSheetPath, err = filepath.Abs(c.Paths.MasterSheetPath); err != nil {
		return fmt.Errorf("failed to resolve master_sheet_path: %w", err)
	}

	// Resolve output folder
	if c.Paths.OutputFolder, err = filepath.Abs(c.Paths.OutputFolder); err != nil {
		return fmt.Errorf("failed to resolve output_folder: %w", err)
	}

	// Resolve log folder
	if c.Paths.LogFolder, err = filepath.Abs(c.Paths.LogFolder); err != nil {
		return fmt.Errorf("failed to resolve log_folder: %w", err)
	}

	// Resolve backup folder
	if c.Paths.BackupFolder, err = filepath.Abs(c.Paths.BackupFolder); err != nil {
		return fmt.Errorf("failed to resolve backup_folder: %w", err)
	}

	return nil
}

// EnsureDirectories creates necessary directories if they don't exist
func (c *Config) EnsureDirectories() error {
	dirs := []string{
		c.Paths.OutputFolder,
		c.Paths.LogFolder,
		c.Paths.BackupFolder,
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}
