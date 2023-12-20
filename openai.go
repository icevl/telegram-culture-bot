package main

import (
	"context"
	"fmt"

	openai "github.com/sashabaranov/go-openai"
)

type OpenAI struct {
	Client *openai.Client
}

func NewOpenAI(token string) *OpenAI {
	return &OpenAI{
		Client: openai.NewClient(token),
	}
}

func (o *OpenAI) GetAnswer(message string) (string, bool) {
	resp, err := o.Client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: message,
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return "", false
	}

	resultText := resp.Choices[0].Message.Content
	response := escapeQuotes(resultText)

	return response, true
}
