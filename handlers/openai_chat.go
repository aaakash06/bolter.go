package handlers

import (
	"bolter/utils"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

func OpenAIChat(w http.ResponseWriter, r *http.Request) {
	// extract the body from the request
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ChatErrorResponse{
			Error: err.Error(),
		})
		// http.Error(w, "Error reading request body", http.StatusBadRequest)
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

	// create client
	client := openai.NewClient(
		option.WithAPIKey(os.Getenv("OPENAI_API_KEY")),
	)

	// create params
	param := openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{},
		Model: openai.ChatModelGPT4o,
		MaxCompletionTokens: openai.Int(1000),
	}
	param.Messages = append(param.Messages, openai.SystemMessage(utils.GetSystemPrompt("")))
	for _, message := range data.Messages {
		if message.Role == "user" {
			param.Messages = append(param.Messages, openai.UserMessage(message.Content))
		} else {
			param.Messages = append(param.Messages, openai.AssistantMessage(message.Content))
		}
	}

	// json.NewEncoder(w).Encode(param)

	// // make api call
	chatCompletion, err := client.Chat.Completions.New(context.TODO(), param)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ChatErrorResponse{
			Error: err.Error(),
		})
		return
	}

	// // println(chatCompletion.Choices[0].Message.Content)

	// // Return the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if len(chatCompletion.Choices) > 0 {
		stepsXML := chatCompletion.Choices[0].Message.Content
		json.NewEncoder(w).Encode(ChatResponse{
			Steps: stepsXML,
		})
	} else {
		json.NewEncoder(w).Encode(ChatErrorResponse{
			Error: "No response choices returned",
		})
	}
}