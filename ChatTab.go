package main

import (
	"aigen/aigenUi"
	"aigen/essentialsGen"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"log"
)

// ChatTab Chat Bubble is a container that contains a label and an avatar
// The label contains the message
// The avatar contains the avatar of the sender
// The isUser parameter determines if the message is from the user or the bot
// If the message is from the user, the bubble will be on the right side of the chat window
// If the message is from the bot, the bubble will be on the left side of the chat window
func ChatTab() (*fyne.Container, *container.TabItem) {
	//Create the chat tab - now with better layout for desktop
	chat := container.NewVBox()
	chat.Refresh()

	// Use a subtle background instead of an image for better readability
	aiGen := container.NewTabItem("Sage Chat", chat)
	aiGen.Icon = theme.HomeIcon()

	messagesFromDB, err := getMessages()

	if err != nil {
		log.Printf("Error getting messages: %v", err)
	}
	
	//Loop Through Messages From DB and Display with GIGACHAD COLORS
	for _, message := range messagesFromDB {
		if message.Sender == "YOU" {
			aigenUi.AddGigachadBubble(chat, message.Content, true, "")
		} else {
			if message.Media != "NULL" {
				addMediaChatBubble(chat, message.Media, false)
			} else {
				// Use last known model or default
				model := "Gemini"
				aigenUi.AddGigachadBubble(chat, message.Content, false, model)
			}
		}
	}

	essentialsGen.StartUpCall()

	return chat, aiGen
}
