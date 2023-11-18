package ses

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
)

// SESClient is the interface that contains SendEmail function.
// We can mock this interface for testing purposes.
type SESClient interface {
	SendEmail(ctx context.Context, params *ses.SendEmailInput, optFns ...func(*ses.Options)) (*ses.SendEmailOutput, error)
}

// SESSender is the struct that implements Sender interface. It contains the SESClient interface.
type SESSender struct {
	client SESClient
}

// NewSESSender creates a new SESSender. It requires a SESClient interface.
func NewSESSender(client SESClient) *SESSender {
	return &SESSender{
		client,
	}
}

// Send sends the email to the recipients.
func (s *SESSender) Send(to, cc, bcc []string, from, projectName, version, url string) error {
	content := fmt.Sprintf("%s %s is out! Check it out at %s", projectName, version, url)
	subject := fmt.Sprintf("New release alert for project %s!", projectName)

	input := &ses.SendEmailInput{
		Destination: &types.Destination{
			ToAddresses:  to,
			CcAddresses:  cc,
			BccAddresses: bcc,
		},
		Message: &types.Message{
			Body: &types.Body{
				Text: &types.Content{
					Data: &content,
				},
			},
			Subject: &types.Content{
				Data: &subject,
			},
		},
		Source: &from,
	}

	_, err := s.client.SendEmail(context.Background(), input)
	return err
}
