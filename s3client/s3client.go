// internal/s3client/s3client.go

package s3client

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3ClientInterface interface {
	GeneratePresignedURL(key string, expirationTime *int) (string, time.Duration, error)
	GeneratePresignedDownloadURL(key string, expirationTime *int) (string, time.Duration, error)
}

type S3Client struct {
	Client *s3.S3
	Bucket string
}

const DefaultPresignedURLExpiration = 24 * time.Hour

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

// Function to create Upload presigned Url on S3
func (s *S3Client) GeneratePresignedURL(key string, expirationTime *int) (string, time.Duration, error) {
	req, _ := s.Client.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
		// ContentType: aws.String(contentType),
	})
	// Use default duration if no expiration time is provided
	duration := DefaultPresignedURLExpiration
	if expirationTime != nil && *expirationTime > 0 {
		duration = time.Duration(*expirationTime) * time.Second
	}
	maxDuration := 7 * 24 * time.Hour
	if duration > maxDuration {
		duration = maxDuration
	}

	url, err := req.Presign(duration)
	if err != nil {
		return "", time.Duration(0), err
	}

	return url, duration, nil
}

// Function to create Download presigned Url on S3
func (s *S3Client) GeneratePresignedDownloadURL(key string, expirationTime *int) (string, time.Duration, error) {
	req, _ := s.Client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
	})
	// Use default duration if no expiration time is provided
	duration := DefaultPresignedURLExpiration
	if expirationTime != nil && *expirationTime > 0 {
		duration = time.Duration(*expirationTime) * time.Second
	}
	maxDuration := 7 * 24 * time.Hour
	if duration > maxDuration {
		duration = maxDuration
	}
	url, err := req.Presign(duration)
	if err != nil {
		return "", time.Duration(0), err
	}

	return url, duration, nil
}

var _ S3ClientInterface = (*S3Client)(nil) // Ensure S3Client implements S3ClientInterface
