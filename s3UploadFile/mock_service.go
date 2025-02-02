package s3uploadfile

import (
	"context"
	"time"

	"github.com/OliPou/s3are/internal/database"
)

// Mock S3 Client
type MockS3Client struct {
	GeneratePresignedURLFunc         func(key string, expirationTime *int) (string, time.Duration, error)
	GeneratePresignedDownloadURLFunc func(key string, expirationTime *int) (string, time.Duration, error)
}

func (m *MockS3Client) GeneratePresignedURL(key string, expirationTime *int) (string, time.Duration, error) {
	return m.GeneratePresignedURLFunc(key, nil)
}

func (m *MockS3Client) GeneratePresignedDownloadURL(key string, expirationTime *int) (string, time.Duration, error) {
	return m.GeneratePresignedDownloadURLFunc(key, nil)
}

// Mock DB
type MockDB struct {
	CreateUploadedFileFunc func(ctx context.Context, arg database.CreateUploadedFileParams) (database.UploadedFile, error)
	UpdateUploadedFileFunc func(ctx context.Context, arg database.UpdateUploadedFileParams) (database.UploadedFile, error)
	GetUploadedFileFunc    func(ctx context.Context, arg database.GetUploadedFileParams) (database.UploadedFile, error)
}

func (m *MockDB) CreateUploadedFile(ctx context.Context, arg database.CreateUploadedFileParams) (database.UploadedFile, error) {
	return m.CreateUploadedFileFunc(ctx, arg)
}

func (m *MockDB) UpdateUploadedFile(ctx context.Context, arg database.UpdateUploadedFileParams) (database.UploadedFile, error) {
	return m.UpdateUploadedFileFunc(ctx, arg)
}

func (m *MockDB) GetUploadedFile(ctx context.Context, arg database.GetUploadedFileParams) (database.UploadedFile, error) {
	return m.GetUploadedFileFunc(ctx, arg)
}

// Verify that MockDB implements DBInterface
var _ DBInterface = (*MockDB)(nil)

// Verify that MockS3Client implements S3ClientInterface
var _ S3ClientInterface = (*MockS3Client)(nil)
