package main

import (
	"database/sql"
	"log"
	_ "github.com/mattn/go-sqlite3"
)

// migrateDatabase handles all database migrations
func migrateDatabase() {
	log.Println("Starting database migrations...")
	
	// Migrate messages table to add conversation_id
	migrateMessagesTable()
	
	// Migrate settings table to add selection column
	migrateSettingsTable()
	
	// Migrate conversations table to add summary and last_llm
	migrateConversationsTable()
	
	// Initialize llmSelection with default value
	initializeLLMSelection()
	
	log.Println("Database migrations completed.")
}

func migrateConversationsTable() {
	db, err := sql.Open("sqlite3", "DB/messages.db")
	if err != nil {
		log.Printf("Error opening database for conversations migration: %v", err)
		return
	}
	defer db.Close()
	
	// Check if conversations table exists
	var tableExists int
	err = db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='conversations'").Scan(&tableExists)
	if err != nil {
		log.Printf("Error checking conversations table existence: %v", err)
		return
	}
	
	if tableExists == 0 {
		// Table doesn't exist, create it with full schema
		_, err = db.Exec(`
			CREATE TABLE conversations (
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
		return
	}
	
	// Check if summary column exists
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM pragma_table_info('conversations') WHERE name='summary'").Scan(&count)
	if err == nil && count == 0 {
		log.Println("Adding summary column to conversations table...")
		_, err = db.Exec("ALTER TABLE conversations ADD COLUMN summary TEXT DEFAULT ''")
		if err != nil {
			log.Printf("Error adding summary column: %v", err)
		}
	}
	
	// Check if last_llm column exists
	count = 0
	err = db.QueryRow("SELECT COUNT(*) FROM pragma_table_info('conversations') WHERE name='last_llm'").Scan(&count)
	if err == nil && count == 0 {
		log.Println("Adding last_llm column to conversations table...")
		_, err = db.Exec("ALTER TABLE conversations ADD COLUMN last_llm TEXT DEFAULT ''")
		if err != nil {
			log.Printf("Error adding last_llm column: %v", err)
		}
	}
}

func migrateMessagesTable() {
	db, err := sql.Open("sqlite3", MessagesDB)
	if err != nil {
		log.Printf("Error opening messages database: %v", err)
		return
	}
	defer db.Close()
	
	// Check if conversation_id column exists
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM pragma_table_info('messages') WHERE name='conversation_id'").Scan(&count)
	if err != nil {
		log.Printf("Error checking for conversation_id column: %v", err)
		return
	}
	
	if count == 0 {
		log.Println("Adding conversation_id column to messages table...")
		_, err = db.Exec("ALTER TABLE messages ADD COLUMN conversation_id TEXT DEFAULT 'default'")
		if err != nil {
			log.Printf("Error adding conversation_id column: %v", err)
		} else {
			log.Println("Successfully added conversation_id column")
		}
	}
}

func migrateSettingsTable() {
	db, err := sql.Open("sqlite3", SettingsDB)
	if err != nil {
		log.Printf("Error opening settings database: %v", err)
		return
	}
	defer db.Close()
	
	// Check if selection column exists
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM pragma_table_info('settings') WHERE name='selection'").Scan(&count)
	if err != nil {
		log.Printf("Error checking for selection column: %v", err)
		return
	}
	
	if count == 0 {
		log.Println("Adding selection column to settings table...")
		_, err = db.Exec("ALTER TABLE settings ADD COLUMN selection TEXT DEFAULT 'OpenAI'")
		if err != nil {
			log.Printf("Error adding selection column: %v", err)
		} else {
			log.Println("Successfully added selection column")
		}
	}
}

func initializeLLMSelection() {
	db, err := sql.Open("sqlite3", "DB/llmSelection.db")
	if err != nil {
		log.Printf("Error opening llmSelection database: %v", err)
		return
	}
	defer db.Close()
	
	// Check if there's already a selection
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM llmSelection").Scan(&count)
	if err != nil {
		log.Printf("Error checking llmSelection: %v", err)
		return
	}
	
	if count == 0 {
		log.Println("Initializing llmSelection with default value...")
		_, err = db.Exec("INSERT INTO llmSelection (selection) VALUES ('OpenAI')")
		if err != nil {
			log.Printf("Error initializing llmSelection: %v", err)
		} else {
			log.Println("Successfully initialized llmSelection")
		}
	}
}

