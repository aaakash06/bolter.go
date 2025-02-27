package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// Request structure for OpenRouter API
type OpenRouterRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Response structure for OpenRouter API
type OpenRouterResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func TemplateHandler(w http.ResponseWriter, r *http.Request) {
	// Create request body
	requestBody := OpenRouterRequest{
		Model: "meta-llama/llama-3.1-8b-instruct:free",
		Messages: []Message{
			{
				Role:    "user",
				Content: "what is 2+3",
			},
		},
	}

	// Convert request to JSON
	requestJSON, err := json.Marshal(requestBody)
	if err != nil {
		http.Error(w, "Failed to marshal request: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewBuffer(requestJSON))
	if err != nil {
		http.Error(w, "Failed to create request: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Add headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("OPENAI_API_KEY"))
	// Optional OpenRouter specific headers
	req.Header.Set("HTTP-Referer", "https://your-site.com") // Optional: your site URL
	req.Header.Set("X-Title", "My Application")             // Optional: your app name

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to send request: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read the response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		http.Error(w, fmt.Sprintf("API error (status %d): %s", resp.StatusCode, string(respBody)), http.StatusInternalServerError)
		return
	}

	// Parse the response
	var openRouterResp OpenRouterResponse
	err = json.Unmarshal(respBody, &openRouterResp)
	if err != nil {
		http.Error(w, "Failed to parse response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the response content
	if len(openRouterResp.Choices) > 0 {
		content := openRouterResp.Choices[0].Message.Content
		fmt.Println("Response:", content)

		// Return the response to the client
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": content})
	} else {
		http.Error(w, "No response choices returned", http.StatusInternalServerError)
	}
}
