package s3uploadfile

type ApiConfig struct {
	S3Client S3ClientInterface
	DB       DBInterface
}
