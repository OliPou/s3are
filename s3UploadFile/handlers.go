package s3uploadfile

import (
	"fmt"
	"net/http"

	"github.com/OliPou/s3are/internal/common"
	"github.com/OliPou/s3are/internal/database"
	"github.com/OliPou/s3are/s3client"
	"github.com/gin-gonic/gin"
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
	var params UploadsFileParams
	if err := common.ValidateRequest(c, &params); err != nil {
		return
	}
	// Generate UUID first so we can use it in the filename
	uploadInfo, err := UploadRequest(c, params, consumer, apiCfg)
	if err != nil {
		common.RespondError(c, http.StatusInternalServerError, fmt.Sprintf("Error generating presigned URL: %v", err))
		return
	}
	common.RespondWithJSON(c, http.StatusCreated, uploadInfo)
}

func (apiCfg *ApiConfig) HandlerRequestUploadCompleted(c *gin.Context, consumer string) {
	var params UploadCompletedParams
	if err := common.ValidateRequest(c, &params); err != nil {
		return
	}
	uploadedFile, err := UploadedCompleted(c, params, apiCfg)
	if err != nil {
		common.RespondError(c, http.StatusInternalServerError, fmt.Sprintf("Error generating presigned URL: %v", err))
		return
	}

	common.RespondWithJSON(c, http.StatusOK, uploadedFile)
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

	common.RespondWithJSON(c, http.StatusOK, DatabaseUploadFileToUploadFile(uploadedFile))

}
