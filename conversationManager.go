package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"fyne.io/fyne/v2"
)

var currentConversationID = "default"

// Conversation represents a single conversation
type Conversation struct {
	ID        string
	Title     string
	Model     string
	Summary   string
	LastLLM   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// CreateNewConversation creates a new conversation and returns its ID
func CreateNewConversation(model string) (string, error) {
	db, err := sql.Open("sqlite3", "DB/messages.db")
	if err != nil {
		return "", err
	}
	defer db.Close()

	// Create conversations table if it doesn't exist with expanded schema
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS conversations (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			model TEXT DEFAULT 'Gemini',
			summary TEXT,
			last_llm TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Printf("Error creating conversations table: %v", err)
	}

	// Generate a unique ID
	conversationID := fmt.Sprintf("conv_%d", time.Now().Unix())
	title := fmt.Sprintf("Chat %s", time.Now().Format("15:04"))

	// Insert new conversation with summary and last LLM
	_, err = db.Exec(`
		INSERT INTO conversations (id, title, model, summary, last_llm, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, conversationID, title, model, "", model, time.Now(), time.Now())
	if err != nil {
		return "", err
	}

	currentConversationID = conversationID
	return conversationID, nil
}

// LoadConversation loads messages for a specific conversation
func LoadConversation(chat *fyne.Container, conversationID string) error {
	// Clear current chat
	chat.Objects = []fyne.CanvasObject{}

	// Load messages for this conversation
	messages, err := getMessagesForConversation(conversationID)
	if err != nil {
		return err
	}

	// Display messages
	for _, msg := range messages {
		if msg.Sender == "YOU" {
			addChatBubble(chat, msg.Content, true)
		} else {
			addChatBubble(chat, msg.Content, false)
		}
	}

	currentConversationID = conversationID
	return nil
}

// GetConversations returns all conversations
func GetConversations() ([]Conversation, error) {
	db, err := sql.Open("sqlite3", "DB/messages.db")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query(`
		SELECT id, title, model, created_at, updated_at
		FROM conversations
		ORDER BY updated_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var conversations []Conversation
	for rows.Next() {
		var conv Conversation
		err := rows.Scan(&conv.ID, &conv.Title, &conv.Model, &conv.CreatedAt, &conv.UpdatedAt)
		if err != nil {
			continue
		}
		conversations = append(conversations, conv)
	}

	return conversations, nil
}

// getMessagesForConversation gets messages for a specific conversation
func getMessagesForConversation(conversationID string) ([]Message, error) {
	db, err := sql.Open("sqlite3", MessagesDB)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query(`
		SELECT id, sender, content, media, conversation_id, created_at
		FROM messages
		WHERE conversation_id = ?
		ORDER BY created_at
	`, conversationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		err := rows.Scan(&msg.ID, &msg.Sender, &msg.Content, &msg.Media, &msg.ConversationID, &msg.CreatedAt)
		if err != nil {
			continue
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

// SwitchConversation switches to a different conversation
func SwitchConversation(app fyne.App, chat *fyne.Container, conversationID string) error {
	// Load the conversation
	err := LoadConversation(chat, conversationID)
	if err != nil {
		return err
	}

	// Refresh the window
	w := app.Driver().AllWindows()[0]
	w.Content().Refresh()

	return nil
}

// UpdateConversationTitle updates the title of a conversation
func UpdateConversationTitle(conversationID, title string) error {
	db, err := sql.Open("sqlite3", "DB/messages.db")
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(`
		UPDATE conversations
		SET title = ?, updated_at = ?
		WHERE id = ?
	`, title, time.Now(), conversationID)

	return err
}

// AddMessage adds a message with the current conversation ID
func AddMessageWithConversation(sender, content string) error {
	db, err := sql.Open("sqlite3", MessagesDB)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(`
		INSERT INTO messages (sender, content, conversation_id)
		VALUES (?, ?, ?)
	`, sender, content, currentConversationID)

	return err
}

// GetCurrentConversationID returns the current conversation ID
func GetCurrentConversationID() string {
	return currentConversationID
}

// NewConversation creates a new conversation
func NewConversationChat(app fyne.App, chat *fyne.Container, model string) (*fyne.Container, func()) {
	// Get current model
	if model == "" {
		model = "Gemini"
	}
	
	// Create new conversation
	convID, err := CreateNewConversation(model)
	if err != nil {
		log.Printf("Error creating conversation: %v", err)
		convID = "default"
	}

	// Clear chat
	chat.Objects = []fyne.CanvasObject{}

	// Return update function
	return chat, func() {
		LoadConversation(chat, convID)
	}
}

