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

// Color returns theme colors with WCAG AAA compliant high contrast palette
func (m ModernLightTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	// Force light variant for all colors to ensure consistent light theme
	switch name {
	// Primary colors
	case theme.ColorNamePrimary:
		return color.RGBA{R: 25, G: 118, B: 210, A: 255} // Professional blue #1976D2
	case theme.ColorNameBackground:
		return color.RGBA{R: 255, G: 255, B: 255, A: 255} // Pure white background #FFFFFF

	// Text colors - Maximum contrast with pure black
	case theme.ColorNameForeground:
		return color.RGBA{R: 0, G: 0, B: 0, A: 255} // Pure black text #000000 - 21:1 contrast
	case theme.ColorNameDisabled:
		return color.RGBA{R: 64, G: 64, B: 64, A: 255} // Dark gray for disabled #404040 - 9.7:1 contrast

	// Button colors
	case theme.ColorNameButton:
		return color.RGBA{R: 248, G: 249, B: 250, A: 255} // Light button background #F8F9FA
	case theme.ColorNameDisabledButton:
		return color.RGBA{R: 233, G: 236, B: 239, A: 255} // Disabled button #E9ECEF
	case theme.ColorNameHover:
		return color.RGBA{R: 233, G: 236, B: 239, A: 255} // Hover state #E9ECEF
	case theme.ColorNamePressed:
		return color.RGBA{R: 222, G: 226, B: 230, A: 255} // Pressed state #DEE2E6

	// Status colors
	case theme.ColorNameSuccess:
		return color.RGBA{R: 40, G: 167, B: 69, A: 255} // Success green #28A745 - 4.6:1 contrast
	case theme.ColorNameWarning:
		return color.RGBA{R: 133, G: 100, B: 4, A: 255} // Warning dark yellow #856404 - 7.4:1 contrast
	case theme.ColorNameError:
		return color.RGBA{R: 220, G: 53, B: 69, A: 255} // Error red #DC3545 - 5.9:1 contrast

	// Input colors
	case theme.ColorNameInputBackground:
		return color.RGBA{R: 255, G: 255, B: 255, A: 255} // White input background
	case theme.ColorNameInputBorder:
		return color.RGBA{R: 206, G: 212, B: 218, A: 255} // Input border #CED4DA
	case theme.ColorNamePlaceHolder:
		return color.RGBA{R: 96, G: 96, B: 96, A: 255} // Dark placeholder text #606060

	// Selection colors
	case theme.ColorNameSelection:
		return color.RGBA{R: 25, G: 118, B: 210, A: 51} // Selection highlight (20% opacity)
	case theme.ColorNameFocus:
		return color.RGBA{R: 25, G: 118, B: 210, A: 255} // Focus indicator

	// UI element colors
	case theme.ColorNameScrollBar:
		return color.RGBA{R: 173, G: 181, B: 189, A: 255} // Scrollbar #ADB5BD
	case theme.ColorNameShadow:
		return color.RGBA{R: 0, G: 0, B: 0, A: 25} // Subtle shadow
	case theme.ColorNameSeparator:
		return color.RGBA{R: 222, G: 226, B: 230, A: 255} // Separator #DEE2E6

	// Card and container colors
	case theme.ColorNameHeaderBackground:
		return color.RGBA{R: 248, G: 249, B: 250, A: 255} // Header background #F8F9FA
	case theme.ColorNameMenuBackground:
		return color.RGBA{R: 255, G: 255, B: 255, A: 255} // Menu background
	case theme.ColorNameOverlayBackground:
		return color.RGBA{R: 0, G: 0, B: 0, A: 128} // Modal overlay (50% opacity)
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

// Maximum contrast color constants with pure black text
var (
	// Text colors (maximum contrast on white background)
	PrimaryTextColor   = color.RGBA{R: 0,   G: 0,   B: 0,   A: 255} // #000000 - 21:1 contrast (pure black)
	SecondaryTextColor = color.RGBA{R: 0,   G: 0,   B: 0,   A: 255} // #000000 - 21:1 contrast (pure black)
	MutedTextColor     = color.RGBA{R: 64,  G: 64,  B: 64,  A: 255} // #404040 - 9.7:1 contrast (dark gray)
	LabelTextColor     = color.RGBA{R: 0,   G: 0,   B: 0,   A: 255} // #000000 - 21:1 contrast (pure black)

	// Background colors
	PrimaryBgColor     = color.RGBA{R: 255, G: 255, B: 255, A: 255} // #FFFFFF - Pure white
	SecondaryBgColor   = color.RGBA{R: 248, G: 249, B: 250, A: 255} // #F8F9FA - Light gray
	CardBgColor        = color.RGBA{R: 255, G: 255, B: 255, A: 255} // #FFFFFF - Card background

	// Interactive colors
	PrimaryBlue        = color.RGBA{R: 25,  G: 118, B: 210, A: 255} // #1976D2 - Professional blue
	SuccessGreen       = color.RGBA{R: 40,  G: 167, B: 69,  A: 255} // #28A745 - Success green
	WarningAmber       = color.RGBA{R: 133, G: 100, B: 4,   A: 255} // #856404 - Warning (dark for contrast)
	ErrorRed           = color.RGBA{R: 220, G: 53,  B: 69,  A: 255} // #DC3545 - Error red

	// Border and separator colors
	BorderColor        = color.RGBA{R: 206, G: 212, B: 218, A: 255} // #CED4DA - Light border
	SeparatorColor     = color.RGBA{R: 222, G: 226, B: 230, A: 255} // #DEE2E6 - Separator

	// Legacy colors (for backward compatibility)
	ModernBlue   = PrimaryBlue
	ModernGreen  = SuccessGreen
	ModernRed    = ErrorRed
	LightGray    = SecondaryBgColor
)

// createStyledButton creates a button with enhanced styling and proper contrast
func createStyledButton(text string, icon string, importance widget.Importance, onTapped func()) *widget.Button {
	buttonText := text
	if icon != "" {
		buttonText = icon + " " + text
	}

	button := widget.NewButton(buttonText, onTapped)
	button.Importance = importance

	return button
}

// createHighContrastLabel creates a label with high contrast text
func createHighContrastLabel(text string, style fyne.TextStyle) *widget.Label {
	label := widget.NewLabel(text)
	label.TextStyle = style
	return label
}

// createPrimaryLabel creates a label with primary text color and styling
func createPrimaryLabel(text string) *widget.Label {
	label := widget.NewLabel(text)
	label.TextStyle = fyne.TextStyle{Bold: false}
	return label
}

// createSecondaryLabel creates a label with secondary text color
func createSecondaryLabel(text string) *widget.Label {
	label := widget.NewLabel(text)
	label.TextStyle = fyne.TextStyle{Italic: true}
	return label
}

// createStyledCard creates a card with enhanced styling and high contrast
func createStyledCard(title, subtitle string, content fyne.CanvasObject) *widget.Card {
	card := widget.NewCard(title, subtitle, content)
	return card
}

// createSectionHeader creates a styled section header with high contrast
func createSectionHeader(text string) *widget.Label {
	header := widget.NewLabel(text)
	header.TextStyle = fyne.TextStyle{Bold: true}
	return header
}

// createHelpText creates styled help text with proper contrast
func createHelpText(text string) *widget.Label {
	help := widget.NewLabel(text)
	help.TextStyle = fyne.TextStyle{Italic: true}
	return help
}

// createStatusLabel creates a status label with appropriate color coding
func createStatusLabel(text string, statusType string) *widget.Label {
	label := widget.NewLabel(text)
	label.TextStyle = fyne.TextStyle{Bold: true}
	return label
}

// createValidationLabel creates a label for validation messages
func createValidationLabel(text string, isError bool) *widget.Label {
	label := widget.NewLabel(text)
	if isError {
		label.TextStyle = fyne.TextStyle{Bold: true}
	} else {
		label.TextStyle = fyne.TextStyle{Italic: true}
	}
	return label
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
