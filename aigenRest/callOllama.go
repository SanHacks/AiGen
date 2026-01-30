package aigenRest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

// CallOllama calls Ollama with better error handling
func CallOllama(message string) (string, error) {
	// Check if Ollama is running first
	if !CheckOllamaRunning() {
		return "", fmt.Errorf("Ollama is not running. Please start Ollama with: ollama serve")
	}

	// Try different models in order
	models := []string{"llama2", "llama3", "mistral", "codellama", "gemma"}
	
	for _, model := range models {
		response, err := tryOllamaModel(message, model)
		if err == nil {
			return response, nil
		}
		log.Printf("Model %s failed: %v", model, err)
	}
	
	return "", fmt.Errorf("failed to connect to Ollama with any model")
}

// tryOllamaModel tries a specific model
func tryOllamaModel(message, model string) (string, error) {
	url := "http://localhost:11434/api/generate"
	
	payload := map[string]interface{}{
		"model":  model,
		"stream": false,
		"prompt": message,
	}

	jsonData, _ := json.Marshal(payload)

	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	
	req, err := http.NewRequest("POST", url, strings.NewReader(string(jsonData)))
	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Ollama returned status: %d", res.StatusCode)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return "", err
	}

	// Try to get response
	if completion, ok := response["response"].(string); ok {
		return completion, nil
	}
	if completion, ok := response["completion"].(string); ok {
		return completion, nil
	}

	return "", fmt.Errorf("no completion found in response")
}
