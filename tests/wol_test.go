package tests

import (
	"testing"
	"wakey/internal/wol"
)

func TestIsOnline(t *testing.T) {
	// Note: This is a simple test and may need to be adjusted for your environment
	online := wol.IsOnline("1.1.1.1")
	if !online {
		t.Errorf("Expected 1.1.1.1 to be online")
	}
}
