package handlers

import (
	"bolter/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

// Response is a simple structure for JSON responses
type Response struct {
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

func TemplateHandler(w http.ResponseWriter, r *http.Request) {
	// Get the singleton client
	client := utils.GetOpenRouterClient()

	// Create messages
	messages := []utils.Message{
		// utils.SystemMessage("You are a helpful assistant."),
		utils.UserMessage("Return either node or react based on what do you think this project should be. Only return a single word either 'node' or 'react'. Do not return anything extra."),
	}

	// Call the API
	resp, err := client.ChatCompletion("meta-llama/llama-3.1-8b-instruct:free", messages)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Error: err.Error(),
		})
		return
	}

	// Return the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if len(resp.Choices) > 0 {
		json.NewEncoder(w).Encode(Response{
			Message: resp.Choices[0].Message.Content,
		})
	} else {
		json.NewEncoder(w).Encode(Response{
			Error: "No response choices returned",
		})
	}
}

// StreamingHandler handles streaming chat completions
func StreamingHandler(w http.ResponseWriter, r *http.Request) {
	// Set headers for streaming
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Ensure we can flush
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	// Get the singleton client
	client := utils.GetOpenRouterClient()

	// Create messages
	messages := []utils.Message{
		utils.SystemMessage("You are a helpful assistant."),
		utils.UserMessage("Explain quantum computing in simple terms."),
	}

	// Stream handler function
	streamHandler := func(content string, finishReason *string) error {
		// Send the content
		if content != "" {
			data, _ := json.Marshal(map[string]string{"content": content})
			_, err := fmt.Fprintf(w, "data: %s\n\n", data)
			if err != nil {
				return err
			}
		}

		// Handle finish reason
		if finishReason != nil {
			data, _ := json.Marshal(map[string]string{"finish_reason": *finishReason})
			_, err := fmt.Fprintf(w, "data: %s\n\n", data)
			if err != nil {
				return err
			}
		}

		flusher.Flush()
		return nil
	}

	// Stream the response
	err := client.StreamChatCompletion("meta-llama/llama-3.1-8b-instruct:free", messages, streamHandler)
	if err != nil {
		// Send error as an event
		errData, _ := json.Marshal(map[string]string{"error": err.Error()})
		fmt.Fprintf(w, "data: %s\n\n", errData)
		flusher.Flush()
	}

	// Send DONE event
	fmt.Fprint(w, "data: [DONE]\n\n")
	flusher.Flush()
}
