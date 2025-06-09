// Package gui provides the graphical user interface for the Mark Master Sheet Consolidator.
// It uses the Fyne framework to create a user-friendly interface for configuration and processing.
package gui

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"mark-master-sheet/internal/config"
	"mark-master-sheet/internal/logger"
	"mark-master-sheet/internal/processor"
)

// App represents the main GUI application
type App struct {
	fyneApp    fyne.App
	window     fyne.Window
	config     *config.Config
	logger     *logger.Logger
	processor  *processor.Processor
	
	// UI Components
	masterFileEntry     *widget.Entry
	studentFolderEntry  *widget.Entry
	outputFolderEntry   *widget.Entry
	backupFolderEntry   *widget.Entry
	
	studentWorksheetEntry *widget.Entry
	masterWorksheetEntry  *widget.Entry
	studentIDCellEntry    *widget.Entry
	studentIDColumnEntry  *widget.Entry
	
	markMappingTable    *widget.List
	markMappings        []MarkMapping
	
	enableBackupCheck   *widget.Check
	skipInvalidCheck    *widget.Check
	maxConcurrentEntry  *widget.Entry
	
	progressBar         *widget.ProgressBar
	statusLabel         *widget.Label
	logOutput          *widget.Entry
	
	// Processing state
	isProcessing        bool
	processingContext   context.Context
	cancelProcessing    context.CancelFunc
}

// MarkMapping represents a mapping between student file cell and master sheet column
type MarkMapping struct {
	StudentCell  string
	MasterColumn string
}

// NewApp creates a new GUI application instance with modern design
func NewApp() *App {
	fyneApp := app.NewWithID("com.vinaykoirala.markmaster")

	window := fyneApp.NewWindow("Mark Master Sheet Consolidator")

	// Apply responsive window sizing with constraints
	constraints := GetWindowConstraints()
	window.Resize(fyne.NewSize(constraints.OptWidth, constraints.OptHeight))
	window.SetFixedSize(false)

	// Set window constraints (Fyne doesn't have direct min/max size, but we handle it in resize)
	window.Canvas().SetOnTypedKey(func(key *fyne.KeyEvent) {
		// Handle keyboard shortcuts and window management
		if key.Name == fyne.KeyF11 {
			// Toggle fullscreen (if supported)
		}
	})

	window.SetMaster()

	app := &App{
		fyneApp: fyneApp,
		window:  window,
		markMappings: getDefaultMarkMappings(),
	}

	// Apply modern theme and responsive behavior
	app.setupResponsiveBehavior()

	return app
}

// getDefaultMarkMappings returns the default mark cell to column mappings
func getDefaultMarkMappings() []MarkMapping {
	return []MarkMapping{
		{"C6", "I"}, {"C7", "J"}, {"C8", "K"}, {"C9", "L"}, {"C10", "M"},
		{"C11", "N"}, {"C12", "O"}, {"C13", "P"}, {"C15", "Q"}, {"C16", "R"},
		{"C17", "S"}, {"C18", "T"}, {"C19", "U"}, {"C20", "V"},
	}
}

// Run starts the GUI application
func (a *App) Run() {
	a.setupUI()
	a.setupMenus()
	a.loadDefaultConfig()
	a.window.ShowAndRun()
}

// setupUI creates and arranges all UI components with modern design
func (a *App) setupUI() {
	// Create main container with enhanced tabs
	tabs := container.NewAppTabs(
		container.NewTabItem("File Paths", a.createFilePathsTab()),
		container.NewTabItem("Excel Settings", a.createExcelSettingsTab()),
		container.NewTabItem("Mark Mappings", a.createMarkMappingsTab()),
		container.NewTabItem("Processing", a.createProcessingTab()),
		container.NewTabItem("Logs", a.createLogsTab()),
	)

	// Set tab location and styling
	tabs.SetTabLocation(container.TabLocationTop)

	// Create enhanced status bar
	statusBar := a.createStatusBar()

	// Create header with application title and info
	header := a.createHeader()

	// Main layout with responsive design
	content := container.NewBorder(
		header,        // top
		statusBar,     // bottom
		nil,           // left
		nil,           // right
		container.NewPadded(tabs), // center with padding
	)

	a.window.SetContent(content)

	// Apply modern theme and styling
	a.applyModernStyling()
}

// createHeader creates a modern header with application branding
func (a *App) createHeader() *fyne.Container {
	// Application title with modern styling
	titleLabel := widget.NewLabel("Mark Master Sheet Consolidator")
	titleLabel.TextStyle = fyne.TextStyle{Bold: true}

	// Subtitle with version info
	subtitleLabel := widget.NewLabel("Professional Excel Processing Tool v1.0.0")

	// Author info
	authorLabel := widget.NewLabel("by Vinay Koirala")

	// Create header layout
	titleContainer := container.NewVBox(
		titleLabel,
		subtitleLabel,
		authorLabel,
	)

	// Add some spacing and styling
	header := container.NewBorder(
		nil, nil, nil, nil,
		container.NewPadded(titleContainer),
	)

	return header
}

// applyModernStyling applies modern visual styling to the application
func (a *App) applyModernStyling() {
	// Apply custom theme if available
	a.applyCustomTheme()

	// Set window icon (if available)
	// a.window.SetIcon(resourceIconPng) // Uncomment when icon is available
}

// setupResponsiveBehavior configures responsive window behavior
func (a *App) setupResponsiveBehavior() {
	constraints := GetWindowConstraints()

	// Monitor window resize events for responsive behavior
	a.window.Canvas().SetOnTypedRune(func(r rune) {
		// Handle responsive layout changes based on window size
		size := a.window.Canvas().Size()

		// Adjust layout based on window size
		if size.Width < constraints.MinWidth || size.Height < constraints.MinHeight {
			// Compact layout for small screens
			a.adjustForSmallScreen()
		} else if size.Width > constraints.OptWidth {
			// Expanded layout for large screens
			a.adjustForLargeScreen()
		}
	})
}

// adjustForSmallScreen optimizes layout for smaller screens
func (a *App) adjustForSmallScreen() {
	// Implement compact layout adjustments
	// This could include hiding certain elements or changing layouts
}

// adjustForLargeScreen optimizes layout for larger screens
func (a *App) adjustForLargeScreen() {
	// Implement expanded layout adjustments
	// This could include showing additional information or wider layouts
}

// createFilePathsTab creates the enhanced file paths configuration tab
func (a *App) createFilePathsTab() *fyne.Container {
	// Master file selection with enhanced styling
	a.masterFileEntry = widget.NewEntry()
	a.masterFileEntry.SetPlaceHolder("Select master Excel file (.xlsx, .xls)...")
	masterFileButton := widget.NewButton("Browse", func() {
		a.selectMasterFile()
	})
	masterFileButton.Importance = widget.MediumImportance

	// Student folder selection
	a.studentFolderEntry = widget.NewEntry()
	a.studentFolderEntry.SetPlaceHolder("Select student files folder (recursive scan)...")
	studentFolderButton := widget.NewButton("Browse", func() {
		a.selectStudentFolder()
	})
	studentFolderButton.Importance = widget.MediumImportance

	// Output folder selection
	a.outputFolderEntry = widget.NewEntry()
	a.outputFolderEntry.SetPlaceHolder("Select output folder for processed files...")
	outputFolderButton := widget.NewButton("Browse", func() {
		a.selectOutputFolder()
	})
	outputFolderButton.Importance = widget.MediumImportance

	// Backup folder selection
	a.backupFolderEntry = widget.NewEntry()
	a.backupFolderEntry.SetPlaceHolder("Select backup folder for safety copies...")
	backupFolderButton := widget.NewButton("Browse", func() {
		a.selectBackupFolder()
	})
	backupFolderButton.Importance = widget.MediumImportance

	// Enhanced layout with better spacing and visual hierarchy
	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Master Excel File *:", Widget: container.NewBorder(nil, nil, nil, masterFileButton, a.masterFileEntry)},
			{Text: "Student Files Folder *:", Widget: container.NewBorder(nil, nil, nil, studentFolderButton, a.studentFolderEntry)},
			{Text: "Output Folder:", Widget: container.NewBorder(nil, nil, nil, outputFolderButton, a.outputFolderEntry)},
			{Text: "Backup Folder:", Widget: container.NewBorder(nil, nil, nil, backupFolderButton, a.backupFolderEntry)},
		},
	}

	// Add help text
	helpText := widget.NewLabel("* Required fields. Configure the file paths for processing student marks.")
	helpText.TextStyle = fyne.TextStyle{Italic: true}

	return container.NewVBox(
		widget.NewCard("File and Folder Configuration",
			"Configure input and output locations for mark processing",
			container.NewVBox(form, widget.NewSeparator(), helpText)),
	)
}

// createExcelSettingsTab creates the enhanced Excel configuration tab
func (a *App) createExcelSettingsTab() *fyne.Container {
	// Excel settings entries with enhanced validation
	a.studentWorksheetEntry = widget.NewEntry()
	a.studentWorksheetEntry.SetText("Grading Sheet")
	a.studentWorksheetEntry.SetPlaceHolder("Name of worksheet in student files")

	a.masterWorksheetEntry = widget.NewEntry()
	a.masterWorksheetEntry.SetText("001")
	a.masterWorksheetEntry.SetPlaceHolder("Name of worksheet in master file")

	a.studentIDCellEntry = widget.NewEntry()
	a.studentIDCellEntry.SetText("B2")
	a.studentIDCellEntry.SetPlaceHolder("e.g., B2, C3, A1")

	a.studentIDColumnEntry = widget.NewEntry()
	a.studentIDColumnEntry.SetText("B")
	a.studentIDColumnEntry.SetPlaceHolder("e.g., A, B, C")

	// Enhanced validation with visual feedback
	a.studentIDCellEntry.OnChanged = func(text string) {
		a.validateCellReference(text, "Student ID Cell")
	}

	a.studentIDColumnEntry.OnChanged = func(text string) {
		a.validateColumnReference(text, "Student ID Column")
	}

	// Enhanced form with better descriptions
	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Student Worksheet Name:", Widget: a.studentWorksheetEntry},
			{Text: "Master Worksheet Name:", Widget: a.masterWorksheetEntry},
			{Text: "Student ID Cell Location:", Widget: a.studentIDCellEntry},
			{Text: "Student ID Column (Master):", Widget: a.studentIDColumnEntry},
		},
	}

	// Add configuration examples
	exampleText := widget.NewLabel("Examples: Cell references like B2, C3 | Column references like A, B, AA")
	exampleText.TextStyle = fyne.TextStyle{Italic: true}

	return container.NewVBox(
		widget.NewCard("Excel Worksheet Configuration",
			"Configure worksheet names and cell locations for data extraction",
			container.NewVBox(form, widget.NewSeparator(), exampleText)),
	)
}

// createStatusBar creates the enhanced bottom status bar
func (a *App) createStatusBar() *fyne.Container {
	a.progressBar = widget.NewProgressBar()
	a.progressBar.Hide()

	a.statusLabel = widget.NewLabel("Ready")
	a.statusLabel.TextStyle = fyne.TextStyle{Bold: true}

	// Add version and author info to status bar
	versionLabel := widget.NewLabel("v1.0.0 | © Vinay Koirala")
	versionLabel.TextStyle = fyne.TextStyle{Italic: true}

	return container.NewBorder(
		nil, nil,
		a.statusLabel,
		versionLabel,
		container.NewPadded(a.progressBar),
	)
}

// validateCellReference validates Excel cell reference format
func (a *App) validateCellReference(cellRef, fieldName string) {
	if cellRef == "" {
		return
	}
	
	// Basic validation for Excel cell reference (e.g., A1, B2, AA10)
	valid := true
	if len(cellRef) < 2 {
		valid = false
	} else {
		// Check if it starts with letters and ends with numbers
		i := 0
		for i < len(cellRef) && cellRef[i] >= 'A' && cellRef[i] <= 'Z' {
			i++
		}
		if i == 0 || i == len(cellRef) {
			valid = false
		} else {
			for j := i; j < len(cellRef); j++ {
				if cellRef[j] < '0' || cellRef[j] > '9' {
					valid = false
					break
				}
			}
		}
	}
	
	if !valid {
		a.showError(fmt.Sprintf("Invalid cell reference format for %s: %s\nExpected format: A1, B2, AA10, etc.", fieldName, cellRef))
	}
}

// showError displays an error dialog
func (a *App) showError(message string) {
	dialog.ShowError(fmt.Errorf(message), a.window)
}

// showInfo displays an info dialog
func (a *App) showInfo(title, message string) {
	dialog.ShowInformation(title, message, a.window)
}

// updateStatus updates the status label with text indicators and optionally the progress bar
func (a *App) updateStatus(status string, progress ...float64) {
	// Add status indicators based on content
	var statusPrefix string
	switch {
	case strings.Contains(strings.ToLower(status), "error") || strings.Contains(strings.ToLower(status), "failed"):
		statusPrefix = "[ERROR]"
	case strings.Contains(strings.ToLower(status), "warning"):
		statusPrefix = "[WARNING]"
	case strings.Contains(strings.ToLower(status), "processing") || strings.Contains(strings.ToLower(status), "building"):
		statusPrefix = "[PROCESSING]"
	case strings.Contains(strings.ToLower(status), "completed") || strings.Contains(strings.ToLower(status), "success"):
		statusPrefix = "[SUCCESS]"
	default:
		statusPrefix = "[READY]"
	}

	a.statusLabel.SetText(fmt.Sprintf("%s %s", statusPrefix, status))

	if len(progress) > 0 {
		a.progressBar.SetValue(progress[0])
		if progress[0] > 0 {
			a.progressBar.Show()
		} else {
			a.progressBar.Hide()
		}
	}
}

// createMarkMappingsTab creates the enhanced mark mappings configuration tab
func (a *App) createMarkMappingsTab() *fyne.Container {
	// Create header labels for the mapping table
	headerContainer := container.NewHBox(
		widget.NewLabel("Student Cell"),
		widget.NewLabel(""),  // spacer
		widget.NewLabel("Master Column"),
		widget.NewLabel(""),  // spacer for remove button
	)

	// Create enhanced list for mark mappings with better visibility
	a.markMappingTable = widget.NewList(
		func() int {
			return len(a.markMappings)
		},
		func() fyne.CanvasObject {
			// Create entry fields with better sizing
			studentCell := widget.NewEntry()
			studentCell.SetPlaceHolder("C6")
			studentCell.Resize(fyne.NewSize(80, 32))

			masterColumn := widget.NewEntry()
			masterColumn.SetPlaceHolder("I")
			masterColumn.Resize(fyne.NewSize(80, 32))

			// Create arrow label
			arrowLabel := widget.NewLabel("->")
			arrowLabel.Alignment = fyne.TextAlignCenter

			// Create remove button
			removeBtn := widget.NewButton("Remove", nil)
			removeBtn.Importance = widget.DangerImportance
			removeBtn.Resize(fyne.NewSize(80, 32))

			// Create container with proper spacing and alignment
			return container.NewBorder(
				nil, nil, nil, nil,
				container.NewHBox(
					container.NewBorder(nil, nil, nil, nil, studentCell),
					container.NewBorder(nil, nil, nil, nil, arrowLabel),
					container.NewBorder(nil, nil, nil, nil, masterColumn),
					container.NewBorder(nil, nil, nil, nil, removeBtn),
				),
			)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			if id >= len(a.markMappings) {
				return
			}

			// Navigate to the actual container with entries
			borderContainer := obj.(*fyne.Container)
			hboxContainer := borderContainer.Objects[0].(*fyne.Container)

			studentCellContainer := hboxContainer.Objects[0].(*fyne.Container)
			studentCell := studentCellContainer.Objects[0].(*widget.Entry)

			masterColumnContainer := hboxContainer.Objects[2].(*fyne.Container)
			masterColumn := masterColumnContainer.Objects[0].(*widget.Entry)

			removeBtnContainer := hboxContainer.Objects[3].(*fyne.Container)
			removeBtn := removeBtnContainer.Objects[0].(*widget.Button)

			mapping := a.markMappings[id]
			studentCell.SetText(mapping.StudentCell)
			masterColumn.SetText(mapping.MasterColumn)

			// Update mapping when entries change with enhanced validation
			studentCell.OnChanged = func(text string) {
				if id < len(a.markMappings) {
					a.markMappings[id].StudentCell = text
					a.validateCellReference(text, "Student Cell")
				}
			}

			masterColumn.OnChanged = func(text string) {
				if id < len(a.markMappings) {
					a.markMappings[id].MasterColumn = text
					a.validateColumnReference(text, "Master Column")
				}
			}

			removeBtn.OnTapped = func() {
				a.removeMarkMapping(id)
			}
		},
	)

	// Enhanced buttons for managing mappings
	addButton := widget.NewButton("Add Mapping", func() {
		a.addMarkMapping()
	})
	addButton.Importance = widget.HighImportance

	resetButton := widget.NewButton("Reset to Default", func() {
		a.resetMarkMappings()
	})
	resetButton.Importance = widget.MediumImportance

	buttonContainer := container.NewHBox(addButton, resetButton)

	// Create scrollable container for the mappings
	scrollContainer := container.NewScroll(a.markMappingTable)
	scrollContainer.SetMinSize(fyne.NewSize(600, 300))

	// Add mapping count info
	mappingInfo := widget.NewLabel(fmt.Sprintf("Current mappings: %d", len(a.markMappings)))
	mappingInfo.TextStyle = fyne.TextStyle{Italic: true}

	// Main container with proper layout
	content := container.NewVBox(
		headerContainer,
		widget.NewSeparator(),
		scrollContainer,
		widget.NewSeparator(),
		mappingInfo,
		buttonContainer,
	)

	return container.NewVBox(
		widget.NewCard("Mark Cell Mappings",
			"Configure how student marks map to master sheet columns (Student Cell -> Master Column)",
			content),
	)
}

// createProcessingTab creates the enhanced processing configuration and control tab
func (a *App) createProcessingTab() *fyne.Container {
	// Enhanced processing options with better descriptions
	a.enableBackupCheck = widget.NewCheck("Enable Backup (Recommended)", nil)
	a.enableBackupCheck.SetChecked(true)

	a.skipInvalidCheck = widget.NewCheck("Skip Invalid Files (Continue on errors)", nil)
	a.skipInvalidCheck.SetChecked(true)

	a.maxConcurrentEntry = widget.NewEntry()
	a.maxConcurrentEntry.SetText("10")
	a.maxConcurrentEntry.SetPlaceHolder("1-20")
	a.maxConcurrentEntry.Validator = func(text string) error {
		if val, err := strconv.Atoi(text); err != nil || val < 1 || val > 20 {
			return fmt.Errorf("must be a number between 1 and 20")
		}
		return nil
	}

	// Enhanced processing buttons
	dryRunButton := widget.NewButton("Dry Run (Test)", func() {
		a.startProcessing(true)
	})
	dryRunButton.Importance = widget.MediumImportance

	processButton := widget.NewButton("Process Files", func() {
		a.startProcessing(false)
	})
	processButton.Importance = widget.HighImportance

	stopButton := widget.NewButton("Stop", func() {
		a.stopProcessing()
	})
	stopButton.Importance = widget.DangerImportance
	stopButton.Disable()

	// Enhanced configuration buttons
	loadConfigButton := widget.NewButton("Load Config", func() {
		a.loadConfigFromFile()
	})
	loadConfigButton.Importance = widget.MediumImportance

	saveConfigButton := widget.NewButton("Save Config", func() {
		a.saveConfigToFile()
	})
	saveConfigButton.Importance = widget.MediumImportance

	// Enhanced layout with better organization
	optionsForm := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Backup Options:", Widget: a.enableBackupCheck},
			{Text: "Error Handling:", Widget: a.skipInvalidCheck},
			{Text: "Concurrent Processing:", Widget: a.maxConcurrentEntry},
		},
	}

	// Add help text for options
	optionsHelp := widget.NewLabel("Configure how the application handles processing and errors")
	optionsHelp.TextStyle = fyne.TextStyle{Italic: true}

	processingButtons := container.NewHBox(dryRunButton, processButton, stopButton)
	configButtons := container.NewHBox(loadConfigButton, saveConfigButton)

	// Processing controls section
	processingSection := container.NewVBox(
		widget.NewLabel("Processing Controls:"),
		processingButtons,
		widget.NewLabel("Use 'Dry Run' to test your configuration before processing actual files."),
	)

	// Configuration management section
	configSection := container.NewVBox(
		widget.NewLabel("Configuration Management:"),
		configButtons,
		widget.NewLabel("Save your current settings or load a previously saved configuration."),
	)

	return container.NewVBox(
		widget.NewCard("Processing Options",
			"Configure processing behavior and safety options",
			container.NewVBox(optionsForm, widget.NewSeparator(), optionsHelp)),
		widget.NewCard("Actions & Controls",
			"Start processing or manage configuration files",
			container.NewVBox(
				processingSection,
				widget.NewSeparator(),
				configSection,
			)),
	)
}

// createLogsTab creates the enhanced logs and output tab
func (a *App) createLogsTab() *fyne.Container {
	a.logOutput = widget.NewMultiLineEntry()
	a.logOutput.SetText("Application ready. Configure settings and start processing.\n")
	a.logOutput.Wrapping = fyne.TextWrapWord
	a.logOutput.MultiLine = true

	// Enhanced buttons for log management
	clearButton := widget.NewButton("Clear Logs", func() {
		a.logOutput.SetText("Application ready. Configure settings and start processing.\n")
	})
	clearButton.Importance = widget.MediumImportance

	exportButton := widget.NewButton("Export Logs", func() {
		a.exportLogs()
	})
	exportButton.Importance = widget.MediumImportance

	// Create button container
	buttonContainer := container.NewHBox(clearButton, exportButton)

	// Create scrollable log container with better sizing
	logScroll := container.NewScroll(a.logOutput)
	logScroll.SetMinSize(fyne.NewSize(800, 400))

	// Add log statistics
	logStats := widget.NewLabel("Ready to process files...")
	logStats.TextStyle = fyne.TextStyle{Italic: true}

	return container.NewVBox(
		widget.NewCard("Processing Logs",
			"Real-time processing status, results, and error messages",
			container.NewVBox(
				logScroll,
				widget.NewSeparator(),
				logStats,
				buttonContainer,
			)),
	)
}

// addMarkMapping adds a new mark mapping row
func (a *App) addMarkMapping() {
	a.markMappings = append(a.markMappings, MarkMapping{
		StudentCell:  "",
		MasterColumn: "",
	})
	a.markMappingTable.Refresh()
	a.updateStatus(fmt.Sprintf("Added new mapping. Total: %d mappings", len(a.markMappings)))
}

// exportLogs exports the current log content to a file
func (a *App) exportLogs() {
	dialog.ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err != nil || writer == nil {
			return
		}
		defer writer.Close()

		logContent := a.logOutput.Text
		if _, err := writer.Write([]byte(logContent)); err != nil {
			a.showError(fmt.Sprintf("Failed to export logs: %v", err))
			return
		}

		a.updateStatus("Logs exported successfully")
	}, a.window)
}

// removeMarkMapping removes a mark mapping at the specified index
func (a *App) removeMarkMapping(index int) {
	if index >= 0 && index < len(a.markMappings) {
		a.markMappings = append(a.markMappings[:index], a.markMappings[index+1:]...)
		a.markMappingTable.Refresh()
		a.updateStatus(fmt.Sprintf("Removed mapping. Total: %d mappings", len(a.markMappings)))
	}
}

// resetMarkMappings resets mark mappings to default values
func (a *App) resetMarkMappings() {
	a.markMappings = getDefaultMarkMappings()
	a.markMappingTable.Refresh()
	a.updateStatus(fmt.Sprintf("Reset to default mappings. Total: %d mappings", len(a.markMappings)))
}

// validateColumnReference validates Excel column reference format
func (a *App) validateColumnReference(colRef, fieldName string) {
	if colRef == "" {
		return
	}

	// Basic validation for Excel column reference (e.g., A, B, AA, AB)
	valid := true
	if len(colRef) == 0 {
		valid = false
	} else {
		for _, char := range colRef {
			if char < 'A' || char > 'Z' {
				valid = false
				break
			}
		}
	}

	if !valid {
		a.showError(fmt.Sprintf("Invalid column reference format for %s: %s\nExpected format: A, B, AA, AB, etc.", fieldName, colRef))
	}
}

// File selection methods
func (a *App) selectMasterFile() {
	dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil || reader == nil {
			return
		}
		defer reader.Close()

		path := reader.URI().Path()
		a.masterFileEntry.SetText(path)
	}, a.window)
}

func (a *App) selectStudentFolder() {
	dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
		if err != nil || uri == nil {
			return
		}

		path := uri.Path()
		a.studentFolderEntry.SetText(path)
	}, a.window)
}

func (a *App) selectOutputFolder() {
	dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
		if err != nil || uri == nil {
			return
		}

		path := uri.Path()
		a.outputFolderEntry.SetText(path)
	}, a.window)
}

func (a *App) selectBackupFolder() {
	dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
		if err != nil || uri == nil {
			return
		}

		path := uri.Path()
		a.backupFolderEntry.SetText(path)
	}, a.window)
}

// setupMenus creates the application menu bar
func (a *App) setupMenus() {
	// File menu
	newItem := fyne.NewMenuItem("New Configuration", func() {
		a.resetToDefaults()
	})

	loadItem := fyne.NewMenuItem("Load Configuration...", func() {
		a.loadConfigFromFile()
	})

	saveItem := fyne.NewMenuItem("Save Configuration...", func() {
		a.saveConfigToFile()
	})

	quitItem := fyne.NewMenuItem("Quit", func() {
		a.fyneApp.Quit()
	})

	fileMenu := fyne.NewMenu("File", newItem, fyne.NewMenuItemSeparator(), loadItem, saveItem, fyne.NewMenuItemSeparator(), quitItem)

	// Edit menu
	resetItem := fyne.NewMenuItem("Reset Mark Mappings", func() {
		a.resetMarkMappings()
	})

	editMenu := fyne.NewMenu("Edit", resetItem)

	// Help menu
	aboutItem := fyne.NewMenuItem("About", func() {
		a.showAbout()
	})

	helpItem := fyne.NewMenuItem("Help", func() {
		a.showHelp()
	})

	helpMenu := fyne.NewMenu("Help", helpItem, aboutItem)

	// Set main menu
	mainMenu := fyne.NewMainMenu(fileMenu, editMenu, helpMenu)
	a.window.SetMainMenu(mainMenu)
}

// resetToDefaults resets all configuration to default values
func (a *App) resetToDefaults() {
	a.masterFileEntry.SetText("")
	a.studentFolderEntry.SetText("")
	a.outputFolderEntry.SetText("./output")
	a.backupFolderEntry.SetText("./backups")

	a.studentWorksheetEntry.SetText("Grading Sheet")
	a.masterWorksheetEntry.SetText("001")
	a.studentIDCellEntry.SetText("B2")
	a.studentIDColumnEntry.SetText("B")

	a.enableBackupCheck.SetChecked(true)
	a.skipInvalidCheck.SetChecked(true)
	a.maxConcurrentEntry.SetText("10")

	a.resetMarkMappings()

	a.updateStatus("Configuration reset to defaults")
}

// showAbout displays the about dialog
func (a *App) showAbout() {
	content := fmt.Sprintf(`Mark Master Sheet Consolidator v1.0.0

A production-ready application for consolidating student marks from individual Excel files into a master spreadsheet.

Author: Vinay Koirala
Email: koiralavinay@gmail.com
Professional: binaya.koirala@iic.edu.np
LinkedIn: veenay
GitHub: v-eenay

Repository: https://github.com/v-eenay/MarkMasterSheetConsolidator.git

All rights reserved to Vinay Koirala`)

	dialog.ShowInformation("About", content, a.window)
}

// showHelp displays the help dialog
func (a *App) showHelp() {
	content := `Mark Master Sheet Consolidator Help

Quick Start:
1. Configure file paths in the "File Paths" tab
2. Set Excel worksheet names in "Excel Settings" tab
3. Configure mark cell mappings in "Mark Mappings" tab
4. Set processing options in "Processing" tab
5. Use "Dry Run" to test, then "Process Files" to execute

File Paths:
- Master Excel File: The main spreadsheet to update
- Student Files Folder: Folder containing student Excel files (scanned recursively)
- Output Folder: Where updated master sheets will be saved
- Backup Folder: Where backup copies will be stored

Excel Settings:
- Student Worksheet Name: Name of worksheet in student files (default: "Grading Sheet")
- Master Worksheet Name: Name of worksheet in master file (default: "001")
- Student ID Cell: Cell containing student ID in student files (default: "B2")
- Student ID Column: Column containing student IDs in master sheet (default: "B")

Mark Mappings:
- Configure which cells in student files map to which columns in master sheet
- Use "Add Mapping" to add new mappings
- Use "Remove" button to delete mappings
- Use "Reset to Default" to restore standard mappings

Processing Options:
- Enable Backup: Creates timestamped backups before changes
- Skip Invalid Files: Continues processing if some files fail
- Max Concurrent Files: Number of files to process simultaneously (1-20)

For more detailed information, see the documentation at:
https://github.com/v-eenay/MarkMasterSheetConsolidator`

	dialog.ShowInformation("Help", content, a.window)
}
