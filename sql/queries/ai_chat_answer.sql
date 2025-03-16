-- name: CreateAIChatAnswer :one
INSERT INTO ai_chat_answers (message, ai_chat_message_id)
VALUES ($1, $2)
RETURNING *;
-- name: DeleteAIChatAnswers :exec
UPDATE ai_chat_answers acmr
SET acmr.deleted_at = NOW()
FROM ai_chat_messages acm
WHERE acmr.ai_chat_message_id = acm.id
  AND acm.ai_chat_id = $1;
-- name: UpdateAIChatAnswer :exec
UPDATE ai_chat_answers
SET rating = $2
WHERE id = $1;