package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

const SitesConfigFile = "config/sites.json"

// Site represents a configurable iframe site
type Site struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	URL         string `json:"url"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

// SitesConfig holds the list of configurable sites
type SitesConfig struct {
	Sites []Site `json:"sites"`
}

var defaultSites = []Site{
	{
		ID:          "snapcart",
		Name:        "SnapCart",
		URL:         "https://snapcartza.co.za/",
		Description: "E-commerce platform",
		Icon:        "🛒",
	},
	{
		ID:          "gemini",
		Name:        "Gemini",
		URL:         "https://gemini.google.com/",
		Description: "Google Gemini AI",
		Icon:        "🤖",
	},
	{
		ID:          "chatgpt",
		Name:        "ChatGPT",
		URL:         "https://chat.openai.com/",
		Description: "OpenAI ChatGPT",
		Icon:        "💬",
	},
}

// LoadSitesConfig loads sites from JSON file
func LoadSitesConfig() (*SitesConfig, error) {
	// Create config directory if it doesn't exist
	configDir := filepath.Dir(SitesConfigFile)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, err
	}
	
	// Check if file exists
	if _, err := os.Stat(SitesConfigFile); os.IsNotExist(err) {
		// Create default config
		config := &SitesConfig{Sites: defaultSites}
		saveSitesConfig(config)
		return config, nil
	}
	
	// Read existing config
	data, err := os.ReadFile(SitesConfigFile)
	if err != nil {
		return nil, err
	}
	
	var config SitesConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	
	return &config, nil
}

// SaveSitesConfig saves sites to JSON file
func saveSitesConfig(config *SitesConfig) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	
	return os.WriteFile(SitesConfigFile, data, 0644)
}

// AddSite adds a new site to the configuration
func AddSite(name, url, description, icon string) error {
	config, err := LoadSitesConfig()
	if err != nil {
		return err
	}
	
	// Generate ID
	id := fmt.Sprintf("site_%d", len(config.Sites)+1)
	
	site := Site{
		ID:          id,
		Name:        name,
		URL:         url,
		Description: description,
		Icon:        icon,
	}
	
	config.Sites = append(config.Sites, site)
	
	return saveSitesConfig(config)
}

// RemoveSite removes a site from the configuration
func RemoveSite(id string) error {
	config, err := LoadSitesConfig()
	if err != nil {
		return err
	}
	
	for i, site := range config.Sites {
		if site.ID == id {
			config.Sites = append(config.Sites[:i], config.Sites[i+1:]...)
			break
		}
	}
	
	return saveSitesConfig(config)
}

// GetSites returns all configured sites
func GetSites() ([]Site, error) {
	config, err := LoadSitesConfig()
	if err != nil {
		log.Printf("Error loading sites config: %v", err)
		return defaultSites, nil
	}
	
	return config.Sites, nil
}

