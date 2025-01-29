package s3uploadfile

import (
	"database/sql"
	"time"

	"github.com/OliPou/s3are/internal/database"
	"github.com/google/uuid"
)

type UploadedFile struct {
	TransactionUuid      uuid.UUID
	Consumer             string
	UserName             string
	FileName             string
	FileSize             sql.NullInt32
	FileType             sql.NullString
	UploadPresignedUrl   string
	DownloadPresignedUrl string
	Status               string
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

func DatabaseUploadFileToUploadFile(dbUploadFile database.UploadedFile) UploadedFile {
	return UploadedFile{
		TransactionUuid:      dbUploadFile.TransactionUuid,
		Consumer:             dbUploadFile.Consumer,
		UserName:             dbUploadFile.UserName,
		FileName:             dbUploadFile.FileName,
		FileSize:             dbUploadFile.FileSize,
		FileType:             dbUploadFile.FileType,
		UploadPresignedUrl:   dbUploadFile.UploadPresignedUrl,
		DownloadPresignedUrl: dbUploadFile.DownloadPresignedUrl.String,
		Status:               dbUploadFile.Status,
		CreatedAt:            dbUploadFile.CreatedAt,
		UpdatedAt:            dbUploadFile.UpdatedAt.Time,
	}
}

type UploadsFileParams struct {
	UserName      string `json:"userName" binding:"required"`
	FileName      string `json:"fileName" binding:"required"`
	FileExtention string `json:"fileExtention" binding:"required"`
}

type GetUploadedFileParams struct {
	TransactionUuid uuid.UUID `json:"transactionUuid" binding:"required"`
	UserName        string    `json:"userName" binding:"required"`
}

type UploadCompletedParams struct {
	FileName string `json:"fileName" binding:"required"`
	FileSize int64  `json:"fileSize" binding:"required"`
	FileType string `json:"fileType" binding:"required"`
}
