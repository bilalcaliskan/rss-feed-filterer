package email

import (
	"errors"
	"testing"

	"github.com/bilalcaliskan/rss-feed-filterer/internal/announce"
	"github.com/stretchr/testify/assert"
)

type MockSender struct {
}

func (s *MockSender) Send(to, cc, bcc []string, from, projectName, version, url string) error {
	return errors.New("injected error")
}

func TestNewEmailAnnouncer(t *testing.T) {
	type args struct {
		sender Sender
		from   string
		to     []string
		cc     []string
		bcc    []string
	}
	tests := []struct {
		name string
		args args
		want *EmailAnnouncer
	}{
		{
			name: "Success",
			args: args{
				sender: nil,
				from:   "from",
				to:     []string{"to"},
				cc:     []string{"cc"},
				bcc:    []string{"bcc"},
			},
			want: &EmailAnnouncer{
				Sender: nil,
				From:   "from",
				To:     []string{"to"},
				Cc:     []string{"cc"},
				Bcc:    []string{"bcc"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewEmailAnnouncer(tt.args.sender, tt.args.from, tt.args.to, tt.args.cc, tt.args.bcc); got == nil {
				t.Errorf("NewEmailAnnouncer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEmailAnnouncer_Notify(t *testing.T) {
	announcer := &EmailAnnouncer{
		Sender: &MockSender{},
		From:   "from",
		To:     []string{"to"},
		Cc:     []string{"cc"},
		Bcc:    []string{"bcc"},
	}

	err := announcer.Notify(&announce.AnnouncerPayload{
		ProjectName: "projectName",
		Version:     "version",
		URL:         "url",
	})

	assert.NotNil(t, err)
}

func TestEmailAnnouncer_IsEnabled(t *testing.T) {
	announcer := &EmailAnnouncer{
		Sender: &MockSender{},
		From:   "from",
		To:     []string{"to"},
		Cc:     []string{"cc"},
		Bcc:    []string{"bcc"},
	}

	assert.True(t, announcer.IsEnabled())
}
