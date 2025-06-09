package gui

import (
	"context"
	"fmt"
	"os"
	"time"

	"mark-master-sheet/internal/config"
	"mark-master-sheet/internal/logger"
	"mark-master-sheet/internal/processor"
)

// startProcessing begins the file processing operation
func (a *App) startProcessing(dryRun bool) {
	if a.isProcessing {
		a.showError("Processing is already in progress")
		return
	}
	
	// Build configuration from UI
	cfg, err := a.buildConfigFromUI()
	if err != nil {
		a.showError(fmt.Sprintf("Configuration error: %v", err))
		return
	}
	
	// Validate file paths exist
	if err := a.validatePaths(cfg); err != nil {
		a.showError(fmt.Sprintf("Path validation failed: %v", err))
		return
	}
	
	// Initialize logger
	if err := cfg.EnsureDirectories(); err != nil {
		a.showError(fmt.Sprintf("Failed to create directories: %v", err))
		return
	}
	
	logger, err := logger.NewLogger(&cfg.Logging, cfg.Paths.LogFolder)
	if err != nil {
		a.showError(fmt.Sprintf("Failed to initialize logger: %v", err))
		return
	}
	
	// Initialize processor
	proc := processor.NewProcessor(cfg, logger)
	
	// Set up processing state
	a.isProcessing = true
	a.config = cfg
	a.logger = logger
	a.processor = proc
	
	// Create cancellable context
	a.processingContext, a.cancelProcessing = context.WithCancel(context.Background())
	
	// Update UI
	a.updateProcessingUI(true, dryRun)
	
	// Start processing in goroutine
	go a.runProcessing(dryRun)
}

// stopProcessing cancels the current processing operation
func (a *App) stopProcessing() {
	if !a.isProcessing {
		return
	}
	
	if a.cancelProcessing != nil {
		a.cancelProcessing()
	}
	
	a.updateStatus("Stopping processing...")
	a.appendLog("Processing cancelled by user\n")
}

// runProcessing executes the actual processing logic
func (a *App) runProcessing(dryRun bool) {
	defer func() {
		a.isProcessing = false
		a.updateProcessingUI(false, dryRun)
	}()
	
	startTime := time.Now()
	
	if dryRun {
		a.updateStatus("Starting dry run...")
		a.appendLog("=== DRY RUN MODE - No changes will be made ===\n")
	} else {
		a.updateStatus("Starting file processing...")
		a.appendLog("=== PROCESSING MODE - Files will be updated ===\n")
	}
	
	// Run the processing
	summary, err := a.processor.ProcessFiles(a.processingContext, dryRun)
	
	duration := time.Since(startTime)
	
	if err != nil {
		if a.processingContext.Err() == context.Canceled {
			a.updateStatus("Processing cancelled")
			a.appendLog(fmt.Sprintf("Processing cancelled after %v\n", duration))
		} else {
			a.updateStatus("Processing failed")
			a.appendLog(fmt.Sprintf("Processing failed: %v\n", err))
			a.showError(fmt.Sprintf("Processing failed: %v", err))
		}
		return
	}
	
	// Display results
	a.updateStatus("Processing completed")
	a.displayProcessingSummary(summary, dryRun, duration)
	
	if !dryRun && summary.FailedFiles == 0 {
		a.showInfo("Success", "Processing completed successfully!")
	} else if summary.FailedFiles > 0 {
		a.showInfo("Completed with Warnings", 
			fmt.Sprintf("Processing completed with %d failed files. Check logs for details.", summary.FailedFiles))
	}
}

// validatePaths validates that required paths exist
func (a *App) validatePaths(cfg *config.Config) error {
	// Check master sheet exists
	if _, err := os.Stat(cfg.Paths.MasterSheetPath); os.IsNotExist(err) {
		return fmt.Errorf("master sheet file not found: %s", cfg.Paths.MasterSheetPath)
	}
	
	// Check student files folder exists
	if _, err := os.Stat(cfg.Paths.StudentFilesFolder); os.IsNotExist(err) {
		return fmt.Errorf("student files folder not found: %s", cfg.Paths.StudentFilesFolder)
	}
	
	return nil
}

// updateProcessingUI updates UI elements during processing
func (a *App) updateProcessingUI(processing bool, dryRun bool) {
	// Find the processing buttons (this is a simplified approach)
	// In a real implementation, you'd store references to these buttons
	if processing {
		a.progressBar.Show()
		a.progressBar.SetValue(0)
		if dryRun {
			a.updateStatus("Running dry run...")
		} else {
			a.updateStatus("Processing files...")
		}
	} else {
		a.progressBar.Hide()
		a.updateStatus("Ready")
	}
}

// appendLog adds text to the log output area
func (a *App) appendLog(text string) {
	currentText := a.logOutput.Text
	a.logOutput.SetText(currentText + text)
	
	// Auto-scroll to bottom
	a.logOutput.CursorRow = len(a.logOutput.Text)
	a.logOutput.Refresh()
}

// displayProcessingSummary shows the processing results
func (a *App) displayProcessingSummary(summary interface{}, dryRun bool, duration time.Duration) {
	a.appendLog("\n=== PROCESSING SUMMARY ===\n")
	
	if dryRun {
		a.appendLog("Mode: DRY RUN (no changes made)\n")
	} else {
		a.appendLog("Mode: PRODUCTION\n")
	}
	
	// Type assertion to access summary fields
	// Note: This would need to be adjusted based on the actual summary type
	a.appendLog(fmt.Sprintf("Duration: %v\n", duration))
	a.appendLog("Processing completed.\n")
	a.appendLog("Check the logs folder for detailed processing information.\n")
	a.appendLog("========================\n\n")
}

// Additional helper methods for UI updates during processing
func (a *App) updateProgress(current, total int, currentFile string) {
	if total > 0 {
		progress := float64(current) / float64(total)
		a.progressBar.SetValue(progress)
		
		percentage := int(progress * 100)
		a.updateStatus(fmt.Sprintf("Processing... %d%% (%d/%d)", percentage, current, total))
		
		if currentFile != "" {
			a.appendLog(fmt.Sprintf("Processing: %s\n", currentFile))
		}
	}
}

// logError logs an error message
func (a *App) logError(message string) {
	a.appendLog(fmt.Sprintf("ERROR: %s\n", message))
}

// logWarning logs a warning message
func (a *App) logWarning(message string) {
	a.appendLog(fmt.Sprintf("WARNING: %s\n", message))
}

// logInfo logs an info message
func (a *App) logInfo(message string) {
	a.appendLog(fmt.Sprintf("INFO: %s\n", message))
}
