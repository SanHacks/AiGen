package main

import (
	"aigen/aigenAudioAutoPlay"
	"aigen/aigenRest"
	"aigen/aigenUi"
	"fmt"
	"log"
	"os"
	"time"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)


// createModelSelector creates a modern model selector
func createModelSelector() *widget.Select {
	selectedModel, err := aigenUi.GetSelectedModel()
	if err != nil {
		selectedModel = "Gemini"
	}
	
	selector := widget.NewSelect([]string{"OpenAI", "Claude", "Gemini", "Ollama"}, func(value string) {
		aigenUi.UpdateSelectedModel(value)
		aigenRest.SendNotificationNow("Switched to " + value)
	})
	
	selector.SetSelected(selectedModel)
	return selector
}

// sendMessage handles sending messages with conversational context
func sendMessage(inputBox *widget.Entry, chat *fyne.Container) {
	message := inputBox.Text
	if message == "" {
		return
	}
	
	// Display user message with GIGACHAD bubbles
	aigenUi.AddGigachadBubble(chat, message, true, "")
	
	// Clear input
	inputBox.SetText("")
	
	// Show typing indicator
	typingIndicator := aigenUi.NewTypingIndicator()
	typingIndicator.Show()
	chat.Add(typingIndicator.GetContainer())
	chat.Refresh()
	
	// Get selected model
	model, err := aigenUi.GetSelectedModel()
	if err != nil {
		model = "Gemini"
	}
	
	// Check if conversational mode is enabled
	// In conversational mode, include last few messages for context
	contextMessage := message
	useConversational := false // TODO: Get from UI state
	
	if useConversational && model == "Gemini" {
		// Get last few messages for context
		lastMessages, err := getLastMessages()
		if err == nil && len(lastMessages) > 0 {
			// Build context from last messages
			contextBuilder := "Previous conversation:\n"
			for _, msg := range lastMessages[len(lastMessages)-5:] { // Last 5 messages
				contextBuilder += fmt.Sprintf("%s: %s\n", msg.Sender, msg.Content)
			}
			contextMessage = contextBuilder + "\nUser: " + message
		}
	}
	
	// Remove typing indicator
	defer func() {
		typingIndicator.Hide()
		chat.Refresh()
	}()
	
	// Call appropriate model
	var response string
	switch model {
	case "OpenAI":
		response, _ = aigenRest.MakeApiCall(message)
	case "Claude":
		response, _ = aigenRest.CallClaude(message)
	case "Gemini":
		if useConversational {
			response, _ = aigenRest.CallGemini(contextMessage)
		} else {
			response, _ = aigenRest.CallGemini(message)
		}
	case "Ollama":
		response, _ = aigenRest.CallOllama(message)
	}
	
	// Display bot response with GIGACHAD colored bubble based on model
	aigenUi.AddGigachadBubble(chat, response, false, model)
	
	// Save to database
	addMessage("YOU", message)
	addMessage("Bot", response)
	
	// Generate speech if audio is enabled
	if getAudioSettings() && len(response) > 0 {
		go generateSpeechAsync(response)
	}
}

// generateSpeechAsync generates speech in the background
func generateSpeechAsync(text string) {
	// Use ElevenLabs for high-quality TTS
	// Default voice: Rachel (female, friendly)
	voiceID := "21m00Tcm4TlvDq8ikWAM" // Rachel voice
	
	audioData, err := aigenRest.CallElevenLabs(text, voiceID)
	if err != nil {
		log.Printf("Error generating speech: %v", err)
		return
	}
	
	// Save audio to file
	filename := fmt.Sprintf("cache/%d.mp3", time.Now().Unix())
	err = os.WriteFile(filename, audioData, 0644)
	if err != nil {
		log.Printf("Error saving audio file: %v", err)
		return
	}
	
	// Play audio using existing function
	go func() {
		err := aigenAudioAutoPlay.PlayAudioPlayback(filename)
		if err != nil {
			log.Printf("Error playing audio: %v", err)
		}
	}()
}

// handleVoiceInput handles voice input
func handleVoiceInput(inputBox *widget.Entry, chat *fyne.Container) {
	aigenRest.SendNotificationNow("Voice recording started - 20 seconds")
	// TODO: Implement voice recording
}

// Modern input container with better UX
func createModernInputContainer(chat *fyne.Container, tabs *container.AppTabs) *fyne.Container {
	inputBox := widget.NewMultiLineEntry()
	inputBox.SetPlaceHolder("Type your message... (Press Enter to send)")
	inputBox.Wrapping = fyne.TextWrapWord
	
	// Handle Enter key to send
	inputBox.OnSubmitted = func(text string) {
		sendMessage(inputBox, chat)
	}
	
	// Send button
	sendBtn := widget.NewButtonWithIcon("Send", theme.MailSendIcon(), func() {
		sendMessage(inputBox, chat)
	})
	sendBtn.Importance = widget.HighImportance
	
	// Voice button
	voiceBtn := widget.NewButtonWithIcon("", theme.MediaRecordIcon(), func() {
		handleVoiceInput(inputBox, chat)
	})
	
	// Audio toggle button (enable/disable audio responses)
	audioToggle := widget.NewCheck("🔊 Audio", func(enabled bool) {
		if enabled {
			aigenRest.SendNotificationNow("Audio responses enabled")
		} else {
			aigenRest.SendNotificationNow("Audio responses disabled")
		}
		// Save preference
		if enabled {
			aigenUi.ChangeSetting(1)
		} else {
			aigenUi.ChangeSetting(0)
		}
	})
	audioToggle.SetChecked(true)
	
	// Conversational mode toggle (continuous conversation with context)
	conversationalMode := widget.NewCheck("💬 Conversational", func(enabled bool) {
		if enabled {
			aigenRest.SendNotificationNow("Conversational mode enabled - context preserved")
		} else {
			aigenRest.SendNotificationNow("Conversational mode disabled")
		}
	})
	conversationalMode.SetChecked(false)
	
	// Model selector for quick access
	modelSelector := createModelSelector()
	
	// Top row: Model selector, conversational mode, and audio toggle
	topRow := container.NewHBox(
		widget.NewLabel("Model:"),
		modelSelector,
		layout.NewSpacer(),
		conversationalMode,
		audioToggle,
	)
	
	// Input row: Voice button, input, and send button
	inputRow := container.NewBorder(
		nil,
		nil,
		voiceBtn,
		sendBtn,
		inputBox,
	)
	
	// Main container
	inputContainer := container.NewVBox(
		topRow,
		widget.NewSeparator(),
		inputRow,
	)
	
	return inputContainer
}

