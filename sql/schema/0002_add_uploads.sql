-- +goose Up
ALTER TABLE uploaded_file
ADD download_expiration_time TIMESTAMP;
ALTER TABLE uploaded_file
ADD upload_expiration_time TIMESTAMP;

-- +goose Down
ALTER TABLE uploaded_file
DROP COLUMN download_expiration_time;
ALTER TABLE uploaded_file
DROP COLUMN upload_expiration_time;
