package gui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// ModernLightTheme provides a custom modern light theme for the application
type ModernLightTheme struct{}

// Color returns theme colors with a modern professional light palette
func (m ModernLightTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	// Force light variant for all colors to ensure consistent light theme
	switch name {
	case theme.ColorNamePrimary:
		return color.RGBA{R: 25, G: 118, B: 210, A: 255} // Professional blue
	case theme.ColorNameBackground:
		return color.RGBA{R: 255, G: 255, B: 255, A: 255} // Pure white background
	case theme.ColorNameButton:
		return color.RGBA{R: 248, G: 249, B: 250, A: 255} // Very light gray button
	case theme.ColorNameDisabledButton:
		return color.RGBA{R: 233, G: 236, B: 239, A: 255} // Light disabled
	case theme.ColorNameForeground:
		return color.RGBA{R: 33, G: 37, B: 41, A: 255} // Dark text for readability
	case theme.ColorNameSuccess:
		return color.RGBA{R: 40, G: 167, B: 69, A: 255} // Green
	case theme.ColorNameWarning:
		return color.RGBA{R: 255, G: 193, B: 7, A: 255} // Amber
	case theme.ColorNameError:
		return color.RGBA{R: 220, G: 53, B: 69, A: 255} // Red
	case theme.ColorNameInputBackground:
		return color.RGBA{R: 255, G: 255, B: 255, A: 255} // White input background
	case theme.ColorNameInputBorder:
		return color.RGBA{R: 206, G: 212, B: 218, A: 255} // Light border
	case theme.ColorNameScrollBar:
		return color.RGBA{R: 206, G: 212, B: 218, A: 255} // Light scrollbar
	case theme.ColorNameShadow:
		return color.RGBA{R: 0, G: 0, B: 0, A: 25} // Subtle shadow
	}

	// Fall back to default theme for other colors but force light variant
	return theme.DefaultTheme().Color(name, theme.VariantLight)
}

// Font returns theme fonts with modern typography
func (m ModernLightTheme) Font(style fyne.TextStyle) fyne.Resource {
	// Use default fonts which handle Unicode characters properly
	return theme.DefaultTheme().Font(style)
}

// Icon returns theme icons
func (m ModernLightTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

// Size returns theme sizes with modern spacing
func (m ModernLightTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNamePadding:
		return 12 // Generous padding for modern look
	case theme.SizeNameInlineIcon:
		return 20 // Properly sized icons
	case theme.SizeNameScrollBar:
		return 16 // Comfortable scroll bars
	case theme.SizeNameSeparatorThickness:
		return 1 // Clean thin separators
	case theme.SizeNameInputBorder:
		return 2 // Visible input borders
	case theme.SizeNameText:
		return 14 // Readable text size
	}

	// Fall back to default theme for other sizes
	return theme.DefaultTheme().Size(name)
}

// applyCustomTheme applies the custom light theme to the application
func (a *App) applyCustomTheme() {
	// Apply the modern light theme to ensure good visibility
	a.fyneApp.Settings().SetTheme(&ModernLightTheme{})
}

// Modern color constants for consistent styling
var (
	ModernBlue   = color.RGBA{R: 33, G: 150, B: 243, A: 255}
	ModernGreen  = color.RGBA{R: 76, G: 175, B: 80, A: 255}
	ModernOrange = color.RGBA{R: 255, G: 152, B: 0, A: 255}
	ModernRed    = color.RGBA{R: 244, G: 67, B: 54, A: 255}
	ModernGray   = color.RGBA{R: 158, G: 158, B: 158, A: 255}
	LightGray    = color.RGBA{R: 245, G: 245, B: 245, A: 255}
)

// createStyledButton creates a button with enhanced styling
func createStyledButton(text string, icon string, importance widget.Importance, onTapped func()) *widget.Button {
	buttonText := text
	if icon != "" {
		buttonText = icon + " " + text
	}
	
	button := widget.NewButton(buttonText, onTapped)
	button.Importance = importance
	
	return button
}

// createStyledCard creates a card with enhanced styling and icons
func createStyledCard(title, subtitle string, content fyne.CanvasObject) *widget.Card {
	card := widget.NewCard(title, subtitle, content)
	return card
}

// createInfoLabel creates a styled info label
func createInfoLabel(text string, style fyne.TextStyle) *widget.Label {
	label := widget.NewLabel(text)
	label.TextStyle = style
	return label
}

// createSectionHeader creates a styled section header
func createSectionHeader(text string) *widget.Label {
	header := widget.NewLabel(text)
	header.TextStyle = fyne.TextStyle{Bold: true}
	return header
}

// createHelpText creates styled help text
func createHelpText(text string) *widget.Label {
	help := widget.NewLabel(text)
	help.TextStyle = fyne.TextStyle{Italic: true}
	return help
}

// ResponsiveContainer creates a container that adapts to screen size
func ResponsiveContainer(objects ...fyne.CanvasObject) *fyne.Container {
	// For now, use VBox but could be enhanced with responsive logic
	return container.NewVBox(objects...)
}

// ModernSpacing provides consistent spacing values
const (
	SmallSpacing  = 4
	MediumSpacing = 8
	LargeSpacing  = 16
	XLargeSpacing = 24
)

// WindowConstraints defines responsive window sizing
type WindowConstraints struct {
	MinWidth  float32
	MinHeight float32
	MaxWidth  float32
	MaxHeight float32
	OptWidth  float32
	OptHeight float32
}

// GetWindowConstraints returns the recommended window constraints
func GetWindowConstraints() WindowConstraints {
	return WindowConstraints{
		MinWidth:  800,
		MinHeight: 600,
		MaxWidth:  1920,
		MaxHeight: 1080,
		OptWidth:  1200,
		OptHeight: 800,
	}
}
