-- name: CreateUploadedFile :one
INSERT INTO uploaded_file (
    transaction_uuid,
    consumer,
    user_name,
    file_name,
    upload_presigned_url,
    status,
    created_at
) VALUES (
    $1, $2, $3, $4, $5, $6, NOW()
)
RETURNING *;

-- name: UpdateUploadedFile :one
UPDATE uploaded_file
SET
    file_size = $2,
    file_type = $3,
    download_presigned_url = $4,
    status = $5,
    updated_at = NOW()
WHERE transaction_uuid = $1
RETURNING *;

-- name: GetUploadedFile :one
SELECT * FROM uploaded_file
WHERE transaction_uuid = $1 and consumer = $2 and user_name = $3
LIMIT 1;
