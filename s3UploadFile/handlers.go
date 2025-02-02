package s3uploadfile

import (
	"fmt"
	"net/http"

	"github.com/OliPou/s3are/internal/common"
	"github.com/OliPou/s3are/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (apiCfg *ApiConfig) HandlerRequestUpload(c *gin.Context, consumer string) {
	var params UploadsFileParams
	if err := common.ValidateRequest(c, &params); err != nil {
		return
	}

	// Generate UUID first so we can use it in the filename
	uploadInfo, err := UploadRequest(c, params, consumer, apiCfg, uuid.New)
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
	transactionUuid, err := uuid.Parse(c.Query("transactionUuid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transactionUuid"})
		return
	}
	userName := c.Query("userName")
	if userName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userName is required"})
		return
	}
	fmt.Printf("getUploadedFileParams: %s, consumer: %s\n", transactionUuid, consumer)
	uploadedFile, err := apiCfg.DB.GetUploadedFile(c, database.GetUploadedFileParams{
		TransactionUuid: transactionUuid,
		Consumer:        consumer,
		UserName:        userName,
	})

	if err != nil {
		common.RespondWithJSON(c, http.StatusOK, fmt.Sprintf("TransactionUuid not found"))
		return
	}

	common.RespondWithJSON(c, http.StatusOK, DatabaseUploadFileToUploadFile(uploadedFile))
}
