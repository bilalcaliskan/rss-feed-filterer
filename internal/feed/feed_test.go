package feed

import (
	"context"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/bilalcaliskan/rss-feed-filterer/internal/announce"
	"github.com/bilalcaliskan/rss-feed-filterer/internal/config"
	"github.com/bilalcaliskan/rss-feed-filterer/internal/storage/aws"
	"github.com/stretchr/testify/assert"
)

func TestFilter(t *testing.T) {
	cases := []struct {
		caseName       string
		shouldPass     bool
		configPath     string
		headBucketFunc func(ctx context.Context, params *s3.HeadBucketInput, optFns ...func(*s3.Options)) (*s3.HeadBucketOutput, error)
		headObjectFunc func(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error)
		getObjectFunc  func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
		putObjectFunc  func(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
	}{
		{
			"Success",
			true,
			"../../test/config.yaml",
			func(ctx context.Context, params *s3.HeadBucketInput, optFns ...func(*s3.Options)) (*s3.HeadBucketOutput, error) {
				return &s3.HeadBucketOutput{}, nil
			},
			func(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error) {
				return &s3.HeadObjectOutput{}, nil
			},
			func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
				content, err := os.ReadFile("../../test/releases.json")
				if err != nil {
					return nil, err
				}

				return &s3.GetObjectOutput{Body: io.NopCloser(strings.NewReader(string(content)))}, nil
			},
			func(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
				return &s3.PutObjectOutput{}, nil
			},
		},
		{
			"Failure caused by bucket does not exists",
			false,
			"../../test/config.yaml",
			func(ctx context.Context, params *s3.HeadBucketInput, optFns ...func(*s3.Options)) (*s3.HeadBucketOutput, error) {
				return nil, &types.NoSuchBucket{}
			},
			nil,
			nil,
			nil,
		},
		{
			"Failure caused by invalid project name",
			true,
			"../../test/invalid_config.yaml",
			func(ctx context.Context, params *s3.HeadBucketInput, optFns ...func(*s3.Options)) (*s3.HeadBucketOutput, error) {
				return &s3.HeadBucketOutput{}, nil
			},
			nil,
			nil,
			nil,
		},
	}

	for _, tc := range cases {
		t.Logf("starting case %s", tc.caseName)
		mockS3 := new(aws.MockS3Client)
		mockS3.HeadBucketAPI = tc.headBucketFunc
		mockS3.HeadObjectAPI = tc.headObjectFunc
		mockS3.GetObjectAPI = tc.getObjectFunc
		mockS3.PutObjectAPI = tc.putObjectFunc

		cfg, err := config.ReadConfig(tc.configPath)
		assert.Nil(t, err)
		assert.NotNil(t, cfg)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		// In a real test, you might want to cancel the context after some time
		// to simulate the completion of all goroutines.

		err = Filter(ctx, *cfg, mockS3, &announce.NoopAnnouncer{})
		if tc.shouldPass {
			assert.Nil(t, err)
		} else {
			assert.NotNil(t, err)
		}
	}
}
