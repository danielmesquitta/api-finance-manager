package entity

type FullAIChat struct {
	AIChat
	HasMessages bool `json:"has_messages"`
}
