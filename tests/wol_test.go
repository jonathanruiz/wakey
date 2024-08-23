package tests

import (
	"testing"
	"wakey/internal/wol"
)

func TestIsOnline(t *testing.T) {
	// Added
	online := wol.IsOnline("microsoft.com")
	if !online {
		t.Errorf("Expected microsoft.com to be online")
	}
}
