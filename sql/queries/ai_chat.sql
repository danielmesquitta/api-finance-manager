-- name: GetLatestAIChatByUserID :one
SELECT ai_chats.*,
  EXISTS (
    SELECT 1
    FROM ai_chat_messages
    WHERE ai_chat_messages.ai_chat_id = ai_chats.id
      AND ai_chat_messages.deleted_at IS NULL
    LIMIT 1
  ) AS has_messages
FROM ai_chats
WHERE ai_chats.deleted_at IS NULL
  AND ai_chats.user_id = $1
ORDER BY ai_chats.created_at DESC
LIMIT 1;
-- name: CreateAIChat :one
INSERT INTO ai_chats (user_id)
VALUES ($1)
RETURNING *;
-- name: UpdateAIChat :exec
UPDATE ai_chats
SET title = $2
WHERE id = $1;
-- name: DeleteAIChat :exec
UPDATE ai_chats
SET deleted_at = now()
WHERE id = $1;
-- name: GetAIChat :one
SELECT *
FROM ai_chats
WHERE id = $1;
