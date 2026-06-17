SELECT
	tag_id,
	name,
	color,
	user_id,
	created_at,
	updated_at
FROM tags
WHERE user_id = $1 AND deleted_at IS NULL
ORDER BY created_at DESC
