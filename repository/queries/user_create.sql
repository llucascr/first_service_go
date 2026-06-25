INSERT INTO users (
    user_id,
    name,
    email,
    password,
    deleted_at,
    created_at,
    updated_at
) VALUES ($1, $2, $3, $4, $5, $6, $7)