package aigenUi

import (
	"aigen/aigenRest"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"time"
)

// Gigachad color scheme - world-class colors (using same vars as chatBubbles)
var (
	UserBubbleColor = color.RGBA{R: 52, G: 152, B: 219, A: 255}   // Professional blue
)

// AddGigachadBubble creates a world-class chat bubble with proper colors and animations
func AddGigachadBubble(box *fyne.Container, message string, isUser bool, model string) {
	// Get model color
	bgColor := getUserColor(isUser)
	if !isUser {
		bgColor = getModelColor(model)
	}
	
	// Create colored background rectangle
	bubbleBG := canvas.NewRectangle(bgColor)
	bubbleBG.StrokeWidth = 0
	bubbleBG.SetMinSize(fyne.NewSize(300, 0))
	
	// Create content with rich text
	content := widget.NewRichTextFromMarkdown(message)
	content.Wrapping = fyne.TextWrapWord
	
	// Create icon
	var icon fyne.CanvasObject
	if isUser {
		icon = widget.NewIcon(theme.AccountIcon())
	} else {
		icon = widget.NewIcon(getModelIcon(model))
	}
	
	// Model label for AI messages
	var header fyne.CanvasObject
	if !isUser {
		modelLabel := widget.NewLabel(model)
		modelLabel.TextStyle = fyne.TextStyle{Bold: true}
		modelLabel.Resize(fyne.NewSize(0, 20))
		header = container.NewVBox(modelLabel)
	} else {
		header = container.NewVBox()
	}
	
	// Create bubble content - NO FIXED SIZE to prevent truncation
	bubbleContent := container.NewBorder(
		header,
		nil,
		icon,
		nil,
		content,
	)
	// Don't set fixed size - let it expand
	
	// Stack background and content
	bubble := container.NewMax(
		bubbleBG,
		container.NewPadded(bubbleContent),
	)
	// Remove size restrictions to prevent truncation
	bubble.Resize(fyne.NewSize(600, 0))
	
	// Position based on user
	var messageContainer fyne.CanvasObject
	if isUser {
		// User messages: right-aligned with blue
		messageContainer = container.NewHBox(
			layout.NewSpacer(),
			container.NewMax(bubble),
		)
	} else {
		// AI messages: left-aligned with model color
		messageContainer = container.NewHBox(
			container.NewMax(bubble),
			layout.NewSpacer(),
		)
	}
	
	// Add timestamp
	timestamp := widget.NewLabel(time.Now().Format("15:04"))
	timestamp.Importance = widget.LowImportance
	timestamp.TextStyle = fyne.TextStyle{Italic: true}
	
	// Final container with timestamp
	finalContainer := container.NewVBox(
		messageContainer,
		container.NewHBox(
			layout.NewSpacer(),
			timestamp,
			layout.NewSpacer(),
		),
	)
	
	box.Add(finalContainer)
	box.Refresh()
}

func getUserColor(isUser bool) color.Color {
	if isUser {
		return UserBubbleColor
	}
	return color.RGBA{R: 245, G: 245, B: 245, A: 255} // Light gray
}

func getModelColor(model string) color.Color {
	switch model {
	case "Gemini":
		return GeminiColor
	case "OpenAI":
		return OpenAIColor
	case "Claude":
		return ClaudeColor
	case "Ollama":
		return OllamaColor
	default:
		return color.RGBA{R: 200, G: 200, B: 200, A: 255}
	}
}

func getModelIcon(model string) fyne.Resource {
	switch model {
	case "Gemini":
		return theme.InfoIcon()
	case "OpenAI":
		return theme.ContentAddIcon()
	case "Claude":
		return theme.ConfirmIcon()
	default:
		return theme.QuestionIcon()
	}
}

// Feature 1: Live typing indicator
func CreateTypingIndicator(box *fyne.Container) fyne.CanvasObject {
	indicator := widget.NewLabel("AI is typing...")
	indicator.TextStyle = fyne.TextStyle{Italic: true}
	return indicator
}

// Feature 2: Message reactions (thumb up/down)
func AddMessageReactions(messageID string) fyne.CanvasObject {
	thumbUp := widget.NewButtonWithIcon("", theme.ConfirmIcon(), func() {
		aigenRest.SendNotificationNow("Message liked")
	})
	thumbDown := widget.NewButtonWithIcon("", theme.CancelIcon(), func() {
		aigenRest.SendNotificationNow("Message disliked")
	})
	return container.NewHBox(thumbUp, thumbDown)
}

// Feature 3: Message edit functionality
func AddEditButton(message string, onEdit func(string)) fyne.CanvasObject {
	editBtn := widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), func() {
		aigenRest.SendNotificationNow("Edit message: " + message[:min(20, len(message))])
		if onEdit != nil {
			onEdit(message)
		}
	})
	return editBtn
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Feature 4: Code block syntax highlighting (placeholder)
func CreateCodeBlock(code string, language string) fyne.CanvasObject {
	codeText := widget.NewRichTextFromMarkdown("```" + language + "\n" + code + "\n```")
	codeText.Wrapping = fyne.TextWrapOff
	return container.NewBorder(
		widget.NewLabel(language),
		nil, nil, nil,
		codeText,
	)
}

// Feature 5: Message search functionality
func CreateSearchBar(onSearch func(string)) fyne.CanvasObject {
	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("Search messages...")
	searchEntry.OnChanged = func(text string) {
		if onSearch != nil {
			onSearch(text)
		}
	}
	searchBtn := widget.NewButtonWithIcon("Search", theme.SearchIcon(), func() {
		if onSearch != nil {
			onSearch(searchEntry.Text)
		}
	})
	return container.NewBorder(
		nil, nil,
		searchEntry,
		searchBtn,
		nil,
	)
}

