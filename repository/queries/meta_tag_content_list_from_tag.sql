SELECT
	content_id,
	tag_id,
	notebook_id,
	user_id,
	created_at,
	updated_at
FROM meta_tag_contents
WHERE tag_id = $1 AND deleted_at IS NULL
ORDER BY created_at DESC
