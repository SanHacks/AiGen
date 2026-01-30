package aigenUi

import (
	"aigen/aigenRest"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"log"
)

// BrowserTab creates a secure in-app browser
func BrowserTab(app fyne.App) *container.TabItem {
	// Create browser instance (will be initialized on demand)
	var browser *aigenRest.SecureBrowser

	// URL input
	urlEntry := widget.NewEntry()
	urlEntry.SetPlaceHolder("Enter URL (e.g., https://example.com)")

	// Content display
	contentDisplay := widget.NewRichTextFromMarkdown("# 🔒 Secure Browser\n\nNavigate safely with AI-powered web browsing.\n\n**Features:**\n• Secure Chromium-based browser\n• Sandboxed environment\n• Content extraction\n• Screenshot capability")

	// Navigate button
	navigateBtn := widget.NewButtonWithIcon("Navigate", theme.InfoIcon(), func() {
		url := urlEntry.Text
		if url == "" {
			dialog.ShowInformation("Error", "Please enter a URL", app.Driver().AllWindows()[0])
			return
		}

		// Initialize browser if not already done
		if browser == nil {
			dialog.ShowInformation("Initializing", "Starting secure browser...", app.Driver().AllWindows()[0])
			
			var err error
			browser, err = aigenRest.NewSecureBrowser()
			if err != nil {
				log.Printf("Browser initialization error: %v", err)
				dialog.ShowError(fmt.Errorf("Failed to initialize browser: %v", err), app.Driver().AllWindows()[0])
				return
			}
			
			dialog.ShowInformation("Ready", "Secure browser initialized!", app.Driver().AllWindows()[0])
		}

		// Navigate to URL
		dialog.ShowInformation("Loading", fmt.Sprintf("Loading %s...", url), app.Driver().AllWindows()[0])
		
		if err := browser.Navigate(url); err != nil {
			dialog.ShowError(fmt.Errorf("Failed to navigate: %v", err), app.Driver().AllWindows()[0])
			return
		}

		// Get page content
		content, err := browser.GetTextContent()
		if err != nil {
			dialog.ShowError(fmt.Errorf("Failed to get content: %v", err), app.Driver().AllWindows()[0])
			return
		}

		contentDisplay.ParseMarkdown(content)
		dialog.ShowInformation("Success", "Page loaded successfully!", app.Driver().AllWindows()[0])
	})

	// Search button
	searchBtn := widget.NewButtonWithIcon("Extract Text", theme.SearchIcon(), func() {
		// This would use AI to extract important content
		dialog.ShowInformation("AI Extraction", "Extracting main content with AI...", app.Driver().AllWindows()[0])
		// TODO: Implement AI-powered content extraction
	})

	// Screenshot button
	screenshotBtn := widget.NewButtonWithIcon("Screenshot", theme.MediaPhotoIcon(), func() {
		if browser == nil {
			dialog.ShowInformation("Error", "Please navigate to a page first", app.Driver().AllWindows()[0])
			return
		}

		filename := fmt.Sprintf("cache/screenshot_%d.png", 0)
		if err := browser.Screenshot(filename); err != nil {
			dialog.ShowError(err, app.Driver().AllWindows()[0])
			return
		}

		dialog.ShowInformation("Screenshot", fmt.Sprintf("Screenshot saved to %s", filename), app.Driver().AllWindows()[0])
	})

	// Layout
	topBar := container.NewBorder(
		nil,
		nil,
		navigateBtn,
		container.NewHBox(searchBtn, screenshotBtn),
		urlEntry,
	)

	scrollContent := container.NewScroll(contentDisplay)

	content := container.NewBorder(
		topBar,
		nil,
		nil,
		nil,
		scrollContent,
	)

	browserTab := container.NewTabItem("🌐 Browser", content)
	
	return browserTab
}

