package openai

import (
	"context"

	ai "github.com/openai/openai-go"
	"github.com/openai/openai-go/azure"
	"github.com/openai/openai-go/option"
	config "github.com/shutils/lazyreview/pkg/config"
)

type Client struct {
	api  ai.Client
	conf config.Config
}

func NewClient(conf config.Config) Client {
	var api ai.Client
	if conf.Type == "azure" {
		azureOpenAIEndpoint := conf.Endpoint
		azureOpenAIAPIVersion := conf.Version
		azureOpenAIKey := conf.Key

		api = *ai.NewClient(
			azure.WithEndpoint(azureOpenAIEndpoint, azureOpenAIAPIVersion),
			azure.WithAPIKey(azureOpenAIKey),
		)
	} else {
		apiKey := conf.Key
		api = *ai.NewClient(
			option.WithAPIKey(apiKey),
		)
	}
	return Client{
		api:  api,
		conf: conf,
	}
}

// ChatGPT API呼び出し
func (c Client) Getreviewfromchatgpt(content string, conf config.Config) (*ai.ChatCompletion, error) {
	prompt := "you are a code reviewer. return the response in japanese."
	if conf.Prompt != "" {
		prompt = conf.Prompt
	}
	return reviewFromChatGPT(c, content, conf, prompt)
}

// ChatGPT API呼び出し
func (c Client) GetReviewFromChatGPTWithPrompt(content string, conf config.Config, prompt string) (*ai.ChatCompletion, error) {
	return reviewFromChatGPT(c, content, conf, prompt)
}

func reviewFromChatGPT(c Client, content string, conf config.Config, prompt string) (*ai.ChatCompletion, error) {
	maxTokens := 1000
	if conf.MaxTokens != 0 {
		maxTokens = conf.MaxTokens
	}
	review, err := c.api.Chat.Completions.New(context.TODO(), ai.ChatCompletionNewParams{
		Model: ai.F(c.conf.Model),
		Messages: ai.F([]ai.ChatCompletionMessageParamUnion{
			ai.SystemMessage(prompt),
			ai.UserMessage(content),
		}),
		MaxTokens: ai.Int(int64(maxTokens)),
	})

	if err != nil {
		return review, err
	}
	return review, nil
}
