package aigenRest

import (
	"fmt"
	"net/http"
	"time"
)

// CheckOllamaConnection checks if Ollama is running locally
func CheckOllamaConnection() (bool, error) {
	client := &http.Client{
		Timeout: 2 * time.Second,
	}
	
	resp, err := client.Get("http://localhost:11434")
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode == http.StatusOK {
		return true, nil
	}
	
	return false, fmt.Errorf("Ollama responded with status: %d", resp.StatusCode)
}

// CheckOllamaRunning returns true if Ollama is available
func CheckOllamaRunning() bool {
	ok, _ := CheckOllamaConnection()
	return ok
}

