UPDATE nodes_contents
SET
	content_id = $1,
	updated_at = $2
WHERE node_id = $3 AND deleted_at IS NULL
RETURNING
	node_id,
	content_id,
	user_id,
	notebook_id,
	deleted_at,
	created_at,
	updated_at
