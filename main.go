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
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type SQSClientAPI interface {
	SendMessage(ctx context.Context,
		params *sqs.SendMessageInput,
		optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error)
}

var (
	sqsClient SQSClientAPI
	queueURL  string
)

func init() {
	var ok bool
	queueURL = os.Getenv("SQS_QUEUE_URL")
	if !ok {
		panic("SQS_QUEUE_URL environment variable not set")
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(fmt.Sprintf("unable to load SDK config, %v", err))
	}

	sqsClient = sqs.NewFromConfig(cfg)
}

func HandleRequest(ctx context.Context, s3Event events.S3Event) (string, error) {
	for _, record := range s3Event.Records {
		fmt.Printf("Processing record for S3 object: %s/%s\n", record.S3.Bucket.Name, record.S3.Object.Key)

		messageBody, err := json.Marshal(record)
		if err != nil {
			fmt.Printf("Error marshalling S3 record: %v\n", err)
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
