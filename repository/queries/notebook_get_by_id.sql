SELECT *
FROM notebooks
WHERE notebook_id = $1
AND deleted_at IS NULL;