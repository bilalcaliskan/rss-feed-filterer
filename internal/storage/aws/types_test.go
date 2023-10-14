//go:build unit

package aws

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"
)

func TestMockS3Client_GetObject(t *testing.T) {
	mockSvc := new(MockS3Client)
	mockSvc.GetObjectAPI = func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
		return &s3.GetObjectOutput{}, nil
	}

	res, err := mockSvc.GetObject(context.Background(), &s3.GetObjectInput{})
	assert.NotNil(t, res)
	assert.Nil(t, err)
}

func TestMockS3Client_PutObject(t *testing.T) {
	mockSvc := new(MockS3Client)
	mockSvc.PutObjectAPI = func(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
		return &s3.PutObjectOutput{}, nil
	}

	res, err := mockSvc.PutObject(context.Background(), &s3.PutObjectInput{})
	assert.NotNil(t, res)
	assert.Nil(t, err)
}

func TestMockS3Client_HeadObject(t *testing.T) {
	mockSvc := new(MockS3Client)
	mockSvc.HeadObjectAPI = func(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error) {
		return &s3.HeadObjectOutput{}, nil
	}

	res, err := mockSvc.HeadObject(context.Background(), &s3.HeadObjectInput{})
	assert.NotNil(t, res)
	assert.Nil(t, err)
}

func TestMockS3Client_HeadBucket(t *testing.T) {
	mockSvc := new(MockS3Client)
	mockSvc.HeadBucketAPI = func(ctx context.Context, params *s3.HeadBucketInput, optFns ...func(*s3.Options)) (*s3.HeadBucketOutput, error) {
		return &s3.HeadBucketOutput{}, nil
	}

	res, err := mockSvc.HeadBucket(context.Background(), &s3.HeadBucketInput{})
	assert.NotNil(t, res)
	assert.Nil(t, err)
}
