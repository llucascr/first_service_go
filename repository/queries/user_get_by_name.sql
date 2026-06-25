SELECT user_id, name, email, password, created_at, updated_at, deleted_at
FROM users
WHERE name = $1 AND deleted_at IS NULL