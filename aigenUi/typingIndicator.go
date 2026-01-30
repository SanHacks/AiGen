package aigenUi

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// TypingIndicator shows "AI is typing..."
type TypingIndicator struct {
	container *fyne.Container
	label     *widget.Label
	isVisible bool
}

// NewTypingIndicator creates a new typing indicator
func NewTypingIndicator() *TypingIndicator {
	label := widget.NewLabel("AI is typing...")
	label.TextStyle = fyne.TextStyle{Italic: true}
	
	container := container.NewVBox(label)
	
	return &TypingIndicator{
		container: container,
		label:     label,
		isVisible: false,
	}
}

// Show displays the typing indicator
func (ti *TypingIndicator) Show() {
	if !ti.isVisible {
		ti.isVisible = true
		ti.label.SetText("AI is thinking...")
	}
}

// Hide removes the typing indicator
func (ti *TypingIndicator) Hide() {
	if ti.isVisible {
		ti.isVisible = false
		ti.label.SetText("")
	}
}

// GetContainer returns the container for the indicator
func (ti *TypingIndicator) GetContainer() *fyne.Container {
	return ti.container
}

