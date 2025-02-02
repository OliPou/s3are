package s3uploadfile

import (
	"context"
	"time"

	"github.com/OliPou/s3are/internal/database"
)

type DBInterface interface {
	CreateUploadedFile(context.Context, database.CreateUploadedFileParams) (database.UploadedFile, error)
	UpdateUploadedFile(context.Context, database.UpdateUploadedFileParams) (database.UploadedFile, error)
	GetUploadedFile(context.Context, database.GetUploadedFileParams) (database.UploadedFile, error)
}

type S3ClientInterface interface {
	GeneratePresignedURL(key string, expirationTime *int) (string, time.Duration, error)
	GeneratePresignedDownloadURL(key string, expirationTime *int) (string, time.Duration, error)
}
