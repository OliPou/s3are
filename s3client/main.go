// internal/s3client/s3client.go

package s3client

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Client struct {
	Client *s3.S3
	Bucket string
}

func NewS3Client(region, bucket string) (*S3Client, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return nil, err
	}

	s3Client := s3.New(sess)

	return &S3Client{
		Client: s3Client,
		Bucket: bucket,
	}, nil
}

func (s *S3Client) GeneratePresignedURL(key, contentType string) (string, error) {
	req, _ := s.Client.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
		// ContentType: aws.String(contentType),
	})

	url, err := req.Presign(24 * time.Hour)
	if err != nil {
		return "", err
	}

	return url, nil
}
