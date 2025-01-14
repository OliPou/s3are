package s3uploadfile

import (
	"time"

	"github.com/OliPou/s3are/internal/database"
	"github.com/google/uuid"
)

type UploadedFile struct {
	TransactionUuid      uuid.UUID `json:"uuid"`
	Consumer             string    `json:"consumer"`
	UserName             string    `json:"userName"`
	FileName             string    `json:"fileName"`
	UploadPresignedUrl   string    `json:"uploadPresignedUrl"`
	DownloadPresignedUrl string    `json:"downloadPresignedUrl"`
	Status               string    `json:"status"`
	CreatedAt            time.Time `json:"createdAt"`
	UpdatedAT            time.Time `json:"updatedAt"`
}

func databaseUploadFileToUploadFile(dbUploadFile database.UploadedFile) UploadedFile {
	return UploadedFile{
		TransactionUuid:      dbUploadFile.TransactionUuid,
		Consumer:             dbUploadFile.Consumer,
		UserName:             dbUploadFile.UserName,
		FileName:             dbUploadFile.FileName,
		UploadPresignedUrl:   dbUploadFile.UploadPresignedUrl,
		DownloadPresignedUrl: dbUploadFile.DownloadPresignedUrl.String,
		CreatedAt:            dbUploadFile.CreatedAt,
		UpdatedAT:            dbUploadFile.UpdatedAt.Time,
		Status:               dbUploadFile.Status,
	}
}

type GetUploadedFileParams struct {
	TransactionUuid uuid.UUID `json:"transactionUuid" binding:"required"`
	UserName        string    `json:"userName" binding:"required"`
}
