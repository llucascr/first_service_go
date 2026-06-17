SELECT
	content_id,
	notebook_id,
	user_id,
	icon,
	name,
	created_at,
	updated_at
FROM meta_contents
WHERE notebook_id = $1 AND deleted_at IS NULL
ORDER BY created_at DESC
