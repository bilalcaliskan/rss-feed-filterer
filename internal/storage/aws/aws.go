package aws

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	internaltypes "github.com/bilalcaliskan/rss-feed-filterer/internal/types"
)

func CreateConfig(accessKey, secretKey, region string) (aws.Config, error) {
	appCreds := aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""))
	return config.LoadDefaultConfig(context.Background(),
		config.WithRegion(region),
		config.WithCredentialsProvider(appCreds),
	)
}

func CreateClient(accessKey, secretKey, region string) (*s3.Client, error) {
	cfg, err := CreateConfig(accessKey, secretKey, region)
	if err != nil {
		return nil, err
	}

	return s3.NewFromConfig(cfg), nil
}

func GetReleases(client S3ClientAPI, bucketName, key string) (releases []internaltypes.Release, err error) {
	mu := &sync.Mutex{}

	// fetch all the objects in target bucket
	getResult, err := client.GetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})

	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	mu.Lock()
	if _, err := buf.ReadFrom(getResult.Body); err != nil {
		return nil, err
	}
	mu.Unlock()

	if err := json.Unmarshal(buf.Bytes(), &releases); err != nil {
		return nil, err
	}

	return releases, nil
}

func IsObjectExists(client S3ClientAPI, bucketName, key string) bool {
	_, err := client.HeadObject(context.Background(), &s3.HeadObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})

	if err != nil {
		var nfErr *types.NoSuchKey
		if errors.As(err, &nfErr) {
			// The object does not exist
			return false
		}

		return false
	}

	return true
}

func IsBucketExists(client S3ClientAPI, bucketName string) bool {
	_, err := client.HeadBucket(context.TODO(), &s3.HeadBucketInput{
		Bucket: &bucketName,
	})

	if err != nil {
		var awsErr *types.NoSuchBucket
		if ok := errors.As(err, &awsErr); ok {
			return false
		}

		// Some other error occurred (e.g., Forbidden)
		return false
	}

	return true
}

func PutReleases(client S3ClientAPI, bucketName, key string, releases []internaltypes.Release) error {
	data, err := json.MarshalIndent(&releases, "", "    ")
	if err != nil {
		return err
	}

	input := &s3.PutObjectInput{
		Bucket:        aws.String(bucketName),
		Key:           aws.String(key),
		Body:          strings.NewReader(string(data)),
		ContentLength: aws.Int64(int64(len(data))),
		ContentType:   aws.String("application/json"),
	}

	if _, err := client.PutObject(context.Background(), input); err != nil {
		return err
	}

	return nil
}
