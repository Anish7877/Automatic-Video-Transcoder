package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"videotranscoder/aws-services/s3Buckets"
	"github.com/joho/godotenv"
	"github.com/gin-gonic/gin"
)


func main() {
	godotenv.Load()
	bucketName := os.Getenv("S3_INPUT_BUCKET_NAME")
	if bucketName == "" {
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

		uploadURL, err := s3Service.GeneratePresignedURL(c.Request.Context(), bucketName, payload)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate upload URL"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"upload_url": uploadURL,
		})
	})

	log.Println("Starting server on :8080")
	router.Run(":8080")
}
