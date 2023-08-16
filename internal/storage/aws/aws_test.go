package aws

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	types2 "github.com/bilalcaliskan/rss-feed-filterer/internal/types"
	"github.com/stretchr/testify/assert"
)

// errorReader is a type that implements both the io.Reader and io.Closer interfaces.
// However, its Read method is designed to always return an error, simulating a faulty reader.
type errorReader struct{}

// Read is the implementation of the io.Reader interface's Read method for the errorReader type.
// Instead of performing any actual reading operation, it immediately returns an error.
func (er errorReader) Read(p []byte) (n int, err error) {
	return 0, fmt.Errorf("forced error") // No data read and a forced error is returned.
}

// Close is the implementation of the io.Closer interface's Close method for the errorReader type.
// It currently does nothing and returns no error, but could be modified to return an error if desired.
func (er errorReader) Close() error {
	return nil // Returning nil indicating that the close operation was successful.
	// Replace with an error if you want to simulate a Close operation failure.
}

func TestCreateClient(t *testing.T) {
	client, err := CreateClient("alksdfjalsdkf", "alskdfjalksdfj", "us-east-1")
	assert.NotNil(t, client)
	assert.Nil(t, err)
}

func TestIsObjectExists(t *testing.T) {
	cases := []struct {
		caseName string
		expected bool
		headFunc func(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error)
	}{
		{
			"Object exists",
			true,
			func(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error) {
				return &s3.HeadObjectOutput{}, nil
			},
		},
		{
			"Object does not exists",
			false,
			func(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error) {
				return nil, &types.NoSuchKey{}
			},
		},
		{
			"Other error returned",
			false,
			func(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error) {
				return nil, errors.New("injected error")
			},
		},
	}

	for _, tc := range cases {
		t.Logf("starting case %s", tc.caseName)

		mockS3 := new(MockS3Client)
		mockS3.HeadObjectAPI = tc.headFunc

		assert.Equal(t, tc.expected, IsObjectExists(mockS3, "thisisdemobucket", "thisisdummykey"))
	}
}

func TestIsBucketExists(t *testing.T) {
	cases := []struct {
		caseName string
		expected bool
		headFunc func(ctx context.Context, params *s3.HeadBucketInput, optFns ...func(*s3.Options)) (*s3.HeadBucketOutput, error)
	}{
		{
			"Bucket exists",
			true,
			func(ctx context.Context, params *s3.HeadBucketInput, optFns ...func(*s3.Options)) (*s3.HeadBucketOutput, error) {
				return &s3.HeadBucketOutput{}, nil
			},
		},
		{
			"Bucket does not exists",
			false,
			func(ctx context.Context, params *s3.HeadBucketInput, optFns ...func(*s3.Options)) (*s3.HeadBucketOutput, error) {
				return nil, &types.NoSuchBucket{}
			},
		},
		{
			"Other error returned",
			false,
			func(ctx context.Context, params *s3.HeadBucketInput, optFns ...func(*s3.Options)) (*s3.HeadBucketOutput, error) {
				return nil, errors.New("injected error")
			},
		},
	}

	for _, tc := range cases {
		t.Logf("starting case %s", tc.caseName)

		mockS3 := new(MockS3Client)
		mockS3.HeadBucketAPI = tc.headFunc

		assert.Equal(t, tc.expected, IsBucketExists(mockS3, "thisisdemobucket"))
	}
}

func getTime(str string) *time.Time {
	t, err := time.Parse(time.RFC3339, str)
	if err != nil {
		return nil
	}

	return &t
}

func TestGetReleases(t *testing.T) {
	cases := []struct {
		caseName         string
		expectedReleases []types2.Release
		shouldPass       bool
		getFunc          func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
	}{
		{
			"Success",
			[]types2.Release{
				{
					ProjectName: "user1/project1",
					Version:     "v1.0.0",
					PublishedAt: getTime("2023-08-06T16:57:13Z"),
					UpdatedAt:   getTime("2023-08-06T16:57:13Z"),
					Url:         "https://github.com/user1/project1/releases/tag/v1.0.0",
				},
				{
					ProjectName: "user1/project1",
					Version:     "v1.0.1",
					PublishedAt: getTime("2023-07-16T12:55:36Z"),
					UpdatedAt:   getTime("2023-07-16T12:55:36Z"),
					Url:         "https://github.com/user1/project1/releases/tag/v1.0.1",
				},
			},
			true,
			func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
				content, err := os.ReadFile("../../../testdata/releases.json")
				if err != nil {
					return nil, err
				}

				body := strings.NewReader(string(content))

				return &s3.GetObjectOutput{
					Body: io.NopCloser(body),
				}, nil
			},
		},
		{
			"Failure caused by invalid json",
			nil,
			false,
			func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
				content, err := os.ReadFile("../../../testdata/releases_invalid.json")
				if err != nil {
					return nil, err
				}

				body := strings.NewReader(string(content))

				return &s3.GetObjectOutput{
					Body: io.NopCloser(body),
				}, nil
			},
		},
		{
			"Failure caused by invalid io.Reader source",
			nil,
			false,
			func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
				return &s3.GetObjectOutput{
					Body: errorReader{},
				}, nil
			},
		},
		{
			"Failure",
			nil,
			false,
			func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
				return nil, errors.New("injected error")
			},
		},
	}

	for _, tc := range cases {
		t.Logf("starting case %s", tc.caseName)

		mockS3 := new(MockS3Client)
		mockS3.GetObjectAPI = tc.getFunc

		res, err := GetReleases(mockS3, "thisisdemobucket", "project1")
		assert.Equal(t, tc.expectedReleases, res)

		if tc.shouldPass {
			assert.Nil(t, err)
		} else {
			assert.NotNil(t, err)
		}
	}
}

func TestPutReleases(t *testing.T) {
	cases := []struct {
		caseName string
		expected error
		putFunc  func(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
	}{
		{
			"Success",
			nil,
			func(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
				return &s3.PutObjectOutput{}, nil
			},
		},
		{
			"Failure caused by invalid json",
			errors.New("injected error"),
			func(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
				return nil, errors.New("injected error")
			},
		},
	}

	for _, tc := range cases {
		t.Logf("starting case %s", tc.caseName)

		mockS3 := new(MockS3Client)
		mockS3.PutObjectAPI = tc.putFunc

		assert.Equal(t, tc.expected, PutReleases(mockS3, "thisisdemobucket", "thisisdemokey", []types2.Release{}))
	}
}
