UPDATE meta_tag_contents
SET deleted_at = $1
WHERE content_id = $2 AND tag_id = $3 AND deleted_at IS NULL
