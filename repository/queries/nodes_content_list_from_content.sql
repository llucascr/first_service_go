SELECT
	node_id,
	content_id,
	user_id,
	notebook_id,
	created_at,
	updated_at
FROM nodes_contents
WHERE content_id = $1 AND deleted_at IS NULL
ORDER BY created_at DESC
