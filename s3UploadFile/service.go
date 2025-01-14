package s3uploadfile

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/OliPou/s3are/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func UploadRequest(c *gin.Context, params UploadsFileParams, consumer string, apiCfg *ApiConfig) (UploadedFile, error) {
	transactionUUID := uuid.New()
	fileName := fmt.Sprintf("%s_%s_%s.%s",
		transactionUUID.String(),
		consumer,
		params.UserName,
		params.FileExtention,
	)
	presignedURL, err := apiCfg.S3Client.GeneratePresignedURL(fileName, "application/octet-stream")
	if err != nil {
		log.Fatal("Error generating presigned URL:", err)
		return UploadedFile{}, fmt.Errorf("Error generating presigned URL")
	}
	uploadedFile, err := apiCfg.DB.CreateUploadedFile(c, database.CreateUploadedFileParams{
		TransactionUuid:    transactionUUID,
		Consumer:           consumer,
		UserName:           params.UserName,
		FileName:           fileName,
		UploadPresignedUrl: presignedURL,
		Status:             "Waiting file",
		ContentType:        "application/octet-stream",
		// Add other required fields here
	})
	if err != nil {
		log.Fatal("Error creating uploaded file:", err)
		return UploadedFile{}, fmt.Errorf("Error creating uploaded file: %w", err)
	}
	return DatabaseUploadFileToUploadFile(uploadedFile), nil

}

func UploadedCompleted(c *gin.Context, params UploadCompletedParams, apiCfg *ApiConfig) (UploadedFile, error) {
	uploadedFile, err := apiCfg.DB.UpdateUploadedFile(c, database.UpdateUploadedFileParams{
		TransactionUuid: params.TransactionUuid,
		FileSize: sql.NullInt32{
			Int32: int32(params.FileSize),
			Valid: true,
		},
		FileType: sql.NullString{
			String: params.FileType,
			Valid:  true,
		},
		DownloadPresignedUrl: sql.NullString{
			String: "https://s3.download",
			Valid:  true,
		},
		Status: "File Uploaded",
	})
	if err != nil {
		fmt.Println("Error updating uploaded file:", err)
		return UploadedFile{}, fmt.Errorf("Error updating uploaded file: %w", err)
	}
	return DatabaseUploadFileToUploadFile(uploadedFile), nil
}
