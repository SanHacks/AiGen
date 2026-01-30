package main

import (
	"aigen/aigenRest"
	"aigen/aigenUi"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"log"
)

func bottomInputBox(chat *fyne.Container, tabs *container.AppTabs, aiGen *container.TabItem) *container.Split {
	inputBox := widget.NewMultiLineEntry()
	inputBox.Wrapping = fyne.TextWrapWord
	//TODO: Edit InputBox OnChanged
	inputBox.OnChanged = func(s string) {
		log.Printf("Input changed to: %s", s)
		kitchenLog(s)
	}
	//TODO: Edit InputBox Container
	inputBox.PlaceHolder = "Enter your message here..."
	
	// Model selector dropdown
	selectedModel, err := aigenUi.GetSelectedModel()
	if err != nil {
		log.Printf("Error getting selected model: %v", err)
	}
	modelSelector := widget.NewSelect([]string{"OpenAI", "Claude", "Gemini", "Ollama"}, func(value string) {
		log.Printf("Model selected: %s", value)
		err := aigenUi.UpdateSelectedModel(value)
		if err != nil {
			log.Printf("Error updating model: %v", err)
		}
		aigenRest.SendNotificationNow("Switched to " + value)
	})
	
	// Set the initial selected model
	if selectedModel != "" {
		modelSelector.SetSelected(selectedModel)
	} else {
		modelSelector.SetSelected("OpenAI")
	}
	
	sendButton := sendButton(inputBox, chat)
	voiceNoteButton := voiceChatButton(inputBox, chat)
	mediaUploadButton := imageUploadInput(inputBox, chat)
	voiceNoteButton.Resize(fyne.NewSize(50, 100))
	
	// Top row with model selector and buttons
	buttonRow := container.NewHBox(sendButton)
	
	// Left side: Model selector and input box
	leftSide := container.NewVBox(
		modelSelector,
		inputBox,
	)
	
	// Right side: Buttons
	rightSide := container.NewVBox(
		buttonRow,
		voiceNoteButton,
		mediaUploadButton,
	)
	
	// Main container with left (input + model) and right (buttons)
	inputBoxContainer := container.NewHSplit(leftSide, rightSide)
	
	tabs.OnSelected = func(tab *container.TabItem) {
		if tab == aiGen {
			//Show if tab is home
			inputBoxContainer.Show()
		} else {
			inputBoxContainer.Hide()
		}
	}
	return inputBoxContainer
}
