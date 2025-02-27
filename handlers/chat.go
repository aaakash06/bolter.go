package handlers

import (
	"context"
	"net/http"
	"os"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

func Chat(w http.ResponseWriter, r *http.Request) {
	client := openai.NewClient(
		option.WithBaseURL("https://openrouter.ai/api/v1"),
		option.WithAPIKey(os.Getenv("OPEN_API_KEY")), // defaults to os.LookupEnv("OPENAI_API_KEY")
	)

	chatCompletion, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.UserMessage("What is 2+2?"),
		}),
		Model: openai.F("deepseek/deepseek-r1:free"),
	})
	if err != nil {
		panic(err.Error())
	}
	res := chatCompletion.Choices[0].Message.Content
	println(res)
}
