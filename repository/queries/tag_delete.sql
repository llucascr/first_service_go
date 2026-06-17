UPDATE tags
SET deleted_at = $1
WHERE tag_id = $2 AND deleted_at IS NULL
