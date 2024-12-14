-- name: ListCategories :many
SELECT *
FROM categories;
-- name: CreateManyCategories :copyfrom
INSERT INTO categories (external_id, name)
VALUES ($1, $2);
-- name: SearchCategories :many
SELECT *,
  levenshtein(unaccent(name), unaccent(@search)) AS distance
FROM categories
WHERE (
    @search = ''
    OR levenshtein(unaccent(name), unaccent(@search)) <= 2
  )
ORDER BY distance,
  name
LIMIT $1 OFFSET $2;
-- name: CountSearchCategories :one
SELECT COUNT(*)
FROM categories
WHERE (
    @search::string = ''
    OR levenshtein(unaccent(name), unaccent(@search::string)) <= 2
  );