package openai

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"

	"github.com/danielmesquitta/api-finance-manager/internal/config/env"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/gpt"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"golang.org/x/sync/errgroup"
)

type OpenAI struct {
	Client *openai.Client
}

func NewOpenAI(
	e *env.Env,
) *OpenAI {
	client := openai.NewClient(
		option.WithAPIKey(e.OpenAIAPIKey),
	)

	return &OpenAI{
		Client: &client,
	}
}

func (o *OpenAI) Completion(
	ctx context.Context,
	messages []gpt.Message,
	options ...gpt.Option,
) (*gpt.Message, error) {
	opts := gpt.Options{}
	for _, opt := range options {
		opt(&opts)
	}

	params := o.prepareParams(messages, opts)

	toolsByName := make(map[string]gpt.Tool)
	for _, tool := range opts.Tools {
		toolsByName[tool.Name] = tool
	}

	const maxAttempts = 3
	for range maxAttempts {
		completion, err := o.Client.Chat.Completions.New(ctx, params)
		if err != nil {
			return nil, errs.New(err)
		}
		if len(completion.Choices) == 0 {
			return nil, errs.New(fmt.Errorf("no choices returned"))
		}

		choice := completion.Choices[0].Message
		if len(choice.ToolCalls) == 0 {
			message := gpt.Message{
				Role:    gpt.RoleAssistant,
				Content: choice.Content,
			}
			return &message, nil
		}

		params.Messages = append(params.Messages, choice.ToParam())

		err = o.processToolCalls(ctx, choice.ToolCalls, toolsByName, &params)
		if err != nil {
			return nil, errs.New(err)
		}
	}

	return nil, errs.New("max attempts reached without a valid response")
}

func (o *OpenAI) processToolCalls(
	ctx context.Context,
	toolCalls []openai.ChatCompletionMessageToolCall,
	toolsByName map[string]gpt.Tool,
	params *openai.ChatCompletionNewParams,
) error {
	g, subCtx := errgroup.WithContext(ctx)
	mu := sync.Mutex{}

	for _, tc := range toolCalls {
		g.Go(func() error {
			tool, ok := toolsByName[tc.ID]
			if !ok {
				return fmt.Errorf("tool %s not found", tc.ID)
			}

			var args map[string]any
			if err := json.Unmarshal([]byte(tc.Function.Arguments), &args); err != nil {
				return fmt.Errorf(
					"failed to unmarshal tool %s args: %w",
					tc.ID,
					err,
				)
			}

			toolMessage, err := tool.Func(subCtx, args)
			if err != nil {
				return fmt.Errorf("tool %s failed: %w", tc.ID, err)
			}

			mu.Lock()
			params.Messages = append(
				params.Messages,
				openai.ToolMessage(toolMessage, tc.ID),
			)
			mu.Unlock()
			return nil
		})
	}

	return g.Wait()
}
func (o *OpenAI) prepareParams(
	messages []gpt.Message,
	opts gpt.Options,
) openai.ChatCompletionNewParams {
	openAIMessages := o.parseMessages(messages)
	openAITools := o.parseTools(opts.Tools)

	param := openai.ChatCompletionNewParams{
		Messages: openAIMessages,
		Tools:    openAITools,
	}
	if opts.Temperature > 0 {
		param.Temperature = openai.Float(opts.Temperature)
	}
	if opts.Seed > 0 {
		param.Seed = openai.Int(opts.Seed)
	}

	return param
}

func (o *OpenAI) parseTools(
	tools []gpt.Tool,
) []openai.ChatCompletionToolParam {
	openAITools := []openai.ChatCompletionToolParam{}
	for _, tool := range tools {
		openAITool := openai.ChatCompletionToolParam{
			Function: openai.FunctionDefinitionParam{
				Name:        tool.Name,
				Description: openai.String(tool.Description),
				Parameters:  tool.Args,
			},
		}

		openAITools = append(openAITools, openAITool)
	}

	return openAITools
}

func (o *OpenAI) parseMessages(
	messages []gpt.Message,
) []openai.ChatCompletionMessageParamUnion {
	openAIMessages := []openai.ChatCompletionMessageParamUnion{}
	for _, message := range messages {
		msg, err := o.parseMessage(message)
		if err != nil {
			slog.Error("failed to map gpt message", "error", err)
			continue
		}

		openAIMessages = append(
			openAIMessages,
			*msg,
		)
	}

	return openAIMessages
}

func (o *OpenAI) parseMessage(
	message gpt.Message,
) (*openai.ChatCompletionMessageParamUnion, error) {
	var msg openai.ChatCompletionMessageParamUnion
	switch message.Role {
	case gpt.RoleAssistant:
		msg = openai.AssistantMessage(message.Content)

	case gpt.RoleSystem:
		msg = openai.SystemMessage(message.Content)

	case gpt.RoleUser:
		msg = openai.UserMessage(message.Content)

	default:
		return nil, errs.New(fmt.Errorf("invalid role %s", message.Role))
	}

	return &msg, nil
}

var _ gpt.GPT = (*OpenAI)(nil)
