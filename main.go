package main

import (
	"context"
	"log"
	"os"
	"videotranscoder/aws-services/s3Buckets"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	sdkConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("Error loading default config: %s", err)
	}

	s3Client := s3.NewFromConfig(sdkConfig)
	s3Uploader := manager.NewUploader(s3Client)
	s3Downloader := manager.NewDownloader(s3Client)

	buckets := s3Buckets.S3Buckets{
		S3Client:     s3Client,
		S3Uploader:   s3Uploader,
		S3Downloader: s3Downloader,
	}

	bucketName := "cloud-computing-project-video-uploads"
	objectKey := "hellos3.txt"
	content := "Hello, S3"

	// upload file to a particular bucket
	outKey, err := buckets.S3Upload(context.TODO(), bucketName, objectKey, content)
	if err != nil {
		log.Fatalf("failed to upload file: %v", err)
	}
	log.Printf("Successfully Uploaded file. Output Key: %s", outKey)

	// download data from a particular bucket
	data, err := buckets.S3Download(context.TODO(), bucketName, objectKey)
	if err != nil {
		log.Fatalf("Failed to download file: %v", err)
	}
	err = os.WriteFile(objectKey, data, 0644)
	if err != nil {
		log.Fatalf("Failed to Write to file: %v", err)
	}
	log.Printf("Successfully Download file: %v", objectKey)
}
