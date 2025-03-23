package handlers

import (
	"bolter/utils"
	"encoding/json"
	"io"
	"net/http"
)

// Response is a simple structure for JSON responses
type TemplateResponse struct {
	Prompts   []string `json:"promts"`
	UiPrompts []string `json:"uiPrompts"`
}

type TemplateErrorResponse struct {
	Error string `json:"error,omitempty"`
}

type TemplateRequestBody struct {
	Prompt string `json:"prompt"`
}

func TemplateHandler(w http.ResponseWriter, r *http.Request) {

	// extract the body from the request
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(TemplateErrorResponse{
			Error: err.Error(),
		})
		// http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	// defer the closing of the body
	defer r.Body.Close()

	var data TemplateRequestBody
	// var data map[string]interface{}

	// convert the json to struct
	if err := json.Unmarshal(body, &data); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(TemplateErrorResponse{
			Error: err.Error(),
		})
		// http.Error(w, "Error parsing request body", http.StatusBadRequest)
		return
	}

	// Get the singleton client
	client := utils.GetOpenRouterClient()

	// Create messages
	messages := []utils.Message{
		utils.SystemMessage("Return either node or react based on what do you think this project should be. Only return a single word either 'node' or 'react'. Do not return anything extra."),
		utils.UserMessage(data.Prompt),
	}

	// Call the API
	resp, err := client.ChatCompletion("meta-llama/llama-3.1-8b-instruct:free", messages)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(TemplateErrorResponse{
			Error: err.Error(),
		})
		return
	}

	// Return the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if len(resp.Choices) > 0 {
		tech := resp.Choices[0].Message.Content
		if tech == "react" {
			json.NewEncoder(w).Encode(TemplateResponse{
				Prompts:   []string{utils.BasePrompt, utils.GetFSPrompt(tech)},
				UiPrompts: []string{utils.GetTechStackPrompt(tech)},
			})
		} else if tech == "node" {
			json.NewEncoder(w).Encode(TemplateResponse{
				Prompts:   []string{utils.GetFSPrompt(tech)},
				UiPrompts: []string{utils.GetTechStackPrompt(tech)},
			})
		}
	} else {
		json.NewEncoder(w).Encode(TemplateErrorResponse{
			Error: "No response choices returned",
		})
	}
}
