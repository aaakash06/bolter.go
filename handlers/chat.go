package handlers

import (
	"bolter/utils"
	"encoding/json"
	"io"
	"net/http"
)

// Response is a simple structure for JSON responses
type ChatResponse struct {
	Steps string `json:"steps"`
}

type ChatErrorResponse struct {
	Error string `json:"error,omitempty"`
}

type ChatRequestBody struct {
	Messages []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"messages"`
}


func Chat(w http.ResponseWriter, r *http.Request) {

	// extract the body from the request
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ChatErrorResponse{
			Error: err.Error(),
		})
		return
	}
	// defer the closing of the body
	defer r.Body.Close()

	var data ChatRequestBody

	// convert the json to struct
	if err := json.Unmarshal(body, &data); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ChatErrorResponse{
			Error: err.Error(),
		})
		// http.Error(w, "Error parsing request body", http.StatusBadRequest)
		return
	}

	// if invalid input, return error
	if len(data.Messages) == 0 || data.Messages == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ChatErrorResponse{
			Error: "Input format is invalid",
		})
		return
	}

	// Get the singleton client
	client := utils.GetOpenRouterClient()

	// Create messages
	messages := []utils.Message{
		utils.SystemMessage(utils.GetSystemPrompt("")),
	}

	for _, message := range data.Messages {
		messages = append(messages, message)
	}

	// Call the API
	resp, err := client.ChatCompletion("deepseek/deepseek-r1-zero:free", messages)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ChatErrorResponse{
			Error: err.Error(),
		})
		return
	}

	// Return the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if len(resp.Choices) > 0 {
		stepsXML := resp.Choices[0].Message.Content
		json.NewEncoder(w).Encode(ChatResponse{
			Steps: stepsXML,
		})
	} else {
		json.NewEncoder(w).Encode(ChatErrorResponse{
			Error: "No response choices returned",
		})
	}
}

