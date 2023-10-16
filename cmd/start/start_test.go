//go:build e2e

package start

import (
	"context"
	"errors"
	"io"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/bilalcaliskan/rss-feed-filterer/cmd/root/options"
	"github.com/bilalcaliskan/rss-feed-filterer/internal/announce"
	"github.com/bilalcaliskan/rss-feed-filterer/internal/config"
	"github.com/bilalcaliskan/rss-feed-filterer/internal/logging"
	"github.com/bilalcaliskan/rss-feed-filterer/internal/storage/aws"
	"github.com/stretchr/testify/assert"
)

var mu sync.Mutex

func TestExecuteStartCmd(t *testing.T) {
	cases := []struct {
		caseName       string
		args           []string
		headBucketFunc func(ctx context.Context, params *s3.HeadBucketInput, optFns ...func(*s3.Options)) (*s3.HeadBucketOutput, error)
		headObjectFunc func(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error)
		getObjectFunc  func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
		putObjectFunc  func(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
		shouldPass     bool
	}{
		{
			"Success",
			[]string{},
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
			true,
		},
		{
			"Failure caused by injected error",
			[]string{},
			func(ctx context.Context, params *s3.HeadBucketInput, optFns ...func(*s3.Options)) (*s3.HeadBucketOutput, error) {
				return nil, errors.New("injected error")
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
			false,
		},
	}

	for _, tc := range cases {
		t.Logf("starting case %s\n", tc.caseName)
		StartCmd.SetContext(context.Background())

		conf, err := config.ReadConfig("../../test/config.yaml")
		assert.Nil(t, err)
		assert.NotNil(t, conf)

		mockS3 := new(aws.MockS3Client)
		mockS3.HeadBucketAPI = tc.headBucketFunc
		mockS3.HeadObjectAPI = tc.headObjectFunc
		mockS3.GetObjectAPI = tc.getObjectFunc
		mockS3.PutObjectAPI = tc.putObjectFunc

		announcer := &announce.NoopAnnouncer{}
		logger := logging.GetLogger()

		StartCmd.SetContext(context.WithValue(StartCmd.Context(), options.ConfigKey{}, conf))
		StartCmd.SetContext(context.WithValue(StartCmd.Context(), options.S3ClientKey{}, mockS3))
		StartCmd.SetContext(context.WithValue(StartCmd.Context(), options.AnnouncerKey{}, announcer))
		StartCmd.SetContext(context.WithValue(StartCmd.Context(), options.LoggerKey{}, logger))

		go func() {
			time.Sleep(10 * time.Second)
			mu.Lock()
			cancel()
			mu.Unlock()
		}()

		err = StartCmd.Execute()
		if tc.shouldPass {
			assert.Nil(t, err)
		} else {
			assert.NotNil(t, err)
		}
	}
}
