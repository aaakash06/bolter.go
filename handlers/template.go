package handlers

import (
	"context"
	"net/http"
	"os"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

func TemplateHandler(w http.ResponseWriter, r *http.Request) {

	// init the client with baseURL and api key
	client := openai.NewClient(
		option.WithBaseURL("https://openrouter.ai/api/v1"),
		option.WithAPIKey(os.Getenv("OPEN_API_KEY")),
	)

	// make the call with empty context
	// chatCompletion, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
	// 	Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
	// 		openai.UserMessage("What is 2+2?"),
	// 	}),
	// 	Model: openai.F("deepseek/deepseek-r1:free"),
	// })

	// Define the chat completion request
	chatCompletion, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Model: openai.F("meta-llama/llama-3.1-8b-instruct:free"), // Example OpenRouter model
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.UserMessage("what is 2+3"),
		}),
	})
	if err != nil {
		panic(err.Error())
	}

	res := chatCompletion.Choices[0].Message.Content
	println(res)

	// w.WriteHeader(http.StatusOK)
	// w.Header().Set("Content-Type", "application/json")
	// w.Write([]byte(`{"message": "This is a  project"}`))
}

// func TemplateHandler(w http.ResponseWriter, r *http.Request) {
// 	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
// 	resp, err := client.CreateChatCompletion(
// 		context.Background(),
// 		openai.ChatCompletionRequest{
// 			Model: "deepseek/deepseek-r1:free",
// 			Messages: []openai.ChatCompletionMessage{
// 				{
// 					Role:    openai.ChatMessageRoleUser,
// 					Content: "Hello!",
// 				},
// 			},
// 		},
// 	)

// 	if err != nil {
// 		fmt.Printf("ChatCompletion error: %v\n", err)
// 		return
// 	}

// 	fmt.Println(resp.Choices[0].Message.Content)

// 	w.WriteHeader(http.StatusOK)
// 	w.Header().Set("Content-Type", "application/json")
// 	w.Write([]byte(`{"message": "This is a  project"}`))
// }
