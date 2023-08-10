package email

type EmailAnnouncer struct {
	From     string
	Password string
	To       string
	Server   string
	Port     string
	Payload  EmailPayload
	Enabled  bool
}

type EmailPayload struct {
	Subject     string
	Body        string
	ProjectName string
	Version     string
	URL         string
}

func (ea *EmailAnnouncer) Notify() error {
	if !ea.Enabled {
		// Do nothing if Slack notifications are disabled
		return nil
	}

	return nil
}

func (ea *EmailAnnouncer) IsEnabled() bool {
	return ea.Enabled
}
