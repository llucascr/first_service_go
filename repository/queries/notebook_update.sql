UPDATE notebooks
SET
	icon = $1,
	name = $2,
	image = $3,
	description = $4,
	updated_at = $5
WHERE notebook_id = $6 AND deleted_at IS NULL
RETURNING
	notebook_id,
	user_id,
	icon,
	name,
	image,
	description,
	deleted_at,
	created_at,
	updated_at
