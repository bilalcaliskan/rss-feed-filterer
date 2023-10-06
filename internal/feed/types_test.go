package feed

import (
	"context"
	"errors"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/bilalcaliskan/rss-feed-filterer/internal/announce/slack"
	api "github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"

	"github.com/aws/aws-sdk-go-v2/service/s3"
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

type mockSlackAPI struct {
	mock.Mock
}

func (m *mockSlackAPI) PostWebhook(url string, msg *api.WebhookMessage) error {
	args := m.Called(url, msg)
	return args.Error(0)
}

func TestReleaseChecker_CheckGithubReleases(t *testing.T) {
	var sem = make(chan struct{}, 5)

	cases := []struct {
		caseName       string
		cfg            config.Repository
		announcer      announce.Announcer
		announcerErr   error
		ctxDuration    time.Duration
		parserResponse *gofeed.Feed
		parserErr      error
		headObjectFunc func(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error)
		getObjectFunc  func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
		putObjectFunc  func(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
	}{
		{
			"Success on the first attempt",
			config.Repository{
				Name:                 "project1",
				Description:          "",
				Url:                  "https://github.com/user1/project1",
				CheckIntervalMinutes: 10,
			},
			&slack.SlackAnnouncer{},
			nil,
			10 * time.Second,
			&gofeed.Feed{
				Title:   "Release notes from project1",
				Updated: "2023-08-04T12:15:04+03:00",
				Items: []*gofeed.Item{
					{
						Title:           "v1.0.0",
						Link:            "https://github.com/user1/project1/releases/tag/v1.0.0",
						UpdatedParsed:   getTime("2023-08-04T12:21:41Z"),
						PublishedParsed: getTime("2023-08-04T12:21:41Z"),
					},
					{
						Title:           "v1.0.1",
						Link:            "https://github.com/user1/project1/releases/tag/v1.0.1",
						UpdatedParsed:   getTime("2023-08-04T12:21:41Z"),
						PublishedParsed: getTime("2023-08-04T12:21:41Z"),
					},
					{
						Title:           "v1.0.2",
						Link:            "https://github.com/user1/project1/releases/tag/v1.0.2",
						UpdatedParsed:   getTime("2023-08-04T12:21:41Z"),
						PublishedParsed: getTime("2023-08-04T12:21:41Z"),
					},
				},
			},
			nil,
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
			"Success but announcing failed/skipped",
			config.Repository{
				Name:                 "project1",
				Description:          "",
				Url:                  "https://github.com/user1/project1",
				CheckIntervalMinutes: 10,
			},
			&slack.SlackAnnouncer{},
			errors.New("injected error"),
			10 * time.Second,
			&gofeed.Feed{
				Title:   "Release notes from project1",
				Updated: "2023-08-04T12:15:04+03:00",
				Items: []*gofeed.Item{
					{
						Title:           "v1.0.0",
						Link:            "https://github.com/user1/project1/releases/tag/v1.0.0",
						UpdatedParsed:   getTime("2023-08-04T12:21:41Z"),
						PublishedParsed: getTime("2023-08-04T12:21:41Z"),
					},
					{
						Title:           "v1.0.1",
						Link:            "https://github.com/user1/project1/releases/tag/v1.0.1",
						UpdatedParsed:   getTime("2023-08-04T12:21:41Z"),
						PublishedParsed: getTime("2023-08-04T12:21:41Z"),
					},
					{
						Title:           "v1.0.2",
						Link:            "https://github.com/user1/project1/releases/tag/v1.0.2",
						UpdatedParsed:   getTime("2023-08-04T12:21:41Z"),
						PublishedParsed: getTime("2023-08-04T12:21:41Z"),
					},
				},
			},
			nil,
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
			"Success when ticker ticked once",
			config.Repository{
				Name:                 "project1",
				Description:          "",
				Url:                  "https://github.com/user1/project1",
				CheckIntervalMinutes: 1,
			},
			&announce.NoopAnnouncer{},
			nil,
			65 * time.Second,
			&gofeed.Feed{
				Title:   "Release notes from project1",
				Updated: "2023-08-04T12:21:41+03:00",
				Items: []*gofeed.Item{
					{
						Title:           "v1.0.0",
						Link:            "https://github.com/user1/project1/releases/tag/v1.0.0",
						PublishedParsed: getTime("2023-08-04T12:21:41Z"),
						UpdatedParsed:   getTime("2023-08-04T12:21:41Z"),
					},
					{
						Title:           "v1.0.1",
						Link:            "https://github.com/user1/project1/releases/tag/v1.0.1",
						PublishedParsed: getTime("2023-08-04T12:21:41Z"),
						UpdatedParsed:   getTime("2023-08-04T12:21:41Z"),
					},
				},
			},
			nil,
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
			"Failure caused by get error",
			config.Repository{
				Name:                 "project1",
				Description:          "",
				Url:                  "https://github.com/user1/project1",
				CheckIntervalMinutes: 10,
			},
			&announce.NoopAnnouncer{},
			nil,
			10 * time.Second,
			&gofeed.Feed{Title: "dummy"},
			nil,
			func(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error) {
				return &s3.HeadObjectOutput{}, nil
			},
			func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
				return nil, errors.New("injected error")
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
			&announce.NoopAnnouncer{},
			nil,
			10 * time.Second,
			nil,
			errors.New("injected error"),
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
			"Warning caused by release file does not exists on bucket",
			config.Repository{
				Name:                 "project1",
				Description:          "",
				Url:                  "https://github.com/user1/project1",
				CheckIntervalMinutes: 10,
			},
			&announce.NoopAnnouncer{},
			nil,
			10 * time.Second,
			&gofeed.Feed{Title: "dummy"},
			nil,
			func(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error) {
				return nil, &types.NoSuchKey{}
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
			"Failure caused by put releases error",
			config.Repository{
				Name:                 "project1",
				Description:          "",
				Url:                  "https://github.com/user1/project1",
				CheckIntervalMinutes: 10,
			},
			&announce.NoopAnnouncer{},
			nil,
			10 * time.Second,
			&gofeed.Feed{Title: "dummy"},
			nil,
			func(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error) {
				return nil, &types.NoSuchKey{}
			},
			func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
				content, err := os.ReadFile("../../test/releases.json")
				if err != nil {
					return nil, err
				}

				return &s3.GetObjectOutput{Body: io.NopCloser(strings.NewReader(string(content)))}, nil
			},
			func(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
				return nil, errors.New("injected error")
			},
		},
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

		_, ok := tc.announcer.(*slack.SlackAnnouncer)
		if ok {
			mockSlackAPI := new(mockSlackAPI)
			mockSlackAPI.On("PostWebhook", mock.AnythingOfType("string"), mock.AnythingOfType("*slack.WebhookMessage")).Return(tc.announcerErr)
			tc.announcer = slack.NewSlackAnnouncer("test-webhook-url", true, mockSlackAPI)
		}

		parser := new(mockParser)
		parser.On("ParseURL", mock.AnythingOfType("string")).Return(tc.parserResponse, tc.parserErr)

		rc := NewReleaseChecker(mockS3, tc.cfg, parser, "thisisdummybucket", logging.GetLogger(), tc.announcer)

		ctx, cancel := context.WithTimeout(context.Background(), tc.ctxDuration)
		defer cancel()

		projectName, err := rc.extractProjectName()
		assert.Nil(t, err)
		assert.NotEqual(t, "", projectName)

		rc.CheckGithubReleases(ctx, sem, projectName)
	}
}

func getTime(str string) *time.Time {
	t, err := time.Parse(time.RFC3339, str)
	if err != nil {
		return nil
	}

	return &t
}
