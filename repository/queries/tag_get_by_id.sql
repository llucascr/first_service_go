SELECT
	tag_id,
	name,
	color,
	user_id,
	deleted_at,
	created_at,
	updated_at
FROM tags
WHERE tag_id = $1 AND deleted_at IS NULL
