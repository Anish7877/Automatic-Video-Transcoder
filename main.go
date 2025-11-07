package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type SQSClientAPI interface {
	SendMessage(ctx context.Context,
		params *sqs.SendMessageInput,
		optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error)
}

type S3ClientAPI interface {
	HeadObject(ctx context.Context,
		params *s3.HeadObjectInput,
		optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error)
}

var (
	sqsClient SQSClientAPI
	s3Client  S3ClientAPI
	queueURL  string
)

func init() {
	queueURL = os.Getenv("SQS_QUEUE_URL")
	if queueURL == "" {
		panic("SQS_QUEUE_URL environment variable not set")
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(fmt.Sprintf("unable to load SDK config, %v", err))
	}

	sqsClient = sqs.NewFromConfig(cfg)
	s3Client = s3.NewFromConfig(cfg)
}

type JobMessage struct {
	InputBucket  string `json:"input_bucket"`
	InputKey     string `json:"input_key"`
	OutputBucket string `json:"output_bucket"`
	TargetFormat string `json:"target_format"`
}

func HandleRequest(ctx context.Context, s3Event events.S3Event) (string, error) {
	outputBucket := os.Getenv("OUTPUT_BUCKET_NAME")
	if outputBucket == "" {
		panic("OUTPUT_BUCKET_NAME environment variable not set")
	}

	for _, record := range s3Event.Records {
		bucket := record.S3.Bucket.Name
		key := record.S3.Object.Key

		fmt.Printf("Processing record for S3 object: %s/%s\n", bucket, key)

		headInput := &s3.HeadObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		}

		headOutput, err := s3Client.HeadObject(context.TODO(), headInput)
		if err != nil {
			fmt.Printf("Error calling HeadObject for %s/%s: %v\n", bucket, key, err)
			continue
		}
		targetFormat := "mp4"
		if val, ok := headOutput.Metadata["target-format"]; ok {
			targetFormat = val
		}

		job := JobMessage{
			InputBucket:  bucket,
			InputKey:     key,
			OutputBucket: outputBucket,
			TargetFormat: targetFormat,
		}
		messageBody, err := json.Marshal(job)
		if err != nil {
			fmt.Printf("Error marshalling job message: %v\n", err)
			continue
		}

		sendInput := &sqs.SendMessageInput{
			MessageBody: aws.String(string(messageBody)),
			QueueUrl:    aws.String(queueURL),
		}

		result, err := sqsClient.SendMessage(context.TODO(), sendInput)
		if err != nil {
			fmt.Printf("Error sending message for %s: %v\n", record.S3.Object.Key, err)
			return "Failed to send message", err
		}

		fmt.Printf("Message sent successfully for %s! Message ID: %s\n", record.S3.Object.Key, *result.MessageId)
	}

	return "All records processed successfully", nil
}

func main() {
	lambda.Start(HandleRequest)
}
