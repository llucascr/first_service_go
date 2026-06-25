SELECT
	content_id,
	tag_id,
	notebook_id,
	user_id,
	deleted_at,
	created_at,
	updated_at
FROM meta_tag_contents
WHERE content_id = $1 AND tag_id = $2 AND deleted_at IS NULL
