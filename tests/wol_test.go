package tests

import (
	"testing"
	"wakey/internal/helper/wol"
)

func TestIsOnline(t *testing.T) {
	ipAddress := "1.1.1.1"

	// Note: This is a simple test and may need to be adjusted for your environment
	online := wol.IsOnline(ipAddress)
	if !online {
		t.Errorf("Expected %s to be online", ipAddress)
	}
}
