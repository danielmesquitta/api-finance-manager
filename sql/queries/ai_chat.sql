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
-- name: GetAIChatByID :one
SELECT *
FROM ai_chats
WHERE id = $1
  AND deleted_at IS NULL;
-- name: ListAIChatMessagesAndAnswers :many
WITH combined_messages AS (
  SELECT m.id,
    m.message,
    NULL as rating,
    'USER' as author,
    m.created_at
  FROM ai_chat_messages m
  WHERE m.ai_chat_id = $1
    AND m.deleted_at IS NULL
  UNION ALL
  SELECT r.id,
    r.message,
    r.rating,
    'AI' as author,
    r.created_at
  FROM ai_chat_answers r
    JOIN ai_chat_messages m ON m.id = r.ai_chat_message_id
  WHERE m.ai_chat_id = $1
    AND m.deleted_at IS NULL
    AND r.deleted_at IS NULL
)
SELECT *
FROM combined_messages
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;
-- name: CountAIChatMessagesAndAnswers :one
SELECT (
    SELECT COUNT(*)
    FROM ai_chat_messages m
    WHERE m.ai_chat_id = $1
      AND m.deleted_at IS NULL
  ) + (
    SELECT COUNT(*)
    FROM ai_chat_answers r
      JOIN ai_chat_messages m ON m.id = r.ai_chat_message_id
    WHERE m.ai_chat_id = $1
      AND m.deleted_at IS NULL
      AND r.deleted_at IS NULL
  ) AS total_count;