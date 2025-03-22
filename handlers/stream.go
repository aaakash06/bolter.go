package handlers

import (
	"bolter/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

// StreamingHandler handles streaming chat completions
func StreamingHandlerFunction(w http.ResponseWriter, r *http.Request) {
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
		utils.UserMessage("Explain quantum computing in simple terms. in 20 words or less."),
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
