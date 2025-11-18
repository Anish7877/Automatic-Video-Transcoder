package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
	"videotranscoder/aws-services/s3Buckets"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type JobNotification struct {
	Status       string `json:"status"`
	OutputBucket string `json:"output_bucket"`
	OutputKey    string `json:"output_key"`
	InputBucket  string `json:"input_bucket"`
	InputKey     string `json:"input_key"`
	TargetFormat string `json:"target_format"`
	CompletedAt  string `json:"completed_at"`
}

type SQSService struct {
	client   *sqs.Client
	queueURL string
}

func NewSQSService(ctx context.Context, queueURL string) (*SQSService, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	return &SQSService{
		client:   sqs.NewFromConfig(cfg),
		queueURL: queueURL,
	}, nil
}

func (s *SQSService) PollForJobCompletion(ctx context.Context, expectedOutputKey string, timeout time.Duration) (*JobNotification, error) {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		// Receive message from queue
		result, err := s.client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(s.queueURL),
			MaxNumberOfMessages: 10, // Get multiple messages to search through
			WaitTimeSeconds:     20, // Long polling
			VisibilityTimeout:   30,
		})

		if err != nil {
			log.Printf("Error receiving message: %v", err)
			continue
		}

		// Check each message
		for _, message := range result.Messages {
			var notification JobNotification
			err := json.Unmarshal([]byte(*message.Body), &notification)
			if err != nil {
				log.Printf("Error parsing message: %v", err)
				continue
			}

			// Check if this is the job we're waiting for
			if notification.OutputKey == expectedOutputKey && notification.Status == "completed" {
				// Delete the message from queue
				_, err := s.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
					QueueUrl:      aws.String(s.queueURL),
					ReceiptHandle: message.ReceiptHandle,
				})
				if err != nil {
					log.Printf("Error deleting message: %v", err)
				}

				return &notification, nil
			}
		}

		// Sleep briefly before next poll
		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("timeout waiting for job completion")
}

func uploadFile(presignedURL, filePath, contentType, targetFormat string) error {
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
	req.Header.Set("x-amz-meta-target-format", targetFormat)
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
	outputBucketName := os.Getenv("S3_OUTPUT_BUCKET_NAME")
	jobNotificationQueueURL := os.Getenv("JOB_NOTIFICATION_QUEUE_URL")

	if inputBucketName == "" {
		log.Fatal("S3_INPUT_BUCKET_NAME environment variable not set")
	}
	if outputBucketName == "" {
		log.Fatal("S3_OUTPUT_BUCKET_NAME environment variable not set")
	}
	if jobNotificationQueueURL == "" {
		log.Fatal("JOB_NOTIFICATION_QUEUE_URL environment variable not set")
	}

	ctx := context.Background()

	s3Service, err := s3Buckets.NewS3BucketService(ctx)
	if err != nil {
		log.Fatalf("Failed to create S3 service: %s", err)
	}

	sqsService, err := NewSQSService(ctx, jobNotificationQueueURL)
	if err != nil {
		log.Fatalf("Failed to create SQS service: %s", err)
	}

	router := gin.Default()

	router.POST("/upload", func(c *gin.Context) {
		var payload s3Buckets.UploadRequestPayload
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Generate presigned upload URL
		uploadURL, err := s3Service.GenerateUploadPresignedURL(c.Request.Context(), inputBucketName, payload)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate upload URL"})
			return
		}

		fmt.Println("Upload URL generated:", uploadURL)

		// Upload the file
		err = uploadFile(uploadURL, payload.Filepath, payload.ContentType, payload.TargetFormat)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Upload failed: %v", err)})
			return
		}

		// Calculate the expected output key
		fileName := filepath.Base(payload.Filepath)
		originalExt := filepath.Ext(payload.Filepath)
		Ext := payload.TargetFormat
		expectedOutputKey := strings.TrimSuffix(fileName, originalExt) + "." + Ext

		fmt.Printf("Waiting for job completion for output key: %s\n", expectedOutputKey)

		// Poll SQS queue for job completion (with 5 minute timeout)
		notification, err := sqsService.PollForJobCompletion(c.Request.Context(), expectedOutputKey, 5*time.Minute)
		if err != nil {
			c.JSON(http.StatusRequestTimeout, gin.H{
				"error":   "Job processing timeout or failed",
				"message": err.Error(),
			})
			return
		}

		fmt.Printf("Job completed! Output: %s\n", notification.OutputKey)

		// Generate download presigned URL
		downloadURL, err := s3Service.GenerateDownlaodPresignedURL(c.Request.Context(), outputBucketName, notification.OutputKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to generate download URL"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":      "completed",
			"downloadUrl": downloadURL,
			"outputKey":   notification.OutputKey,
			"completedAt": notification.CompletedAt,
		})
	})

	log.Println("Starting server on :8080")
	router.Run(":8080")
}
