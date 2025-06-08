// Package main provides the command-line interface for the Mark Master Sheet Consolidator.
// This application automates the consolidation of student marks from individual Excel files
// into a master spreadsheet with comprehensive error handling and logging.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"mark-master-sheet/internal/config"
	"mark-master-sheet/internal/logger"
	"mark-master-sheet/internal/processor"
	"mark-master-sheet/pkg/models"
)

var (
	configPath = flag.String("config", "config.toml", "Path to configuration file")
	dryRun     = flag.Bool("dry-run", false, "Run in dry-run mode (no actual changes)")
	showStats  = flag.Bool("stats", false, "Show processing statistics and exit")
	version    = flag.Bool("version", false, "Show version information")
)

const (
	appName    = "Mark Master Sheet Consolidator"
	appVersion = "1.0.0"
	appAuthor  = "Vinay Koirala"
)

func main() {
	flag.Parse()

	// Show version and exit
	if *version {
		fmt.Printf("%s v%s\n", appName, appVersion)
		fmt.Printf("Author: %s\n", appAuthor)
		fmt.Printf("Repository: https://github.com/v-eenay/MarkMasterSheetConsolidator.git\n")
		os.Exit(0)
	}

	// Load configuration
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Ensure required directories exist
	if err := cfg.EnsureDirectories(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create directories: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	log, err := logger.NewLogger(&cfg.Logging, cfg.Paths.LogFolder)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	log.Info("=== Mark Master Sheet Consolidator Started ===")
	log.WithField("version", appVersion).Info("Application version")
	log.WithField("config_path", *configPath).Info("Configuration loaded")

	if *dryRun {
		log.Info("Running in DRY-RUN mode - no changes will be made")
	}

	// Create processor
	proc := processor.NewProcessor(cfg, log)

	// Show statistics and exit if requested
	if *showStats {
		stats := proc.GetProcessingStatistics()
		fmt.Println("=== Processing Statistics ===")
		for key, value := range stats {
			fmt.Printf("%s: %v\n", key, value)
		}
		os.Exit(0)
	}

	// Setup graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle interrupt signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		log.WithField("signal", sig).Info("Received shutdown signal")
		cancel()
	}()

	// Add timeout to context
	if cfg.Processing.TimeoutSeconds > 0 {
		timeoutCtx, timeoutCancel := context.WithTimeout(ctx, 
			time.Duration(cfg.Processing.TimeoutSeconds)*time.Second)
		defer timeoutCancel()
		ctx = timeoutCtx
	}

	// Validate master sheet exists
	if _, err := os.Stat(cfg.Paths.MasterSheetPath); os.IsNotExist(err) {
		log.WithField("path", cfg.Paths.MasterSheetPath).Fatal("Master sheet file not found")
	}

	// Validate student files folder exists
	if _, err := os.Stat(cfg.Paths.StudentFilesFolder); os.IsNotExist(err) {
		log.WithField("path", cfg.Paths.StudentFilesFolder).Fatal("Student files folder not found")
	}

	// Start processing
	log.Info("Starting file processing...")
	startTime := time.Now()

	summary, err := proc.ProcessFiles(ctx, *dryRun)
	if err != nil {
		log.WithError(err).Fatal("Processing failed")
	}

	// Log final summary
	duration := time.Since(startTime)
	log.WithField("duration", duration).Info("Processing completed")

	// Print summary to console
	printSummary(summary, *dryRun)

	// Exit with appropriate code
	if summary.FailedFiles > 0 {
		log.Warn("Processing completed with errors")
		os.Exit(1)
	}

	log.Info("=== Mark Master Sheet Consolidator Completed Successfully ===")
}

// printSummary prints a formatted summary to the console
func printSummary(summary interface{}, dryRun bool) {
	fmt.Println("\n=== Processing Summary ===")

	if dryRun {
		fmt.Println("Mode: DRY-RUN (no changes made)")
	} else {
		fmt.Println("Mode: PRODUCTION")
	}

	// Type assertion to access summary fields
	if s, ok := summary.(*models.ProcessingSummary); ok {
		fmt.Printf("Total Files: %d\n", s.TotalFiles)
		fmt.Printf("Successful: %d\n", s.SuccessfulFiles)
		fmt.Printf("Failed: %d\n", s.FailedFiles)
		fmt.Printf("Skipped: %d\n", s.SkippedFiles)
		
		if !dryRun {
			fmt.Printf("Students Updated: %d\n", s.StudentsUpdated)
			fmt.Printf("Students Not Found: %d\n", s.StudentsNotFound)
		}
		
		fmt.Printf("Duration: %v\n", s.TotalDuration)

		if len(s.Errors) > 0 {
			fmt.Printf("\nErrors (%d):\n", len(s.Errors))
			for i, err := range s.Errors {
				if i < 5 { // Show only first 5 errors
					fmt.Printf("  - %s\n", err)
				} else {
					fmt.Printf("  ... and %d more errors\n", len(s.Errors)-5)
					break
				}
			}
		}

		if len(s.Warnings) > 0 {
			fmt.Printf("\nWarnings (%d):\n", len(s.Warnings))
			for i, warning := range s.Warnings {
				if i < 5 { // Show only first 5 warnings
					fmt.Printf("  - %s\n", warning)
				} else {
					fmt.Printf("  ... and %d more warnings\n", len(s.Warnings)-5)
					break
				}
			}
		}
	}

	fmt.Println("========================")
}
