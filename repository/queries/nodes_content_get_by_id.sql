SELECT
	node_id,
	content_id,
	user_id,
	notebook_id,
	deleted_at,
	created_at,
	updated_at
FROM nodes_contents
WHERE node_id = $1 AND deleted_at IS NULL
