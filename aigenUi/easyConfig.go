package aigenUi

import (
	"log"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// EasyConfigTab creates a user-friendly configuration tab
func EasyConfigTab(app fyne.App) *container.TabItem {
	// API Keys Section
	openaiKey := widget.NewPasswordEntry()
	openaiKey.SetPlaceHolder("Enter your OpenAI API key")
	if key := os.Getenv("OPENAI"); key != "" {
		openaiKey.SetText(key)
	}

	claudeKey := widget.NewPasswordEntry()
	claudeKey.SetPlaceHolder("Enter your Anthropic API key")
	if key := os.Getenv("ANTHROPIC"); key != "" {
		claudeKey.SetText(key)
	}

	geminiKey := widget.NewPasswordEntry()
	geminiKey.SetPlaceHolder("Enter your Gemini API key")
	if key := os.Getenv("GEMINI_API_KEY"); key != "" {
		geminiKey.SetText(key)
	}

	elevenLabsKey := widget.NewPasswordEntry()
	elevenLabsKey.SetPlaceHolder("Enter your ElevenLabs API key (for voice)")
	if key := os.Getenv("ELEVENLABS_API_KEY"); key != "" {
		elevenLabsKey.SetText(key)
	}

	azureSpeechKey := widget.NewPasswordEntry()
	azureSpeechKey.SetPlaceHolder("Enter your Azure Speech key")
	if key := os.Getenv("SPEECH_KEY"); key != "" {
		azureSpeechKey.SetText(key)
	}

	// Save button
	saveButton := widget.NewButtonWithIcon("Save Configuration", theme.DocumentSaveIcon(), func() {
		// Save OpenAI key
		if openaiKey.Text != "" {
			os.Setenv("OPENAI", openaiKey.Text)
			log.Println("OpenAI API key updated")
		}

		// Save Claude key
		if claudeKey.Text != "" {
			os.Setenv("ANTHROPIC", claudeKey.Text)
			log.Println("Anthropic API key updated")
		}

		// Save Gemini key
		if geminiKey.Text != "" {
			os.Setenv("GEMINI_API_KEY", geminiKey.Text)
			log.Println("Gemini API key updated")
		}

		// Save ElevenLabs key
		if elevenLabsKey.Text != "" {
			os.Setenv("ELEVENLABS_API_KEY", elevenLabsKey.Text)
			log.Println("ElevenLabs API key updated")
		}

		// Save Azure Speech key
		if azureSpeechKey.Text != "" {
			os.Setenv("SPEECH_KEY", azureSpeechKey.Text)
			log.Println("Azure Speech key updated")
		}

		// Show success dialog
		dialog.ShowInformation("Success", "Configuration saved! Restart the app to apply changes.", app.Driver().AllWindows()[0])
	})
	saveButton.Importance = widget.HighImportance

	// Instructions card
	instructions := widget.NewCard("📝 How to Get API Keys", "",
		widget.NewLabel(
			"• OpenAI: https://platform.openai.com/api-keys\n"+
				"• Claude: https://console.anthropic.com/\n"+
				"• Gemini: https://makersuite.google.com/app/apikey\n"+
				"• ElevenLabs: https://elevenlabs.io/app/keys\n"+
				"• Azure Speech: https://portal.azure.com/\n\n"+
				"After saving, restart the application to use the new keys."))

	// Create form
	form := container.NewVBox(
		instructions,
		widget.NewSeparator(),
		widget.NewForm(
			widget.NewFormItem("OpenAI API Key", openaiKey),
			widget.NewFormItem("Claude API Key", claudeKey),
			widget.NewFormItem("Gemini API Key", geminiKey),
			widget.NewFormItem("ElevenLabs API Key (Voice)", elevenLabsKey),
			widget.NewFormItem("Azure Speech Key", azureSpeechKey),
		),
		widget.NewSeparator(),
		saveButton,
	)

	// Create scrollable container
	scroll := container.NewScroll(form)

	configTab := container.NewTabItem("🔧 Easy Setup", scroll)
	configTab.Icon = theme.SettingsIcon()
	return configTab
}

