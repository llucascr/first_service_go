UPDATE nodes_contents
SET deleted_at = $1
WHERE node_id = $2 AND deleted_at IS NULL
