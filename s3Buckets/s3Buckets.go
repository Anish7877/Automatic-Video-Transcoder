package s3Buckets

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type UploadRequestPayload struct {
	Filename     string `json:"filename" binding:"required"`
	ContentType  string `json:"contentType" binding:"required"`
	TargetFormat string `json:"target-format" binding:"required"`
}

type S3BucketService struct {
	S3Client        *s3.Client
	S3PresignClient *s3.PresignClient
}

func NewS3BucketService(ctx context.Context) (*S3BucketService, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg)
	presignClient := s3.NewPresignClient(client)

	return &S3BucketService{
		S3Client:        client,
		S3PresignClient: presignClient,
	}, nil
}

func (s *S3BucketService) GeneratePresignedURL(ctx context.Context, bucketName string, payload UploadRequestPayload) (string, error) {

	objectKey := payload.Filename

	userMetadata := map[string]string{
		"target-format": payload.TargetFormat,
	}

	presignRequest := &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(objectKey),
		ContentType: aws.String(payload.ContentType),
		Metadata:    userMetadata,
	}

	req, err := s.S3PresignClient.PresignPutObject(
		ctx,
		presignRequest,
		s3.WithPresignExpires(15*time.Minute),
	)
	if err != nil {
		log.Printf("Couldn't get presigned URL: %v", err)
		return "", err
	}

	return req.URL, nil
}
