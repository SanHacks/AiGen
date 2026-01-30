package aigenUi

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// ThemeToggle creates a dark/light mode toggle
func ThemeToggle(app fyne.App) *widget.Check {
	// Check if current theme is dark
	isDark := app.Settings().Theme() != theme.LightTheme()

	themeCheck := widget.NewCheck("", func(checked bool) {
		if checked {
			app.Settings().SetTheme(theme.LightTheme())
		} else {
			app.Settings().SetTheme(theme.DarkTheme())
		}
	})
	
	themeCheck.SetChecked(!isDark) // Inverted: checked = light mode
	themeCheck.SetText("☀️ Light")

	return themeCheck
}

// ThemePreference stores user's theme preference
type ThemePreference struct {
	Dark bool `json:"dark"`
}

// SaveThemePreference saves the theme preference
func SaveThemePreference(dark bool) error {
	// TODO: Save to database or preferences
	return nil
}

// LoadThemePreference loads the theme preference
func LoadThemePreference() bool {
	// TODO: Load from database or preferences
	return true // Default to dark
}

