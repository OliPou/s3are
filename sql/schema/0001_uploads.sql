-- +goose Up
CREATE TABLE uploaded_file(
    transaction_uuid UUID PRIMARY KEY,
    consumer TEXT NOT NULL,
    user_name TEXT NOT NULL,
    file_name TEXT NOT NULL,
    file_size INT,
    file_type TEXT,
    upload_presigned_url TEXT NOT NULL,
    download_presigned_url TEXT,
    status TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP
);

-- +goose Down
DROP TABLE uploaded_file;
