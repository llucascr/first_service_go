INSERT INTO meta_tag_contents (content_id, tag_id, notebook_id, user_id, deleted_at, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
ON CONFLICT (content_id, tag_id) DO UPDATE
SET
	notebook_id = EXCLUDED.notebook_id,
	user_id = EXCLUDED.user_id,
	deleted_at = NULL,
	updated_at = EXCLUDED.updated_at
WHERE meta_tag_contents.deleted_at IS NOT NULL
RETURNING
	content_id,
	tag_id,
	notebook_id,
	user_id,
	deleted_at,
	created_at,
	updated_at
