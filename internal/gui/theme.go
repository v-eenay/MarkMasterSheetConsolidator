package gui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// ModernTheme provides a custom modern theme for the application
type ModernTheme struct{}

// Color returns theme colors with a modern professional palette
func (m ModernTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNamePrimary:
		if variant == theme.VariantLight {
			return color.RGBA{R: 33, G: 150, B: 243, A: 255} // Modern blue
		}
		return color.RGBA{R: 100, G: 181, B: 246, A: 255} // Light blue for dark theme
	case theme.ColorNameBackground:
		if variant == theme.VariantLight {
			return color.RGBA{R: 250, G: 250, B: 250, A: 255} // Very light gray
		}
		return color.RGBA{R: 33, G: 33, B: 33, A: 255} // Dark gray
	case theme.ColorNameButton:
		if variant == theme.VariantLight {
			return color.RGBA{R: 245, G: 245, B: 245, A: 255} // Light button
		}
		return color.RGBA{R: 66, G: 66, B: 66, A: 255} // Dark button
	case theme.ColorNameDisabledButton:
		if variant == theme.VariantLight {
			return color.RGBA{R: 224, G: 224, B: 224, A: 255} // Disabled light
		}
		return color.RGBA{R: 97, G: 97, B: 97, A: 255} // Disabled dark
	case theme.ColorNameForeground:
		if variant == theme.VariantLight {
			return color.RGBA{R: 33, G: 33, B: 33, A: 255} // Dark text
		}
		return color.RGBA{R: 255, G: 255, B: 255, A: 255} // Light text
	case theme.ColorNameSuccess:
		return color.RGBA{R: 76, G: 175, B: 80, A: 255} // Green
	case theme.ColorNameWarning:
		return color.RGBA{R: 255, G: 193, B: 7, A: 255} // Amber
	case theme.ColorNameError:
		return color.RGBA{R: 244, G: 67, B: 54, A: 255} // Red
	}
	
	// Fall back to default theme for other colors
	return theme.DefaultTheme().Color(name, variant)
}

// Font returns theme fonts with modern typography
func (m ModernTheme) Font(style fyne.TextStyle) fyne.Resource {
	// Use default fonts but could be customized with modern font resources
	return theme.DefaultTheme().Font(style)
}

// Icon returns theme icons
func (m ModernTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

// Size returns theme sizes with modern spacing
func (m ModernTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNamePadding:
		return 8 // Increased padding for modern look
	case theme.SizeNameInlineIcon:
		return 20 // Slightly larger icons
	case theme.SizeNameScrollBar:
		return 12 // Thinner scroll bars
	case theme.SizeNameSeparatorThickness:
		return 1 // Thin separators
	case theme.SizeNameInputBorder:
		return 2 // Slightly thicker input borders
	}
	
	// Fall back to default theme for other sizes
	return theme.DefaultTheme().Size(name)
}

// applyCustomTheme applies the custom theme to the application
func (a *App) applyCustomTheme() {
	// Note: In Fyne v2, custom themes need to be set at the app level
	// This is a placeholder for future theme customization
	// For now, we rely on the enhanced UI components and styling
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
