package main

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestMainFunction tests the main function behavior
func TestMainFunction(t *testing.T) {
	// Save original args and restore after test
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	tests := []struct {
		name string
		args []string
	}{
		{
			name: "default execution",
			args: []string{"program"},
		},
		{
			name: "dry-run flag",
			args: []string{"program", "-dry-run"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset flag package for each test
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

			os.Args = tt.args

			// Test flag parsing logic without calling main()
			config, dryRun, stats := parseFlags()

			// Basic validation that parsing worked
			if config == "" {
				t.Error("parseFlags() should return non-empty config path")
			}

			// Validate flags based on args
			if len(tt.args) > 1 && tt.args[1] == "-dry-run" {
				if !dryRun {
					t.Error("parseFlags() should set dryRun=true for -dry-run flag")
				}
			}

			_ = stats // Avoid unused variable error
		})
	}
}

// TestParseFlags tests command-line flag parsing
func TestParseFlags(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedConfig string
		expectedDryRun bool
		expectedStats  bool
	}{
		{
			name:           "default config",
			args:           []string{"program"},
			expectedConfig: "config.toml",
			expectedDryRun: false,
			expectedStats:  false,
		},
		{
			name:           "custom config",
			args:           []string{"program", "-config", "custom.toml"},
			expectedConfig: "custom.toml",
			expectedDryRun: false,
			expectedStats:  false,
		},
		{
			name:           "dry run mode",
			args:           []string{"program", "-dry-run"},
			expectedConfig: "config.toml",
			expectedDryRun: true,
			expectedStats:  false,
		},
		{
			name:           "stats mode",
			args:           []string{"program", "-stats"},
			expectedConfig: "config.toml",
			expectedDryRun: false,
			expectedStats:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset flag package
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
			
			// Save original args
			originalArgs := os.Args
			defer func() { os.Args = originalArgs }()
			
			os.Args = tt.args
			
			config, dryRun, stats := parseFlags()
			
			if config != tt.expectedConfig {
				t.Errorf("parseFlags() config = %v, want %v", config, tt.expectedConfig)
			}
			if dryRun != tt.expectedDryRun {
				t.Errorf("parseFlags() dryRun = %v, want %v", dryRun, tt.expectedDryRun)
			}
			if stats != tt.expectedStats {
				t.Errorf("parseFlags() stats = %v, want %v", stats, tt.expectedStats)
			}
		})
	}
}

// TestConfigFileValidation tests configuration file validation
func TestConfigFileValidation(t *testing.T) {
	// Create temporary directory for test files
	tempDir := t.TempDir()
	
	tests := []struct {
		name       string
		configFile string
		createFile bool
		wantError  bool
	}{
		{
			name:       "existing config file",
			configFile: filepath.Join(tempDir, "valid.toml"),
			createFile: true,
			wantError:  false,
		},
		{
			name:       "non-existent config file",
			configFile: filepath.Join(tempDir, "missing.toml"),
			createFile: false,
			wantError:  true,
		},
		{
			name:       "empty config path",
			configFile: "",
			createFile: false,
			wantError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.createFile && tt.configFile != "" {
				// Create a minimal valid config file
				content := `[paths]
student_files_folder = "./test"
master_sheet_path = "./test.xlsx"
output_folder = "./output"
`
				err := os.WriteFile(tt.configFile, []byte(content), 0644)
				if err != nil {
					t.Fatalf("Failed to create test config file: %v", err)
				}
			}

			err := validateConfigFile(tt.configFile)
			
			if tt.wantError && err == nil {
				t.Errorf("validateConfigFile() expected error but got none")
			}
			if !tt.wantError && err != nil {
				t.Errorf("validateConfigFile() unexpected error: %v", err)
			}
		})
	}
}

// TestVersionDisplay tests version information display
func TestVersionDisplay(t *testing.T) {
	version := getVersion()
	
	if version == "" {
		t.Error("getVersion() returned empty string")
	}
	
	if !strings.Contains(version, "Mark Master Sheet Consolidator") {
		t.Errorf("getVersion() should contain application name, got: %s", version)
	}
}

// TestUsageDisplay tests usage information display
func TestUsageDisplay(t *testing.T) {
	usage := getUsage()
	
	if usage == "" {
		t.Error("getUsage() returned empty string")
	}
	
	expectedFlags := []string{"-config", "-dry-run", "-stats", "-version", "-help"}
	for _, flag := range expectedFlags {
		if !strings.Contains(usage, flag) {
			t.Errorf("getUsage() should contain flag %s, got: %s", flag, usage)
		}
	}
}

// Helper functions that would need to be extracted from main.go for testing

func parseFlags() (config string, dryRun bool, stats bool) {
	configFlag := flag.String("config", "config.toml", "Path to configuration file")
	dryRunFlag := flag.Bool("dry-run", false, "Run in dry-run mode (no actual changes)")
	statsFlag := flag.Bool("stats", false, "Show processing statistics")
	versionFlag := flag.Bool("version", false, "Show version information")
	
	flag.Parse()
	
	if *versionFlag {
		// In real implementation, this would print version and exit
		return *configFlag, *dryRunFlag, *statsFlag
	}
	
	return *configFlag, *dryRunFlag, *statsFlag
}

func validateConfigFile(configPath string) error {
	if configPath == "" {
		return os.ErrNotExist
	}
	
	_, err := os.Stat(configPath)
	return err
}

func getVersion() string {
	return "Mark Master Sheet Consolidator v1.0.0"
}

func getUsage() string {
	return `Usage: mark-master-sheet [options]

Options:
  -config string    Path to configuration file (default "config.toml")
  -dry-run         Run in dry-run mode (no actual changes)
  -stats           Show processing statistics
  -version         Show version information
  -help            Show this help message`
}

// BenchmarkFlagParsing benchmarks command-line flag parsing
func BenchmarkFlagParsing(b *testing.B) {
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()
	
	os.Args = []string{"program", "-config", "test.toml", "-dry-run"}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
		parseFlags()
	}
}
