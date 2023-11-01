package announce

type Announcer interface {
	Notify(payload AnnouncerPayload) error
	IsEnabled() bool
}

type AnnouncerPayload struct {
	ProjectName string
	Version     string
	URL         string
}

type NoopAnnouncer struct{}

func (n *NoopAnnouncer) Notify(payload AnnouncerPayload) error {
	// Intentionally empty. This method won't do anything for the NoopAnnouncer.
	return nil
}

func (n *NoopAnnouncer) IsEnabled() bool {
	// Intentionally empty. This method won't do anything for the NoopAnnouncer.
	return false
}
