UPDATE meta_contents
SET
	icon = $1,
	name = $2,
	updated_at = $3
WHERE content_id = $4 AND deleted_at IS NULL
RETURNING
	content_id,
	notebook_id,
	user_id,
	icon,
	name,
	deleted_at,
	created_at,
	updated_at
