// Package gui provides the graphical user interface for the Mark Master Sheet Consolidator.
// It uses the Fyne framework to create a user-friendly interface for configuration and processing.
package gui

import (
	"context"
	"fmt"
	"strconv"

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

// NewApp creates a new GUI application instance
func NewApp() *App {
	fyneApp := app.NewWithID("com.vinaykoirala.markmaster")

	window := fyneApp.NewWindow("Mark Master Sheet Consolidator")
	window.Resize(fyne.NewSize(900, 700))
	window.SetMaster()

	return &App{
		fyneApp: fyneApp,
		window:  window,
		markMappings: getDefaultMarkMappings(),
	}
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

// setupUI creates and arranges all UI components
func (a *App) setupUI() {
	// Create main container with tabs
	tabs := container.NewAppTabs(
		container.NewTabItem("File Paths", a.createFilePathsTab()),
		container.NewTabItem("Excel Settings", a.createExcelSettingsTab()),
		container.NewTabItem("Mark Mappings", a.createMarkMappingsTab()),
		container.NewTabItem("Processing", a.createProcessingTab()),
		container.NewTabItem("Logs", a.createLogsTab()),
	)

	// Create status bar
	statusBar := a.createStatusBar()

	// Main layout
	content := container.NewBorder(
		nil,           // top
		statusBar,     // bottom
		nil,           // left
		nil,           // right
		tabs,          // center
	)

	a.window.SetContent(content)
}

// createFilePathsTab creates the file paths configuration tab
func (a *App) createFilePathsTab() *fyne.Container {
	// Master file selection
	a.masterFileEntry = widget.NewEntry()
	a.masterFileEntry.SetPlaceHolder("Select master Excel file...")
	masterFileButton := widget.NewButton("Browse", func() {
		a.selectMasterFile()
	})

	// Student folder selection
	a.studentFolderEntry = widget.NewEntry()
	a.studentFolderEntry.SetPlaceHolder("Select student files folder...")
	studentFolderButton := widget.NewButton("Browse", func() {
		a.selectStudentFolder()
	})

	// Output folder selection
	a.outputFolderEntry = widget.NewEntry()
	a.outputFolderEntry.SetPlaceHolder("Select output folder...")
	outputFolderButton := widget.NewButton("Browse", func() {
		a.selectOutputFolder()
	})

	// Backup folder selection
	a.backupFolderEntry = widget.NewEntry()
	a.backupFolderEntry.SetPlaceHolder("Select backup folder...")
	backupFolderButton := widget.NewButton("Browse", func() {
		a.selectBackupFolder()
	})

	// Layout
	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Master Excel File:", Widget: container.NewBorder(nil, nil, nil, masterFileButton, a.masterFileEntry)},
			{Text: "Student Files Folder:", Widget: container.NewBorder(nil, nil, nil, studentFolderButton, a.studentFolderEntry)},
			{Text: "Output Folder:", Widget: container.NewBorder(nil, nil, nil, outputFolderButton, a.outputFolderEntry)},
			{Text: "Backup Folder:", Widget: container.NewBorder(nil, nil, nil, backupFolderButton, a.backupFolderEntry)},
		},
	}

	return container.NewVBox(
		widget.NewCard("File and Folder Paths", "Configure input and output locations", form),
	)
}

// createExcelSettingsTab creates the Excel configuration tab
func (a *App) createExcelSettingsTab() *fyne.Container {
	// Excel settings entries
	a.studentWorksheetEntry = widget.NewEntry()
	a.studentWorksheetEntry.SetText("Grading Sheet")
	
	a.masterWorksheetEntry = widget.NewEntry()
	a.masterWorksheetEntry.SetText("001")
	
	a.studentIDCellEntry = widget.NewEntry()
	a.studentIDCellEntry.SetText("B2")
	
	a.studentIDColumnEntry = widget.NewEntry()
	a.studentIDColumnEntry.SetText("B")

	// Validation
	a.studentIDCellEntry.OnChanged = func(text string) {
		a.validateCellReference(text, "Student ID Cell")
	}

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Student Worksheet Name:", Widget: a.studentWorksheetEntry},
			{Text: "Master Worksheet Name:", Widget: a.masterWorksheetEntry},
			{Text: "Student ID Cell (in student files):", Widget: a.studentIDCellEntry},
			{Text: "Student ID Column (in master sheet):", Widget: a.studentIDColumnEntry},
		},
	}

	return container.NewVBox(
		widget.NewCard("Excel Configuration", "Configure worksheet names and cell locations", form),
	)
}

// createStatusBar creates the bottom status bar
func (a *App) createStatusBar() *fyne.Container {
	a.progressBar = widget.NewProgressBar()
	a.progressBar.Hide()
	
	a.statusLabel = widget.NewLabel("Ready")
	
	return container.NewBorder(
		nil, nil, a.statusLabel, nil,
		a.progressBar,
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

// updateStatus updates the status label and optionally the progress bar
func (a *App) updateStatus(status string, progress ...float64) {
	a.statusLabel.SetText(status)
	if len(progress) > 0 {
		a.progressBar.SetValue(progress[0])
		if progress[0] > 0 {
			a.progressBar.Show()
		} else {
			a.progressBar.Hide()
		}
	}
}

// createMarkMappingsTab creates the mark mappings configuration tab
func (a *App) createMarkMappingsTab() *fyne.Container {
	// Create list for mark mappings
	a.markMappingTable = widget.NewList(
		func() int {
			return len(a.markMappings)
		},
		func() fyne.CanvasObject {
			studentCell := widget.NewEntry()
			studentCell.SetPlaceHolder("C6")
			masterColumn := widget.NewEntry()
			masterColumn.SetPlaceHolder("I")
			removeBtn := widget.NewButton("Remove", nil)

			return container.NewHBox(
				widget.NewLabel("Student Cell:"),
				studentCell,
				widget.NewLabel("→ Master Column:"),
				masterColumn,
				removeBtn,
			)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			if id >= len(a.markMappings) {
				return
			}

			containerObj := obj.(*fyne.Container)
			studentCell := containerObj.Objects[1].(*widget.Entry)
			masterColumn := containerObj.Objects[3].(*widget.Entry)
			removeBtn := containerObj.Objects[4].(*widget.Button)

			mapping := a.markMappings[id]
			studentCell.SetText(mapping.StudentCell)
			masterColumn.SetText(mapping.MasterColumn)

			// Update mapping when entries change
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

	// Buttons for managing mappings
	addButton := widget.NewButton("Add Mapping", func() {
		a.addMarkMapping()
	})

	resetButton := widget.NewButton("Reset to Default", func() {
		a.resetMarkMappings()
	})

	buttonContainer := container.NewHBox(addButton, resetButton)

	return container.NewVBox(
		widget.NewCard("Mark Cell Mappings", "Configure how student marks map to master sheet columns",
			container.NewBorder(nil, buttonContainer, nil, nil, a.markMappingTable)),
	)
}

// createProcessingTab creates the processing configuration and control tab
func (a *App) createProcessingTab() *fyne.Container {
	// Processing options
	a.enableBackupCheck = widget.NewCheck("Enable Backup", nil)
	a.enableBackupCheck.SetChecked(true)

	a.skipInvalidCheck = widget.NewCheck("Skip Invalid Files", nil)
	a.skipInvalidCheck.SetChecked(true)

	a.maxConcurrentEntry = widget.NewEntry()
	a.maxConcurrentEntry.SetText("10")
	a.maxConcurrentEntry.Validator = func(text string) error {
		if val, err := strconv.Atoi(text); err != nil || val < 1 || val > 20 {
			return fmt.Errorf("must be a number between 1 and 20")
		}
		return nil
	}

	// Processing buttons
	dryRunButton := widget.NewButton("Dry Run", func() {
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

	// Configuration buttons
	loadConfigButton := widget.NewButton("Load Config", func() {
		a.loadConfigFromFile()
	})

	saveConfigButton := widget.NewButton("Save Config", func() {
		a.saveConfigToFile()
	})

	// Layout
	optionsForm := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Backup Options:", Widget: a.enableBackupCheck},
			{Text: "Error Handling:", Widget: a.skipInvalidCheck},
			{Text: "Max Concurrent Files:", Widget: a.maxConcurrentEntry},
		},
	}

	processingButtons := container.NewHBox(dryRunButton, processButton, stopButton)
	configButtons := container.NewHBox(loadConfigButton, saveConfigButton)

	return container.NewVBox(
		widget.NewCard("Processing Options", "Configure processing behavior", optionsForm),
		widget.NewCard("Actions", "Start processing or manage configuration",
			container.NewVBox(
				widget.NewLabel("Processing Controls:"),
				processingButtons,
				widget.NewSeparator(),
				widget.NewLabel("Configuration Management:"),
				configButtons,
			)),
	)
}

// createLogsTab creates the logs and output tab
func (a *App) createLogsTab() *fyne.Container {
	a.logOutput = widget.NewMultiLineEntry()
	a.logOutput.SetText("Application ready. Configure settings and start processing.\n")
	a.logOutput.Wrapping = fyne.TextWrapWord

	clearButton := widget.NewButton("Clear Logs", func() {
		a.logOutput.SetText("")
	})

	return container.NewVBox(
		widget.NewCard("Processing Logs", "Real-time processing status and results",
			container.NewBorder(nil, clearButton, nil, nil,
				container.NewScroll(a.logOutput))),
	)
}

// addMarkMapping adds a new mark mapping row
func (a *App) addMarkMapping() {
	a.markMappings = append(a.markMappings, MarkMapping{
		StudentCell:  "",
		MasterColumn: "",
	})
	a.markMappingTable.Refresh()
}

// removeMarkMapping removes a mark mapping at the specified index
func (a *App) removeMarkMapping(index int) {
	if index >= 0 && index < len(a.markMappings) {
		a.markMappings = append(a.markMappings[:index], a.markMappings[index+1:]...)
		a.markMappingTable.Refresh()
	}
}

// resetMarkMappings resets mark mappings to default values
func (a *App) resetMarkMappings() {
	a.markMappings = getDefaultMarkMappings()
	a.markMappingTable.Refresh()
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
