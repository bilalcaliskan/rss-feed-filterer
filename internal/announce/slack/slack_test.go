//go:build unit

package slack

import (
	"errors"
	"testing"

	api "github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSlackAPI struct {
	mock.Mock
}

func (m *MockSlackAPI) PostWebhook(url string, msg *api.WebhookMessage) error {
	args := m.Called(url, msg)
	return args.Error(0)
}

func TestSlackAnnouncer_Notify(t *testing.T) {
	mockSlackAPI := new(MockSlackAPI)
	announcer := NewSlackAnnouncer("test-webhook-url", true, mockSlackAPI)

	payload := SlackPayload{
		ProjectName: "Test Project",
		Version:     "1.0.0",
		URL:         "https://example.com",
		IconUrl:     "https://icon.example.com",
		Username:    "TestUser",
	}

	expectedMsg := &api.WebhookMessage{
		Attachments: []api.Attachment{},
		Username:    payload.Username,
		IconURL:     payload.IconUrl,
		Text:        "Test Project 1.0.0 is out! Check it out at https://example.com",
	}

	testCases := []struct {
		name        string
		payload     interface{}
		setupMocks  func()
		expectedErr error
	}{
		{
			name:    "successful notification",
			payload: payload,
			setupMocks: func() {
				mockSlackAPI.On("PostWebhook", "test-webhook-url", expectedMsg).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:    "notification with wrong payload type",
			payload: "wrong type",
			setupMocks: func() {
				// No mocks needed for this case
			},
			expectedErr: errors.New("invalid payload type for SlackAnnouncer"),
		},
		{
			name:    "slack API error",
			payload: payload,
			setupMocks: func() {
				mockSlackAPI.On("PostWebhook", "test-webhook-url", expectedMsg).Return(errors.New("slack API error"))
			},
			expectedErr: errors.New("slack API error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Logf("starting case %s", tc.name)
			mockSlackAPI.ExpectedCalls = nil
			mockSlackAPI.Calls = nil
			tc.setupMocks()

			err := announcer.Notify(tc.payload)
			assert.Equal(t, tc.expectedErr, err)
			mockSlackAPI.AssertExpectations(t)
		})
	}
}

func TestSlackAnnouncer_IsEnabled(t *testing.T) {
	sa := NewSlackAnnouncer("asdlfkj", true, &SlackService{})
	assert.True(t, sa.IsEnabled())
}

func TestSlackService_PostWebhook(t *testing.T) {
	ss := &SlackService{}
	assert.NotNil(t, ss.PostWebhook("asdklfj", &api.WebhookMessage{
		Attachments: []api.Attachment{},
		Username:    "alskdfjalskdfj",
		IconURL:     "aldskfasdlfkj",
		Text:        "Test Project 1.0.0 is out! Check it out at https://example.com",
	}))
}
