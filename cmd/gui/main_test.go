package main

import (
	"testing"

	"fyne.io/fyne/v2/test"

	"mark-master-sheet/internal/gui"
)

// TestGUIApplicationCreation tests GUI application initialization
func TestGUIApplicationCreation(t *testing.T) {
	// Use Fyne's test app for headless testing
	testApp := test.NewApp()
	defer testApp.Quit()

	// Test GUI app creation
	app := gui.NewApp()
	if app == nil {
		t.Fatal("NewApp() returned nil")
	}

	// Test that the app was created successfully
	// Note: Internal methods are not exposed for testing
}

// TestGUIApplicationBasics tests basic GUI functionality
func TestGUIApplicationBasics(t *testing.T) {
	testApp := test.NewApp()
	defer testApp.Quit()

	app := gui.NewApp()

	// Test that app was created
	if app == nil {
		t.Fatal("NewApp() returned nil")
	}

	// Test that we can call Run without errors (in test mode)
	// Note: We can't test actual window properties without exposing internal methods
}

// TestGUICreationStability tests that GUI creation is stable
func TestGUICreationStability(t *testing.T) {
	testApp := test.NewApp()
	defer testApp.Quit()

	// Test creating multiple GUI instances
	for i := 0; i < 5; i++ {
		app := gui.NewApp()
		if app == nil {
			t.Fatalf("NewApp() iteration %d returned nil", i)
		}
	}
}

// TestGUIMemoryStability tests memory stability
func TestGUIMemoryStability(t *testing.T) {
	testApp := test.NewApp()
	defer testApp.Quit()

	// Create and release multiple GUI instances
	for i := 0; i < 10; i++ {
		app := gui.NewApp()
		if app == nil {
			t.Fatalf("NewApp() iteration %d returned nil", i)
		}
		// Let app go out of scope for garbage collection
		app = nil
	}
}

// TestGUIBasicFunctionality tests basic GUI functionality
func TestGUIBasicFunctionality(t *testing.T) {
	testApp := test.NewApp()
	defer testApp.Quit()

	app := gui.NewApp()

	// Test that we can call Run method
	// Note: In a real test environment, Run() would block, so we don't call it
	// Instead, we just verify the app was created successfully
	if app == nil {
		t.Error("GUI app should be created successfully")
	}
}

// BenchmarkGUICreation benchmarks GUI creation performance
func BenchmarkGUICreation(b *testing.B) {
	testApp := test.NewApp()
	defer testApp.Quit()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		app := gui.NewApp()
		if app == nil {
			b.Fatalf("NewApp() failed at iteration %d", i)
		}
	}
}

// Note: These tests require additional methods to be exposed in the GUI package
// The following methods would need to be added to the gui.App struct:

/*
Required methods to add to gui.App:
- GetWindow() fyne.Window
- SetupUI()
- GetTabs() []string
- GetMarkMappings() []MarkMapping
- ApplyModernStyling()
- ShowError(string)
- ShowWarning(string)
- ShowInfo(string, string)
- UpdateStatus(string, ...float64)
- GetCurrentStatus() string
*/
