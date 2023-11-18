package ses

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockEmailSender struct {
	mock.Mock
}

func (m *MockEmailSender) SendEmail(ctx context.Context, params *ses.SendEmailInput, optFns ...func(*ses.Options)) (*ses.SendEmailOutput, error) {
	args := m.Called(ctx, params, optFns)
	return args.Get(0).(*ses.SendEmailOutput), args.Error(1)
}

func TestNewSESSender(t *testing.T) {
	sender := NewSESSender(nil)
	assert.NotNil(t, sender)
}

func TestSESSender_Send(t *testing.T) {
	mockSender := new(MockEmailSender)
	mockSender.On("SendEmail", mock.Anything, mock.Anything, mock.Anything).Return(&ses.SendEmailOutput{}, errors.New("injected error"))
	sender := NewSESSender(mockSender)
	assert.NotNil(t, sender)

	err := sender.Send([]string{"bilalcaliskan@protonmail.com"}, []string{}, []string{}, "bilalcaliskan@protonmail.com",
		"x-project", "1.0.0", "https://github.com/x-group/x-project")
	assert.NotNil(t, err)
}
