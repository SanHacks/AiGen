package aigenUi

import (
	"image/color"
	
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// Color scheme for professional chat
var (
	UserBubbleBG = color.RGBA{R: 52, G: 152, B: 219, A: 255}   // Blue
	AIBubbleBG   = color.RGBA{R: 236, G: 240, B: 241, A: 255}  // Light gray
	GeminiColor  = color.RGBA{R: 66, G: 133, B: 244, A: 255}   // Google blue
	OpenAIColor  = color.RGBA{R: 16, G: 163, B: 127, A: 255}   // OpenAI green
	ClaudeColor  = color.RGBA{R: 255, G: 107, B: 107, A: 255}  // Red-orange
	OllamaColor  = color.RGBA{R: 156, G: 39, B: 176, A: 255}   // Purple
)

// CreateUserBubble creates a user message bubble (blue, right side)
func CreateUserBubble(message string) fyne.CanvasObject {
	icon := widget.NewIcon(theme.AccountIcon())
	text := widget.NewRichTextFromMarkdown(message)
	text.Wrapping = fyne.TextWrapWord
	
	bubble := container.NewHBox(icon, text)
	return bubble
}

// CreateAIBubble creates an AI message bubble with model-specific styling
func CreateAIBubble(message string, model string) fyne.CanvasObject {
	// Get appropriate icon for model
	var icon fyne.CanvasObject
	switch model {
	case "Gemini":
		icon = widget.NewIcon(theme.InfoIcon())
	case "OpenAI":
		icon = widget.NewIcon(theme.ContentAddIcon())
	case "Claude":
		icon = widget.NewIcon(theme.ConfirmIcon())
	default:
		icon = widget.NewIcon(theme.QuestionIcon())
	}
	
	text := widget.NewRichTextFromMarkdown(message)
	text.Wrapping = fyne.TextWrapWord
	
	// Model label
	modelLabel := widget.NewLabel(model)
	modelLabel.TextStyle = fyne.TextStyle{Bold: true}
	
	bubble := container.NewVBox(
		modelLabel,
		container.NewHBox(icon, text),
	)
	return bubble
}
