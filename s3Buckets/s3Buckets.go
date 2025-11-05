package s3Buckets

import (
	"bytes"
	"context"
	"errors"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/aws-lambda-go/events"
)

type S3Buckets struct {
	S3Client *s3.Client
	S3Uploader *manager.Uploader
	S3Downloader *manager.Downloader
}

func (buckets S3Buckets) S3Upload(ctx context.Context, bucketName string, objectKey string, content string) (string, error) {
	var outKey string
	const Mibs int64 = 10
	input := &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key: aws.String(objectKey),
		Body: bytes.NewReader([]byte(content)),
		ChecksumAlgorithm: types.ChecksumAlgorithmSha256,
	}
	output, err := buckets.S3Uploader.Upload(ctx, input, func (u *manager.Uploader) {
		u.PartSize = Mibs * 1024 * 1024
		u.Concurrency = 10
	})
	if err != nil {
		var noBucket *types.NoSuchBucket
		if errors.As(err, &noBucket) {
			log.Printf("Bucket %s does not exist.\n", bucketName)
			err = noBucket
		}
	} else {
		err := s3.NewObjectExistsWaiter(buckets.S3Client).Wait(ctx, &s3.HeadObjectInput{
			Bucket: aws.String(bucketName),
			Key: aws.String(objectKey),
		}, time.Minute)
		if err != nil {
			log.Printf("Failed attempt to wait for object %s to exist in %s.\n", objectKey, bucketName)
		} else {
			outKey = *output.Key
		}
	}
	return outKey, err
}

func (buckets S3Buckets) S3Download(ctx context.Context, bucketName string, objectKey string) ([]byte, error){
	const Mibs int64 = 10
	buffer := manager.NewWriteAtBuffer([]byte{})
	_, err := buckets.S3Downloader.Download(ctx, buffer, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key: aws.String(objectKey),
	}, func(d *manager.Downloader) {
		d.PartSize = Mibs * 1024 * 1024
		d.Concurrency = 10
	})
	if err != nil {
		log.Printf("Couldn't download large object from %v:%v. Here's why: %v\n",
			bucketName, objectKey, err)
	}
	return buffer.Bytes(), err
}

func (buckets S3Buckets) S3TriggerLambda(ctx context.Context, s3Event events.S3Event) error {
	sdkConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Printf("failed to load default config: %s", err)
		return err
	}
	s3Client := s3.NewFromConfig(sdkConfig)

	for _, record := range s3Event.Records {
		bucket := record.S3.Bucket.Name
		key := record.S3.Object.URLDecodedKey
		headOutput, err := s3Client.HeadObject(ctx, &s3.HeadObjectInput{
			Bucket: &bucket,
			Key:    &key,
		})
		if err != nil {
			log.Printf("error getting head of object %s/%s: %s", bucket, key, err)
			return err
		}
		log.Printf("successfully retrieved %s/%s of type %s", bucket, key, *headOutput.ContentType)
	}

	return nil
}
