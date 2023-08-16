package announce

import "testing"

func TestNoopAnnouncer_Notify(t *testing.T) {
	n := &NoopAnnouncer{}
	err := n.Notify(nil)
	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}
}

func TestNoopAnnouncer_IsEnabled(t *testing.T) {
	n := &NoopAnnouncer{}
	if n.IsEnabled() {
		t.Fatalf("Expected IsEnabled to return false, but got true")
	}
}
