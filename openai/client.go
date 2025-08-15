package openai

import (
	"context"
	"log"
	"spacedrepetitiongo/config"

	openai "github.com/sashabaranov/go-openai"
)

type OpenAiClient struct {
	Client *openai.Client
}

func NewClient() *OpenAiClient {
	return &OpenAiClient{
		Client: openai.NewClient(
			config.OpenAiApiKey(),
		),
	}
}

func (client *OpenAiClient) Ask(
	systemPrompt string,
	messagePrompt string,
) string {
	resp, err := client.Client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT5,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: systemPrompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: messagePrompt,
				},
			},
		},
	)
	if err != nil {
		log.Printf("Coult not generate text: %v\n", err)
	}
	return resp.Choices[0].Message.Content
}

func (client *OpenAiClient) CreateImage(prompt string) string {
	respUrl, err := client.Client.CreateImage(
		context.Background(),
		openai.ImageRequest{
			Prompt:         prompt,
			Size:           openai.CreateImageSize256x256,
			ResponseFormat: openai.CreateImageResponseFormatURL,
			OutputFormat:   openai.CreateImageOutputFormatJPEG,
			N:              1,
		},
	)
	if err != nil {
		log.Printf("Could not generate image: %v\n", err)
	}
	return respUrl.Data[0].URL
}
