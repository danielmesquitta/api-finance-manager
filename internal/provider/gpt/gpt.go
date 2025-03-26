package gpt

import "context"

type Role = string

const (
	RoleAssistant Role = "assistant"
	RoleSystem    Role = "system"
	RoleUser      Role = "user"
)

type Message struct {
	Role    Role   `json:"role"`
	Content string `json:"content"`
}

type ToolFunc func(ctx context.Context, args map[string]any) (string, error)

type Tool struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Func        ToolFunc
	Args        map[string]any `json:"args"`
}

type Options struct {
	Temperature float64 `json:"temperature"`
	Seed        int64   `json:"seed"`
	Tools       []Tool  `json:"tools"`
}

type Option func(*Options)

func WithTemperature(temperature float64) Option {
	return func(o *Options) {
		o.Temperature = temperature
	}
}

func WithSeed(seed int64) Option {
	return func(o *Options) {
		o.Seed = seed
	}
}

func WithTools(tools []Tool) Option {
	return func(o *Options) {
		o.Tools = tools
	}
}

type GPT interface {
	Completion(
		ctx context.Context,
		messages []Message,
		options ...Option,
	) (*Message, error)
}
