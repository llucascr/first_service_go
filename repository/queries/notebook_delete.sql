UPDATE notebooks
SET deleted_at = $1
WHERE notebook_id = $2 AND deleted_at IS NULL
