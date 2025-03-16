package entity

import (
	"time"

	"github.com/google/uuid"
)

type FullAIChat struct {
	AIChat
	HasMessages bool `json:"has_messages"`
}

type AIChatMessageAuthor = string

const (
	AIChatMessageAuthorUser AIChatMessageAuthor = "USER"
	AIChatMessageAuthorAI   AIChatMessageAuthor = "ARTIFICIAL_INTELLIGENCE"
)

type AIChatMessageAndAnswer struct {
	ID        uuid.UUID           `json:"id,omitempty"`
	Message   string              `json:"message,omitempty"`
	Rating    *string             `json:"rating,omitempty"`
	Author    AIChatMessageAuthor `json:"author,omitempty"`
	CreatedAt time.Time           `json:"created_at,omitempty"`
}
