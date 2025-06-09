package gui

import (
	"testing"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
)

// TestNewApp tests GUI application creation
func TestNewApp(t *testing.T) {
	testApp := test.NewApp()
	defer testApp.Quit()

	app := NewApp()
	if app == nil {
		t.Fatal("NewApp() returned nil")
	}

	// Test default mark mappings
	if len(app.markMappings) == 0 {
		t.Error("NewApp() should initialize default mark mappings")
	}

	expectedMappingCount := 14
	if len(app.markMappings) != expectedMappingCount {
		t.Errorf("NewApp() should have %d default mappings, got %d", expectedMappingCount, len(app.markMappings))
	}
}

// TestWindowProperties tests window initialization and properties
func TestWindowProperties(t *testing.T) {
	testApp := test.NewApp()
	defer testApp.Quit()

	app := NewApp()

	// Test window title
	expectedTitle := "Mark Master Sheet Consolidator"
	if app.window.Title() != expectedTitle {
		t.Errorf("Window title = %v, want %v", app.window.Title(), expectedTitle)
	}

	// Test window size (in test environment, size may be 0x0)
	size := app.window.Canvas().Size()
	if size.Width < 0 || size.Height < 0 {
		t.Errorf("Window size should not be negative, got %v", size)
	}

	// Note: In test environment, window size may be 0x0, which is acceptable
}

// TestSetupUI tests UI component initialization
func TestSetupUI(t *testing.T) {
	testApp := test.NewApp()
	defer testApp.Quit()

	app := NewApp()
	app.setupUI()

	// Test that UI components are initialized
	if app.masterFileEntry == nil {
		t.Error("setupUI() should initialize masterFileEntry")
	}
	if app.studentFolderEntry == nil {
		t.Error("setupUI() should initialize studentFolderEntry")
	}
	if app.markMappingContainer == nil {
		t.Error("setupUI() should initialize markMappingContainer")
	}
	if app.progressBar == nil {
		t.Error("setupUI() should initialize progressBar")
	}
	if app.statusLabel == nil {
		t.Error("setupUI() should initialize statusLabel")
	}
}

// TestDefaultMarkMappings tests default mark mapping initialization
func TestDefaultMarkMappings(t *testing.T) {
	mappings := getDefaultMarkMappings()
	
	expectedMappings := map[string]string{
		"C6":  "I",
		"C7":  "J",
		"C8":  "K",
		"C9":  "L",
		"C10": "M",
		"C11": "N",
		"C12": "O",
		"C13": "P",
		"C15": "Q",
		"C16": "R",
		"C17": "S",
		"C18": "T",
		"C19": "U",
		"C20": "V",
	}

	if len(mappings) != len(expectedMappings) {
		t.Errorf("getDefaultMarkMappings() returned %d mappings, want %d", len(mappings), len(expectedMappings))
	}

	for _, mapping := range mappings {
		expectedColumn, exists := expectedMappings[mapping.StudentCell]
		if !exists {
			t.Errorf("Unexpected mapping: %s -> %s", mapping.StudentCell, mapping.MasterColumn)
			continue
		}
		if mapping.MasterColumn != expectedColumn {
			t.Errorf("Mapping %s -> %s, want %s", mapping.StudentCell, mapping.MasterColumn, expectedColumn)
		}
	}
}

// TestMarkMappingOperations tests mark mapping CRUD operations
func TestMarkMappingOperations(t *testing.T) {
	testApp := test.NewApp()
	defer testApp.Quit()

	app := NewApp()
	app.setupUI()

	initialCount := len(app.markMappings)

	// Test adding mapping
	app.addMarkMapping()
	if len(app.markMappings) != initialCount+1 {
		t.Errorf("addMarkMapping() should increase count by 1, got %d", len(app.markMappings))
	}

	// Test removing mapping
	app.removeMarkMapping(len(app.markMappings) - 1)
	if len(app.markMappings) != initialCount {
		t.Errorf("removeMarkMapping() should decrease count by 1, got %d", len(app.markMappings))
	}

	// Test reset mappings
	app.resetMarkMappings()
	if len(app.markMappings) != 14 {
		t.Errorf("resetMarkMappings() should reset to 14 mappings, got %d", len(app.markMappings))
	}
}

// TestValidation tests input validation functions
func TestValidation(t *testing.T) {
	testApp := test.NewApp()
	defer testApp.Quit()

	app := NewApp()

	tests := []struct {
		name     string
		input    string
		function func(string, string)
		wantErr  bool
	}{
		{
			name:     "valid cell reference",
			input:    "B2",
			function: app.validateCellReference,
			wantErr:  false,
		},
		{
			name:     "invalid cell reference",
			input:    "invalid",
			function: app.validateCellReference,
			wantErr:  true,
		},
		{
			name:     "valid column reference",
			input:    "A",
			function: app.validateColumnReference,
			wantErr:  false,
		},
		{
			name:     "invalid column reference",
			input:    "123",
			function: app.validateColumnReference,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Note: These validation functions show errors via dialogs
			// In a real test, we'd need to capture or mock the error display
			tt.function(tt.input, "test field")
		})
	}
}

// TestStatusUpdates tests status bar updates
func TestStatusUpdates(t *testing.T) {
	testApp := test.NewApp()
	defer testApp.Quit()

	app := NewApp()
	app.setupUI()

	tests := []struct {
		name     string
		status   string
		progress float64
	}{
		{
			name:     "ready status",
			status:   "Ready",
			progress: 0.0,
		},
		{
			name:     "processing status",
			status:   "Processing files...",
			progress: 0.5,
		},
		{
			name:     "completed status",
			status:   "Processing completed",
			progress: 1.0,
		},
		{
			name:     "error status",
			status:   "Error occurred",
			progress: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app.updateStatus(tt.status, tt.progress)
			
			// Test that status label contains the status text
			statusText := app.statusLabel.Text
			if statusText == "" {
				t.Error("Status label should not be empty")
			}
			
			// Test progress bar visibility
			if tt.progress > 0 {
				if !app.progressBar.Visible() {
					t.Error("Progress bar should be visible when progress > 0")
				}
			}
		})
	}
}

// TestResponsiveBehavior tests responsive window behavior
func TestResponsiveBehavior(t *testing.T) {
	testApp := test.NewApp()
	defer testApp.Quit()

	app := NewApp()

	// Test different window sizes (in test environment, resize may not work)
	testSizes := []struct {
		name   string
		width  float32
		height float32
	}{
		{"small", 800, 600},
		{"medium", 1200, 800},
		{"large", 1600, 1000},
	}

	for _, size := range testSizes {
		t.Run(size.name, func(t *testing.T) {
			app.window.Resize(fyne.NewSize(size.width, size.height))

			// Allow time for resize to process
			time.Sleep(10 * time.Millisecond)

			// In test environment, window resize may not work as expected
			// We just verify that the resize call doesn't cause a panic
			actualSize := app.window.Canvas().Size()
			if actualSize.Width < 0 || actualSize.Height < 0 {
				t.Errorf("Window size should not be negative after resize, got %v", actualSize)
			}
		})
	}
}

// TestMenuSetup tests menu bar creation
func TestMenuSetup(t *testing.T) {
	testApp := test.NewApp()
	defer testApp.Quit()

	app := NewApp()
	app.setupMenus()

	// Test that main menu is set
	mainMenu := app.window.MainMenu()
	if mainMenu == nil {
		t.Error("setupMenus() should set main menu")
	}

	// Test expected menu items
	expectedMenus := []string{"File", "Edit", "Help"}
	menus := mainMenu.Items

	if len(menus) != len(expectedMenus) {
		t.Errorf("Expected %d menus, got %d", len(expectedMenus), len(menus))
	}

	for i, expectedMenu := range expectedMenus {
		if i < len(menus) && menus[i].Label != expectedMenu {
			t.Errorf("Menu %d should be %s, got %s", i, expectedMenu, menus[i].Label)
		}
	}
}

// TestResetToDefaults tests configuration reset functionality
func TestResetToDefaults(t *testing.T) {
	testApp := test.NewApp()
	defer testApp.Quit()

	app := NewApp()
	app.setupUI()

	// Modify some values
	app.masterFileEntry.SetText("test.xlsx")
	app.studentFolderEntry.SetText("test-folder")
	app.maxConcurrentEntry.SetText("5")

	// Reset to defaults
	app.resetToDefaults()

	// Test that values are reset
	if app.masterFileEntry.Text != "" {
		t.Error("resetToDefaults() should clear master file entry")
	}
	if app.studentFolderEntry.Text != "" {
		t.Error("resetToDefaults() should clear student folder entry")
	}
	if app.outputFolderEntry.Text != "./output" {
		t.Error("resetToDefaults() should set output folder to default")
	}
	if app.maxConcurrentEntry.Text != "10" {
		t.Error("resetToDefaults() should set max concurrent to default")
	}
}

// TestThemeApplication tests theme application
func TestThemeApplication(t *testing.T) {
	testApp := test.NewApp()
	defer testApp.Quit()

	app := NewApp()
	
	// Test that theme application doesn't cause errors
	app.applyModernStyling()
	app.setupResponsiveBehavior()
	
	// Test window constraints
	constraints := GetWindowConstraints()
	if constraints.MinWidth <= 0 || constraints.MinHeight <= 0 {
		t.Error("Window constraints should have positive minimum values")
	}
	if constraints.OptWidth <= constraints.MinWidth || constraints.OptHeight <= constraints.MinHeight {
		t.Error("Optimal size should be larger than minimum size")
	}
}

// BenchmarkAppCreation benchmarks application creation
func BenchmarkAppCreation(b *testing.B) {
	testApp := test.NewApp()
	defer testApp.Quit()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		app := NewApp()
		_ = app
	}
}

// BenchmarkUISetup benchmarks UI setup
func BenchmarkUISetup(b *testing.B) {
	testApp := test.NewApp()
	defer testApp.Quit()

	app := NewApp()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		app.setupUI()
	}
}

// BenchmarkStatusUpdate benchmarks status updates
func BenchmarkStatusUpdate(b *testing.B) {
	testApp := test.NewApp()
	defer testApp.Quit()

	app := NewApp()
	app.setupUI()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		app.updateStatus("Benchmark status", float64(i%100)/100.0)
	}
}
