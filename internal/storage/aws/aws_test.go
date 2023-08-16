package aws

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateClient(t *testing.T) {
	client, err := CreateClient("alksdfjalsdkf", "alskdfjalksdfj", "us-east-1")
	assert.NotNil(t, client)
	assert.Nil(t, err)
}
