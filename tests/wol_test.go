package tests

import (
	"testing"
	"wakey/internal/wol"
)

func TestIsOnline(t *testing.T) {
	ipAddress := "4.213.106.96"

	// Note: This is a simple test and may need to be adjusted for your environment
	online := wol.IsOnline(ipAddress)
	if !online {
		t.Errorf("Expected %s to be online", ipAddress)
	}
}
