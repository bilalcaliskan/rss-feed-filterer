package feed

import (
	"context"
	"errors"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/bilalcaliskan/rss-feed-filterer/internal/announce"
	"github.com/bilalcaliskan/rss-feed-filterer/internal/config"
	"github.com/bilalcaliskan/rss-feed-filterer/internal/logging"
	"github.com/bilalcaliskan/rss-feed-filterer/internal/storage/aws"
	"github.com/mmcdole/gofeed"
	"github.com/stretchr/testify/mock"
)

type mockParser struct {
	mock.Mock
}

func (m *mockParser) ParseURL(url string) (*gofeed.Feed, error) {
	args := m.Called(url)
	return args.Get(0).(*gofeed.Feed), args.Error(1)
}

func TestReleaseChecker_CheckGithubReleases(t *testing.T) {
	var sem = make(chan struct{}, 5)

	cases := []struct {
		caseName       string
		cfg            config.Repository
		ctxDuration    time.Duration
		parserResponse *gofeed.Feed
		parserErr      error
		headObjectFunc func(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error)
		getObjectFunc  func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
		putObjectFunc  func(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
	}{
		{
			"Success",
			config.Repository{
				Name:                 "project1",
				Description:          "",
				Url:                  "https://github.com/user1/project1",
				CheckIntervalMinutes: 10,
			},
			10 * time.Second,
			&gofeed.Feed{Title: "dummy"},
			nil,
			func(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error) {
				return &s3.HeadObjectOutput{}, nil
			},
			func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
				content, err := os.ReadFile("../../testdata/releases.json")
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
			"Success when ticker ticked once",
			config.Repository{
				Name:                 "project1",
				Description:          "",
				Url:                  "https://github.com/user1/project1",
				CheckIntervalMinutes: 1,
			},
			70 * time.Second,
			&gofeed.Feed{Title: "dummy"},
			nil,
			func(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error) {
				return &s3.HeadObjectOutput{}, nil
			},
			func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
				content, err := os.ReadFile("../../testdata/releases.json")
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
			"Failure caused by invalid project name",
			config.Repository{
				Name:                 "project1",
				Description:          "",
				Url:                  "https://github.com/user1",
				CheckIntervalMinutes: 10,
			},
			10 * time.Second,
			&gofeed.Feed{Title: "dummy"},
			nil,
			func(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error) {
				return &s3.HeadObjectOutput{}, nil
			},
			func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
				content, err := os.ReadFile("../../testdata/releases.json")
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
			"Failure caused by parser error",
			config.Repository{
				Name:                 "project1",
				Description:          "",
				Url:                  "https://github.com/user1/project1",
				CheckIntervalMinutes: 10,
			},
			10 * time.Second,
			nil,
			errors.New("injected error"),
			func(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error) {
				return &s3.HeadObjectOutput{}, nil
			},
			func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
				content, err := os.ReadFile("../../testdata/releases.json")
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
			"Warning caused by release file does not exists on bucket",
			config.Repository{
				Name:                 "project1",
				Description:          "",
				Url:                  "https://github.com/user1/project1",
				CheckIntervalMinutes: 10,
			},
			10 * time.Second,
			&gofeed.Feed{Title: "dummy"},
			nil,
			func(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error) {
				return nil, &types.NoSuchKey{}
			},
			func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
				content, err := os.ReadFile("../../testdata/releases.json")
				if err != nil {
					return nil, err
				}

				return &s3.GetObjectOutput{Body: io.NopCloser(strings.NewReader(string(content)))}, nil
			},
			func(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
				return &s3.PutObjectOutput{}, nil
			},
		},
		//{
		//	"Failure caused by get releases error",
		//	config.Repository{
		//		Name:                 "project1",
		//		Description:          "",
		//		Url:                  "https://github.com/user1/project1",
		//		CheckIntervalMinutes: 10,
		//	},
		//	10 * time.Second,
		//	&gofeed.Feed{Title: "dummy"},
		//	nil,
		//	func(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error) {
		//		return &s3.HeadObjectOutput{}, nil
		//	},
		//	func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
		//		return nil, errors.New("injected error")
		//	},
		//	nil,
		//},
	}

	for i := 0; i < 5; i++ {
		sem <- struct{}{}
	}

	for _, tc := range cases {
		t.Logf("starting case %s", tc.caseName)

		mockS3 := new(aws.MockS3Client)
		mockS3.HeadObjectAPI = tc.headObjectFunc
		mockS3.GetObjectAPI = tc.getObjectFunc
		mockS3.PutObjectAPI = tc.putObjectFunc

		parser := new(mockParser)
		parser.On("ParseURL", mock.AnythingOfType("string")).Return(tc.parserResponse, tc.parserErr)

		rc := NewReleaseChecker(mockS3, tc.cfg, parser, "thisisdummybucket", logging.GetLogger(), &announce.NoopAnnouncer{})

		ctx, cancel := context.WithTimeout(context.Background(), tc.ctxDuration)
		defer cancel()

		rc.CheckGithubReleases(ctx, sem)
	}
}
