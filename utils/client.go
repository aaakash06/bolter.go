package utils

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
)

// OpenRouterClient is a singleton wrapper for making API calls to OpenRouter
type OpenRouterClient struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
	referer    string
	appTitle   string
}

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest represents the request structure for chat completions
type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream,omitempty"`
	// Add other parameters as needed (temperature, max_tokens, etc.)
}

// ChatResponse represents the non-streaming response structure
type ChatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// StreamChunk represents a chunk in a streaming response
type StreamChunk struct {
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
		FinishReason *string `json:"finish_reason"`
	} `json:"choices"`
}

// StreamHandler is a function type for handling streaming responses
type StreamHandler func(content string, finishReason *string) error

// Singleton instance
var (
	instance *OpenRouterClient
	once     sync.Once
)

// GetOpenRouterClient returns the singleton instance of OpenRouterClient
func GetOpenRouterClient() *OpenRouterClient {
	once.Do(func() {
		instance = &OpenRouterClient{
			apiKey:     os.Getenv("OPEN_API_KEY"),
			baseURL:    "https://openrouter.ai/api/v1",
			httpClient: &http.Client{},
			referer:    "https://your-site.com", // Change to your site
			appTitle:   "My Application",        // Change to your app name
		}
	})
	return instance
}

// SetAPIKey allows changing the API key
func (c *OpenRouterClient) SetAPIKey(apiKey string) {
	c.apiKey = apiKey
}

// SetReferer sets the HTTP Referer header
func (c *OpenRouterClient) SetReferer(referer string) {
	c.referer = referer
}

// SetAppTitle sets the X-Title header
func (c *OpenRouterClient) SetAppTitle(title string) {
	c.appTitle = title
}

// createRequest creates a new HTTP request with appropriate headers
func (c *OpenRouterClient) createRequest(method, endpoint string, body []byte) (*http.Request, error) {
	req, err := http.NewRequest(method, c.baseURL+endpoint, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("HTTP-Referer", c.referer)
	req.Header.Set("X-Title", c.appTitle)

	return req, nil
}

// ChatCompletion sends a request to the chat completions endpoint and returns the response
func (c *OpenRouterClient) ChatCompletion(model string, messages []Message) (*ChatResponse, error) {
	// Create request body
	reqBody := ChatRequest{
		Model:    model,
		Messages: messages,
		Stream:   false,
	}

	// Convert to JSON
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	req, err := c.createRequest("POST", "/chat/completions", jsonBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Send request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	// Parse response
	var chatResp ChatResponse
	if err := json.Unmarshal(respBody, &chatResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &chatResp, nil
}

// StreamChatCompletion streams the response from the chat completions endpoint
func (c *OpenRouterClient) StreamChatCompletion(model string, messages []Message, handler StreamHandler) error {
	// Create request body
	reqBody := ChatRequest{
		Model:    model,
		Messages: messages,
		Stream:   true,
	}

	// Convert to JSON
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	req, err := c.createRequest("POST", "/chat/completions", jsonBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Send request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	// Process the streaming response
	reader := bufio.NewReader(resp.Body)
	for {
		// Read a line
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("error reading stream: %w", err)
		}

		// Skip empty lines
		line = bytes.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		// Check for data prefix
		if bytes.HasPrefix(line, []byte("data: ")) {
			data := bytes.TrimPrefix(line, []byte("data: "))

			// Check if done
			if string(data) == "[DONE]" {
				break
			}

			// Parse the chunk
			var chunk StreamChunk
			if err := json.Unmarshal(data, &chunk); err != nil {
				return fmt.Errorf("failed to parse chunk: %w", err)
			}

			// Process content
			if len(chunk.Choices) > 0 {
				content := chunk.Choices[0].Delta.Content
				finishReason := chunk.Choices[0].FinishReason

				// Call the handler
				if err := handler(content, finishReason); err != nil {
					return fmt.Errorf("handler error: %w", err)
				}
			}
		}
	}

	return nil
}

// Simple helper to create a user message
func UserMessage(content string) Message {
	return Message{
		Role:    "user",
		Content: content,
	}
}

// Simple helper to create an assistant message
func AssistantMessage(content string) Message {
	return Message{
		Role:    "assistant",
		Content: content,
	}
}

// Simple helper to create a system message
func SystemMessage(content string) Message {
	return Message{
		Role:    "system",
		Content: content,
	}
}
