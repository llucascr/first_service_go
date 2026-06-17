UPDATE tags
SET
	name = $1,
	color = $2,
	updated_at = $3
WHERE tag_id = $4 AND deleted_at IS NULL
RETURNING
	tag_id,
	name,
	color,
	user_id,
	deleted_at,
	created_at,
	updated_at
