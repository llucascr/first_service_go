UPDATE meta_tag_contents
SET
	notebook_id = $1,
	updated_at = $2
WHERE content_id = $3 AND tag_id = $4 AND deleted_at IS NULL
RETURNING
	content_id,
	tag_id,
	notebook_id,
	user_id,
	deleted_at,
	created_at,
	updated_at
