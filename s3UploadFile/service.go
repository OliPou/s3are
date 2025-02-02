package s3uploadfile

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/OliPou/s3are/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UUIDGenerator func() uuid.UUID

func UploadRequest(c *gin.Context, params UploadsFileParams, consumer string, apiCfg *ApiConfig, generateUUID UUIDGenerator) (UploadedFile, error) {
	var presignedURL string
	var duration time.Duration
	var err error
	transactionUUID := generateUUID()
	fileName := fmt.Sprintf("%s_%s_%s_%s.%s",
		transactionUUID.String(),
		consumer,
		params.UserName,
		params.FileName,
		params.FileExtention,
	)
	// if params.LinkExpirationDuration > 0 {
	presignedURL, duration, err = apiCfg.S3Client.GeneratePresignedURL(fileName, params.LinkExpirationDuration)
	if err != nil {
		fmt.Printf("error generating presigned URL: %v", err)
		return UploadedFile{}, fmt.Errorf("error generating presigned URL")
	}
	// } else {
	// 	presignedURL, duration, err = apiCfg.S3Client.GeneratePresignedURL(fileName, nil)
	// 	if err != nil {
	// 		fmt.Printf("error generating presigned URL: %v", err)
	// 		return UploadedFile{}, fmt.Errorf("error generating presigned URL")
	// 	}
	// }
	expirationTime := sql.NullTime{
		Time:  time.Now().Add(duration),
		Valid: true,
	}
	uploadedFile, err := apiCfg.DB.CreateUploadedFile(c, database.CreateUploadedFileParams{
		TransactionUuid:      transactionUUID,
		Consumer:             consumer,
		UserName:             params.UserName,
		FileName:             fileName,
		UploadPresignedUrl:   presignedURL,
		UploadExpirationTime: expirationTime,
		Status:               "Waiting file",
	})
	if err != nil {
		fmt.Printf("Error creating uploaded file: %v", err)
		return UploadedFile{}, fmt.Errorf("error creating uploaded file")
	}
	return DatabaseUploadFileToUploadFile(uploadedFile), nil
}

func UploadedCompleted(c *gin.Context, params UploadCompletedParams, apiCfg *ApiConfig) (UploadedFile, error) {
	presignedURL, duration, _ := apiCfg.S3Client.GeneratePresignedDownloadURL(string(params.FileName), nil)
	expirationTime := sql.NullTime{
		Time:  time.Now().Add(duration),
		Valid: true,
	}
	inputFileName := strings.Split(params.FileName, "_")
	uuidStr := inputFileName[0]
	transactionUuid, err := uuid.Parse(uuidStr)
	if err != nil {
		fmt.Printf("Error updating uploaded file: %v", err)
		return UploadedFile{}, fmt.Errorf("error updating uploaded file")
	}
	uploadedFile, err := apiCfg.DB.UpdateUploadedFile(c, database.UpdateUploadedFileParams{
		TransactionUuid: transactionUuid,
		FileSize: sql.NullInt32{
			Int32: int32(params.FileSize),
			Valid: true,
		},
		FileType: sql.NullString{
			String: params.FileType,
			Valid:  true,
		},
		DownloadPresignedUrl: sql.NullString{
			String: presignedURL,
			Valid:  true,
		},
		Status:                 "File Uploaded",
		DownloadExpirationTime: expirationTime,
	})
	if err != nil {
		fmt.Println("Error updating uploaded file:", err)
		return UploadedFile{}, fmt.Errorf("error updating uploaded file: %w", err)
	}
	return DatabaseUploadFileToUploadFile(uploadedFile), nil
}
