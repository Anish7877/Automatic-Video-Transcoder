package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"videotranscoder/aws-services/s3Buckets"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func uploadFile(presignedURL, filePath, contentType string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	stats, _ := file.Stat()

	req, err := http.NewRequest("PUT", presignedURL, file)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", contentType)
	req.ContentLength = stats.Size()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	fmt.Println("Upload successful!")
	return nil
}

func main() {
	godotenv.Load()
	inputBucketName := os.Getenv("S3_INPUT_BUCKET_NAME")
	outputBucketName := os.Getenv("S3_INPUT_BUCKET_NAME")
	if inputBucketName == "" {
		log.Fatal("S3_BUCKET_NAME environment variable not set")
	}

	ctx := context.Background()
	s3Service, err := s3Buckets.NewS3BucketService(ctx)
	if err != nil {
		log.Fatalf("Failed to create S3 service: %s", err)
	}

	router := gin.Default()

	router.POST("/upload", func(c *gin.Context) {
		var payload s3Buckets.UploadRequestPayload

		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		uploadURL, err := s3Service.GenerateUploadPresignedURL(c.Request.Context(), inputBucketName, payload)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate upload URL"})
			return
		}

		uploadFile(uploadURL, payload.Filepath, payload.ContentType);
		objectKey := filepath.Base(payload.Filepath)
		downloadUrl, err := s3Service.GenerateDownlaodPresignedURL(c.Request.Context(), outputBucketName, objectKey)
		c.JSON(http.StatusOK, gin.H{
			"downloadUrl" : downloadUrl,
		})
	})

	log.Println("Starting server on :8080")
	router.Run(":8080")
}
