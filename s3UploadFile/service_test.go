package s3uploadfile

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/OliPou/s3are/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUploadRequest(t *testing.T) {
	// Setup fixed UUID for testing
	fixedUUID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	mockUUIDGenerator := func() uuid.UUID {
		return fixedUUID
	}

	// Setup mock S3 client
	mockS3Client := &MockS3Client{
		GeneratePresignedURLFunc: func(key, contentType string) (string, error) {
			return "http://mock-presigned-url", nil
		},
	}

	// Setup mock DB
	mockDB := &MockDB{
		CreateUploadedFileFunc: func(ctx context.Context, arg database.CreateUploadedFileParams) (database.UploadedFile, error) {
			return database.UploadedFile{
				TransactionUuid:    fixedUUID,
				Consumer:           "test-consumer",
				UserName:           "test-user",
				FileName:           "test-file.txt",
				UploadPresignedUrl: "http://mock-presigned-url",
				Status:             "Waiting file",
				CreatedAt:          time.Now(),
			}, nil
		},
	}

	// Setup API config with mocks
	apiCfg := &ApiConfig{
		S3Client: mockS3Client,
		DB:       mockDB,
	}

	// Create gin context
	c, _ := gin.CreateTestContext(nil)

	// Test parameters
	params := UploadsFileParams{
		UserName:      "test-user",
		FileName:      "test-file",
		FileExtention: "txt",
	}

	// Execute test
	result, err := UploadRequest(c, params, "test-consumer", apiCfg, mockUUIDGenerator)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, fixedUUID, result.TransactionUuid)
	assert.Equal(t, "test-consumer", result.Consumer)
	assert.Equal(t, "test-user", result.UserName)
	assert.Equal(t, "Waiting file", result.Status)
}

func TestUploadedCompleted(t *testing.T) {
	// Setup fixed UUID for testing
	fixedUUID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	fileName := fixedUUID.String() + "_oly_filename.txt"
	// Setup mock DB
	// Setup mock S3 client
	mockS3Client := &MockS3Client{
		GeneratePresignedDownloadURLFunc: func(key string) (string, error) {
			return "http://mock-presigned-url", nil
		},
	}
	mockDB := &MockDB{
		UpdateUploadedFileFunc: func(ctx context.Context, arg database.UpdateUploadedFileParams) (database.UploadedFile, error) {
			return database.UploadedFile{
				TransactionUuid:      fixedUUID,
				Consumer:             "test-consumer",
				UserName:             "test-user",
				FileName:             "test-file.txt",
				FileSize:             sql.NullInt32{Int32: 1000, Valid: true},
				FileType:             sql.NullString{String: "text/plain", Valid: true},
				DownloadPresignedUrl: sql.NullString{String: "https://s3.download", Valid: true},
				Status:               "File Uploaded",
				CreatedAt:            time.Now(),
			}, nil
		},
	}

	// Setup API config with mocks
	apiCfg := &ApiConfig{
		DB:       mockDB,
		S3Client: mockS3Client,
	}

	// Create gin context
	c, _ := gin.CreateTestContext(nil)

	// Test parameters
	params := UploadCompletedParams{
		FileName: fileName,
		FileSize: 782,
		FileType: "text/plain",
	}

	// Execute test
	result, err := UploadedCompleted(c, params, apiCfg)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, fixedUUID, result.TransactionUuid)
	assert.Equal(t, "File Uploaded", result.Status)
	assert.Equal(t, int32(1000), result.FileSize.Int32)
	assert.Equal(t, "text/plain", result.FileType.String)
}

func TestGetUploadedFile(t *testing.T) {
	fixedUUID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")

	// Setup mock DB
	mockDB := &MockDB{
		GetUploadedFileFunc: func(ctx context.Context, arg database.GetUploadedFileParams) (database.UploadedFile, error) {
			return database.UploadedFile{
				TransactionUuid: fixedUUID,
				Consumer:        "test-consumer",
				UserName:        "test-user",
				FileName:        "test-file.txt",
				Status:          "File Uploaded",
				CreatedAt:       time.Now(),
			}, nil
		},
	}

	// Setup API config with mocks
	apiCfg := &ApiConfig{
		DB: mockDB,
	}

	// Create gin context
	c, _ := gin.CreateTestContext(nil)

	// Test parameters
	params := database.GetUploadedFileParams{
		TransactionUuid: fixedUUID,
		Consumer:        "test-consumer",
		UserName:        "test-user",
	}

	// Execute test
	result, err := apiCfg.DB.GetUploadedFile(c, params)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, fixedUUID, result.TransactionUuid)
	assert.Equal(t, "test-consumer", result.Consumer)
	assert.Equal(t, "test-user", result.UserName)
	assert.Equal(t, "File Uploaded", result.Status)
}
