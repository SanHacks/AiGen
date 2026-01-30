package aigenUi

import (
	"aigen/aigenRest"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// Conversation represents a single conversation thread
type Conversation struct {
	ID    string
	Title string
	Model string
}

// Sidebar creates a modern conversation history sidebar
func Sidebar(app fyne.App, conversations []Conversation, onSelect func(string), chat *fyne.Container) fyne.CanvasObject {
	// Conversation list
	list := widget.NewList(
		func() int {
			return len(conversations)
		},
		func() fyne.CanvasObject {
			icon := widget.NewIcon(theme.ContentAddIcon())
			title := widget.NewLabel("Conversation title")
			title.Truncation = fyne.TextTruncateEllipsis
			model := widget.NewLabel("Model")
			model.Truncation = fyne.TextTruncateEllipsis
			vbox := container.NewVBox(title, model)
			return container.NewHBox(icon, vbox)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			if id < len(conversations) {
				conv := conversations[id]
				// Update labels in the vbox
				if hbox, ok := obj.(*fyne.Container); ok && len(hbox.Objects) >= 2 {
					if vbox, ok := hbox.Objects[1].(*fyne.Container); ok {
						if len(vbox.Objects) >= 2 {
							vbox.Objects[0].(*widget.Label).SetText(conv.Title)
							vbox.Objects[1].(*widget.Label).SetText(conv.Model)
						}
					}
				}
			}
		},
	)

	list.OnSelected = func(id widget.ListItemID) {
		if id < len(conversations) {
			onSelect(conversations[id].ID)
		}
	}

	// New chat button - CLEARS CHAT to start fresh conversation
	newChatBtn := widget.NewButtonWithIcon("New Chat", theme.ContentAddIcon(), func() {
		// Clear the chat container
		chat.Objects = []fyne.CanvasObject{}
		chat.Refresh()
		
		// Reset conversation ID
		// currentConversationID = "new_" + time.Now().Format("20060102150405")
		
		aigenRest.SendNotificationNow("New conversation started!")
	})

	// Theme toggle button
	themeBtn := widget.NewButtonWithIcon("Theme", theme.ColorPaletteIcon(), func() {
		// Toggle theme
		if app != nil {
			currentTheme := app.Settings().Theme()
			if currentTheme == theme.LightTheme() {
				app.Settings().SetTheme(theme.DarkTheme())
				aigenRest.SendNotificationNow("🌙 Dark mode enabled")
			} else {
				app.Settings().SetTheme(theme.LightTheme())
				aigenRest.SendNotificationNow("☀️ Light mode enabled")
			}
		}
	})
	
	// Header with theme button
	header := container.NewBorder(
		container.NewBorder(
			nil, nil,
			widget.NewIcon(theme.SettingsIcon()),
			container.NewHBox(themeBtn, newChatBtn),
			widget.NewLabel("Conversations"),
		),
		nil, nil, nil, nil,
	)

	// Sidebar content
	content := container.NewBorder(
		header,
		nil,
		nil,
		nil,
		list,
	)

	return content
}

// ModernSidebarWithChat creates sidebar with chat container reference
func ModernSidebarWithChat(application fyne.App, chat *fyne.Container) fyne.CanvasObject {
	// Sample conversations for now
	// TODO: Load from database via dbHistory.go
	conversations := []Conversation{
		{ID: "1", Title: "Quantum Computing Basics", Model: "Gemini"},
		{ID: "2", Title: "Python Data Analysis", Model: "OpenAI"},
		{ID: "3", Title: "Web Development Tips", Model: "Claude"},
	}

	// Create sidebar
	onSelect := func(id string) {
		// Handle conversation selection
		aigenRest.SendNotificationNow("Selected: " + id)
	}

	sidebar := Sidebar(application, conversations, onSelect, chat)
	
	// Set width for desktop
	sidebar.Resize(fyne.NewSize(250, 0))
	
	return sidebar
}

