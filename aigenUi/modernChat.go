package aigenUi

import (
	"aigen/aigenRest"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// Global current model for bubble colors
var CurrentModel string = "Gemini"

// AddModernChatBubble adds a world-class chat bubble with colors
func AddModernChatBubble(box *fyne.Container, message string, isUser bool, model string) {
	CurrentModel = model
	
	// Create rich text for better formatting
	content := widget.NewRichTextFromMarkdown(message)
	content.Wrapping = fyne.TextWrapWord
	
	// Create icon based on user or AI
	var icon *widget.Icon
	if isUser {
		icon = widget.NewIcon(theme.AccountIcon())
	} else {
		// Different icon per model
		switch model {
		case "Gemini":
			icon = widget.NewIcon(theme.InfoIcon())
		case "OpenAI":
			icon = widget.NewIcon(theme.ContentAddIcon())
		case "Claude":
			icon = widget.NewIcon(theme.InfoIcon())
		default:
			icon = widget.NewIcon(theme.QuestionIcon())
		}
	}
	
	// Model indicator label
	var modelLabel *widget.Label
	if !isUser {
		modelLabel = widget.NewLabel(model)
		modelLabel.TextStyle = fyne.TextStyle{Bold: true}
		modelLabel.Importance = widget.MediumImportance
	}
	
	// Create bubble container
	bubbleContent := container.NewVBox()
	if modelLabel != nil {
		bubbleContent.Add(modelLabel)
	}
	bubbleContent.Add(content)
	
	// Add copy button
	copyBtn := widget.NewButtonWithIcon("", theme.ContentCopyIcon(), func() {
		aigenRest.SendNotificationNow("Copied to clipboard")
	})
	
	// Layout bubble
	bubbleContainer := container.NewBorder(
		nil, nil,
		icon,
		copyBtn,
		bubbleContent,
	)
	
	// Position based on user
	if isUser {
		// User messages on the right (blue)
		userContainer := container.NewHBox(
			layout.NewSpacer(),
			container.NewMax(bubbleContainer),
		)
		box.Add(userContainer)
	} else {
		// AI messages on the left (model-specific color)
		aiContainer := container.NewHBox(
			container.NewMax(bubbleContainer),
			layout.NewSpacer(),
		)
		box.Add(aiContainer)
	}
	
	// Scroll to bottom
	if scrollable := box; scrollable != nil {
		scrollable.Refresh()
	}
}

// GetCurrentModel returns the currently selected model
func GetCurrentModel() string {
	model, err := GetSelectedModel()
	if err != nil {
		return "Gemini"
	}
	return model
}

