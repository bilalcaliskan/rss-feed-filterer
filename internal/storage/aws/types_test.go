package aws

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"
)

func TestMockS3Client_GetObject(t *testing.T) {
	f := func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
		return &s3.GetObjectOutput{}, nil
	}

	mock := new(MockS3Client)
	mock.GetObjectAPI = f

	res, err := mock.GetObject(context.Background(), &s3.GetObjectInput{})
	assert.NotNil(t, res)
	assert.Nil(t, err)
}

func TestMockS3Client_PutObject(t *testing.T) {
	f := func(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
		return &s3.PutObjectOutput{}, nil
	}

	mock := new(MockS3Client)
	mock.PutObjectAPI = f

	res, err := mock.PutObject(context.Background(), &s3.PutObjectInput{})
	assert.NotNil(t, res)
	assert.Nil(t, err)
}
