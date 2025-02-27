-- name: CreateFeedback :exec
INSERT INTO feedbacks (message, user_id)
VALUES ($1, $2);