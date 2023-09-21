//go:build e2e

package e2e

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/bilalcaliskan/rss-feed-filterer/cmd/root"
	internalaws "github.com/bilalcaliskan/rss-feed-filterer/internal/storage/aws"
)

func TestExecuteRootCmd(t *testing.T) {
	//rootOpts := options.GetRootOptions()
	ctx := context.Background()
	root.RootCmd.SetContext(ctx)

	cases := []struct {
		caseName      string
		args          []string
		shouldPass    bool
		getObjectFunc func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
	}{
		{
			"aldsfkaslfkj",
			[]string{"--config-file=../../resources/sample_config.yaml"},
			true,
			func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
				return &s3.GetObjectOutput{}, nil
			},
		},
	}

	for _, tc := range cases {
		t.Logf("starting case %s\n", tc.caseName)

		mockS3 := new(internalaws.MockS3Client)
		mockS3.GetObjectAPI = tc.getObjectFunc

		//root.RootCmd.SetContext(context.WithValue(root.RootCmd.Context(), S3ClientKey{}, mockS3))
		//root.RootCmd.SetArgs(tc.args)
		//err := root.RootCmd.Execute()
		//if tc.shouldPass {
		//	assert.Nil(t, err)
		//} else {
		//	assert.NotNil(t, err)
		//}
	}
}
