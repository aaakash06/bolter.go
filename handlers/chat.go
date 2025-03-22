package handlers

import (
	"bolter/utils"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
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
	body, err := ioutil.ReadAll(r.Body)
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


	
  // make api call
	client := openai.NewClient(
		option.WithAPIKey("OPENAI_API_KEY"), // defaults to os.LookupEnv("OPENAI_API_KEY")
	)
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
	chatCompletion, err := client.Chat.Completions.New(context.TODO(), param)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ChatErrorResponse{
			Error: err.Error(),
		})
		return
	}
	println(chatCompletion.Choices[0].Message.Content)

	// Return the response
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
// func Chat(w http.ResponseWriter, r *http.Request) {

// 	// extract the body from the request
// 	body, err := ioutil.ReadAll(r.Body)
// 	if err != nil {
// 		w.Header().Set("Content-Type", "application/json")
// 		w.WriteHeader(http.StatusInternalServerError)
// 		json.NewEncoder(w).Encode(ChatErrorResponse{
// 			Error: err.Error(),
// 		})
// 		// http.Error(w, "Error reading request body", http.StatusBadRequest)
// 		return
// 	}

// 	// defer the closing of the body
// 	defer r.Body.Close()

// 	var data ChatRequestBody

// 	// convert the json to struct
// 	if err := json.Unmarshal(body, &data); err != nil {
// 		w.Header().Set("Content-Type", "application/json")
// 		w.WriteHeader(http.StatusInternalServerError)
// 		json.NewEncoder(w).Encode(ChatErrorResponse{
// 			Error: err.Error(),
// 		})
// 		// http.Error(w, "Error parsing request body", http.StatusBadRequest)
// 		return
// 	}

// 	// Get the singleton client
// 	client := utils.GetOpenRouterClient()

// 	// Create messages
// 	messages := []utils.Message{
// 		utils.SystemMessage(utils.GetSystemPrompt("")),
// 	}

// 	for _, message := range data.Messages {
// 		messages = append(messages, message)
// 	}

// 	// Call the API
// 	resp, err := client.ChatCompletion("deepseek/deepseek-r1-zero:free", messages)
// 	if err != nil {
// 		w.Header().Set("Content-Type", "application/json")
// 		w.WriteHeader(http.StatusInternalServerError)
// 		json.NewEncoder(w).Encode(ChatErrorResponse{
// 			Error: err.Error(),
// 		})
// 		return
// 	}

// 	// Return the response
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)

// 	if len(resp.Choices) > 0 {
// 		stepsXML := resp.Choices[0].Message.Content
// 		json.NewEncoder(w).Encode(ChatResponse{
// 			Steps: stepsXML,
// 		})
// 	} else {
// 		json.NewEncoder(w).Encode(ChatErrorResponse{
// 			Error: "No response choices returned",
// 		})
// 	}
// }
