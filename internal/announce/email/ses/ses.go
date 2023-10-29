package ses

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
)

type SESSender struct {
	client *ses.Client
}

func NewSESSender(client *ses.Client) *SESSender {
	return &SESSender{
		client,
	}
}

func (s *SESSender) Send(to, cc, bcc []string, from, subject, content string) error {
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
