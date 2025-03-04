package openai

import (
	"github.com/danielmesquitta/api-finance-manager/internal/config"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type OpenAI struct {
	Client *openai.Client
}

func NewOpenAI(
	e *config.Env,
) *OpenAI {
	client := openai.NewClient(
		option.WithAPIKey(e.OpenAIAPIKey),
	)

	return &OpenAI{
		Client: client,
	}
}
