SELECT * FROM users
WHERE name = $1 AND deleted_at IS NULL