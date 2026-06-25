SELECT
	content_id,
	notebook_id,
	user_id,
	icon,
	name,
	deleted_at,
	created_at,
	updated_at
FROM meta_contents
WHERE content_id = $1 AND deleted_at IS NULL
