package s3uploadfile

import (
	"context"

	"github.com/OliPou/s3are/internal/database"
)

// Mock S3 Client
type MockS3Client struct {
	GeneratePresignedURLFunc         func(key, contentType string) (string, error)
	GeneratePresignedDownloadURLFunc func(key string) (string, error)
}

func (m *MockS3Client) GeneratePresignedURL(key, contentType string) (string, error) {
	return m.GeneratePresignedURLFunc(key, contentType)
}

func (m *MockS3Client) GeneratePresignedDownloadURL(key string) (string, error) {
	return m.GeneratePresignedDownloadURLFunc(key)
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
