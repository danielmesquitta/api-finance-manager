-- name: GenerateAIChatMessage :one
INSERT INTO ai_chat_messages (message, ai_chat_id)
VALUES ($1, $2)
RETURNING *;
-- name: DeleteAIChatMessages :exec
UPDATE ai_chat_messages
SET deleted_at = NOW()
WHERE ai_chat_id = $1;