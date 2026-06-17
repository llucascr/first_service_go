UPDATE meta_contents
SET deleted_at = $1
WHERE content_id = $2 AND deleted_at IS NULL
