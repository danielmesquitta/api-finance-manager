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
	ID        uuid.UUID           `db:"id"         json:"id,omitempty"`
	Message   string              `db:"message"    json:"message,omitempty"`
	Rating    *string             `db:"rating"     json:"rating,omitempty"`
	Author    AIChatMessageAuthor `db:"author"     json:"author,omitempty"`
	CreatedAt time.Time           `db:"created_at" json:"created_at,omitempty"`
}
