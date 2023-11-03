//go:build unit

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
	cases := []struct {
		caseName             string
		cfg                  config.Repository
		announcers           []announce.Announcer
		announcerErr         error
		ctxDuration          time.Duration
		parserResponse       *gofeed.Feed
		parserErr            error
		checkIntervalSeconds int
		headObjectFunc       func(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error)
		getObjectFunc        func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
		putObjectFunc        func(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
	}{
		{
			"Success on the first attempt",
			config.Repository{
				Name:                 "project1",
				Description:          "",
				Url:                  "https://github.com/user1/project1",
				CheckIntervalMinutes: 1,
			},
			[]announce.Announcer{&slack.SlackAnnouncer{}},
			nil,
			10 * time.Second,
			&gofeed.Feed{
				Title:   "Release notes from project1",
				Updated: "2023-08-04T12:15:04+03:00",
				Items: []*gofeed.Item{
					{
						Title:           "v1.0.0",
						Link:            "https://github.com/user1/project1/releases/tag/v1.0.0",
						UpdatedParsed:   getTimeFromString("2023-08-04T12:21:41Z"),
						PublishedParsed: getTimeFromString("2023-08-04T12:21:41Z"),
					},
					{
						Title:           "v1.0.1",
						Link:            "https://github.com/user1/project1/releases/tag/v1.0.1",
						UpdatedParsed:   getTimeFromString("2023-08-04T12:21:41Z"),
						PublishedParsed: getTimeFromString("2023-08-04T12:21:41Z"),
					},
					{
						Title:           "v1.0.2",
						Link:            "https://github.com/user1/project1/releases/tag/v1.0.2",
						UpdatedParsed:   getTimeFromString("2023-08-04T12:21:41Z"),
						PublishedParsed: getTimeFromString("2023-08-04T12:21:41Z"),
					},
				},
			},
			nil,
			10,
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
				CheckIntervalMinutes: 1,
			},
			[]announce.Announcer{&slack.SlackAnnouncer{}},
			errors.New("injected error"),
			10 * time.Second,
			&gofeed.Feed{
				Title:   "Release notes from project1",
				Updated: "2023-08-04T12:15:04+03:00",
				Items: []*gofeed.Item{
					{
						Title:           "v1.0.0",
						Link:            "https://github.com/user1/project1/releases/tag/v1.0.0",
						UpdatedParsed:   getTimeFromString("2023-08-04T12:21:41Z"),
						PublishedParsed: getTimeFromString("2023-08-04T12:21:41Z"),
					},
					{
						Title:           "v1.0.1",
						Link:            "https://github.com/user1/project1/releases/tag/v1.0.1",
						UpdatedParsed:   getTimeFromString("2023-08-04T12:21:41Z"),
						PublishedParsed: getTimeFromString("2023-08-04T12:21:41Z"),
					},
					{
						Title:           "v1.0.2",
						Link:            "https://github.com/user1/project1/releases/tag/v1.0.2",
						UpdatedParsed:   getTimeFromString("2023-08-04T12:21:41Z"),
						PublishedParsed: getTimeFromString("2023-08-04T12:21:41Z"),
					},
				},
			},
			nil,
			10,
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
			[]announce.Announcer{&announce.NoopAnnouncer{}},
			nil,
			65 * time.Second,
			&gofeed.Feed{
				Title:   "Release notes from project1",
				Updated: "2023-08-04T12:21:41+03:00",
				Items: []*gofeed.Item{
					{
						Title:           "v1.0.0",
						Link:            "https://github.com/user1/project1/releases/tag/v1.0.0",
						PublishedParsed: getTimeFromString("2023-08-04T12:21:41Z"),
						UpdatedParsed:   getTimeFromString("2023-08-04T12:21:41Z"),
					},
					{
						Title:           "v1.0.1",
						Link:            "https://github.com/user1/project1/releases/tag/v1.0.1",
						PublishedParsed: getTimeFromString("2023-08-04T12:21:41Z"),
						UpdatedParsed:   getTimeFromString("2023-08-04T12:21:41Z"),
					},
				},
			},
			nil,
			10,
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
				CheckIntervalMinutes: 1,
			},
			[]announce.Announcer{&announce.NoopAnnouncer{}},
			nil,
			10 * time.Second,
			&gofeed.Feed{Title: "dummy"},
			nil,
			10,
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
				CheckIntervalMinutes: 1,
			},
			[]announce.Announcer{&announce.NoopAnnouncer{}},
			nil,
			10 * time.Second,
			nil,
			errors.New("injected error"),
			10,
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
				CheckIntervalMinutes: 1,
			},
			[]announce.Announcer{&announce.NoopAnnouncer{}},
			nil,
			10 * time.Second,
			&gofeed.Feed{Title: "dummy"},
			nil,
			10,
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
				CheckIntervalMinutes: 1,
			},
			[]announce.Announcer{&announce.NoopAnnouncer{}},
			nil,
			10 * time.Second,
			&gofeed.Feed{Title: "dummy"},
			nil,
			10,
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

	// fill the semaphore channel with empty objects
	//for i := 0; i < 5; i++ {
	//	sem <- struct{}{}
	//}

	for _, tc := range cases {
		t.Logf("starting case %s", tc.caseName)

		// create a mock S3 client
		mockS3 := new(aws.MockS3Client)
		// override the headObjectFunc with mock headObjectFunc
		mockS3.HeadObjectAPI = tc.headObjectFunc
		// override the getObjectFunc with mock getObjectFunc
		mockS3.GetObjectAPI = tc.getObjectFunc
		// override the putObjectFunc with mock putObjectFunc
		mockS3.PutObjectAPI = tc.putObjectFunc

		var anns []announce.Announcer

		// create a mock announcer
		for _, a := range tc.announcers {
			_, ok := a.(*slack.SlackAnnouncer)
			if ok {
				// create a mock Slack API
				mockSlackAPI := new(mockSlackAPI)
				// override the PostWebhook with mock PostWebhook
				mockSlackAPI.On("PostWebhook", mock.AnythingOfType("string"), mock.AnythingOfType("*slack.WebhookMessage")).Return(tc.announcerErr)
				// override the announcer with mock announcer
				anns = append(anns, slack.NewSlackAnnouncer("test-webhook-url", "foo", "aldskfadsfljk", mockSlackAPI))
				continue
			}

			anns = append(anns, a)
		}

		// create a mock parser
		parser := new(mockParser)
		// override the ParseURL with mock ParseURL
		parser.On("ParseURL", mock.AnythingOfType("string")).Return(tc.parserResponse, tc.parserErr)

		// create a release checker
		rc := NewReleaseChecker(mockS3, tc.cfg, parser, "thisisdummybucket", logging.GetLogger(), anns)

		// create a context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), tc.ctxDuration)
		defer cancel()

		// extract project name from repository URL
		projectName, err := rc.extractProjectName()
		assert.Nil(t, err)
		assert.NotEqual(t, "", projectName)

		// check the feed
		rc.CheckGithubReleases(ctx, projectName)
	}
}

// getTimeFromString parses a string to time.Time and returns pointer to it.
func getTimeFromString(str string) *time.Time {
	t, err := time.Parse(time.RFC3339, str)
	if err != nil {
		return nil
	}

	return &t
}
