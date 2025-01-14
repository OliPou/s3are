package s3uploadfile

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/OliPou/s3are/internal/common"
	"github.com/OliPou/s3are/internal/database"
	"github.com/OliPou/s3are/s3client"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ApiConfig struct {
	DB       *database.Queries
	S3Client *s3client.S3Client
}

// Example usage:
func HandlerHealthz(c *gin.Context) {
	// Success response example
	status := struct {
		Status string `json:"status"`
		Ready  string `json:"ready"`
	}{
		Status: "Ok",
		Ready:  "true",
	}
	common.RespondWithJSON(c, http.StatusOK, status)
}

func (apiCfg *ApiConfig) HandlerRequestUpload(c *gin.Context, consumer string) {
	type parameters struct {
		UserName      string `json:"userName"`
		FileExtention string `json:"fileExtention"`
	}
	// Get the first value of the header or default to empty string
	decoder := json.NewDecoder(c.Request.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		common.RespondError(c, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON:", err))
		return
	}
	// Generate UUID first so we can use it in the filename
	transactionUUID := uuid.New()

	fileName := fmt.Sprintf("%s_%s_%s.%s",
		transactionUUID.String(),
		consumer,
		params.UserName,
		params.FileExtention,
	)
	presignedURL, err := apiCfg.S3Client.GeneratePresignedURL(fileName, "application/octet-stream")
	if err != nil {
		common.RespondError(c, http.StatusInternalServerError, fmt.Sprintf("Error generating presigned URL: %v", err))
		return
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
		common.RespondError(c, http.StatusInternalServerError, fmt.Sprintf("Error creating upload record: %v", err))
		return
	}

	common.RespondWithJSON(c, http.StatusCreated, databaseUploadFileToUploadFile(uploadedFile))
}

func (apiCfg *ApiConfig) HandlerRequestUploadComplete(c *gin.Context, consumer string) {
	type parameters struct {
		TransactionUuid uuid.UUID `json:"transactionUuid" binding:"required"`
		FileSize        int64     `json:"fileSize" binding:"required"`
		FileType        string    `json:"fileType" binding:"required"`
	}

	var params parameters
	if err := common.ValidateRequest(c, &params); err != nil {
		return
	}

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
		common.RespondError(c, http.StatusInternalServerError, fmt.Sprintf("Error getting upload record: %v", err))
		return
	}

	common.RespondWithJSON(c, http.StatusOK, databaseUploadFileToUploadFile(uploadedFile))
}

func (apiCfg *ApiConfig) HandlerFileStatus(c *gin.Context, consumer string) {
	var params GetUploadedFileParams
	if err := common.ValidateRequest(c, &params); err != nil {
		return
	}
	fmt.Printf("getUploadedFileParams: %s, consumer: %s\n", params.TransactionUuid, consumer)
	uploadedFile, err := apiCfg.DB.GetUploadedFile(c, database.GetUploadedFileParams{
		TransactionUuid: params.TransactionUuid,
		Consumer:        consumer,
		UserName:        params.UserName,
	})

	if err != nil {
		common.RespondWithJSON(c, http.StatusOK, fmt.Sprintf("TransactionUuid not found"))
		return
	}

	common.RespondWithJSON(c, http.StatusOK, databaseUploadFileToUploadFile(uploadedFile))

}
