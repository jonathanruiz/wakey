package tests

import (
	"testing"
	"wakey/internal/common/wol"
)

func TestIsOnline(t *testing.T) {
	ipAddress := "1.1.1.1"

	// Note: This is a simple test and may need to be adjusted for your environment
	online := wol.IsOnline(ipAddress)
	if !online {
		t.Errorf("Expected %s to be online", ipAddress)
	}
}

func TestMagicPacketSize(t *testing.T) {
	packet, err := wol.New("AA:BB:CC:DD:EE:FF")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	b, err := packet.Marshal()
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	if len(b) != 102 {
		t.Errorf("expected 102 bytes, got %d", len(b))
	}
}

func TestMagicPacketHeader(t *testing.T) {
	packet, err := wol.New("AA:BB:CC:DD:EE:FF")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	b, err := packet.Marshal()
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	for i := 0; i < 6; i++ {
		if b[i] != 0xFF {
			t.Errorf("header byte[%d]: expected 0xFF, got 0x%02X", i, b[i])
		}
	}
}

func TestMagicPacketPayload(t *testing.T) {
	mac := [6]byte{0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF}

	packet, err := wol.New("AA:BB:CC:DD:EE:FF")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	b, err := packet.Marshal()
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	// Payload starts at byte 6 and repeats the MAC 16 times
	for rep := 0; rep < 16; rep++ {
		for i := 0; i < 6; i++ {
			offset := 6 + rep*6 + i
			if b[offset] != mac[i] {
				t.Errorf("payload rep %d byte %d: expected 0x%02X, got 0x%02X", rep, i, mac[i], b[offset])
			}
		}
	}
}

func TestNewInvalidMAC(t *testing.T) {
	cases := []string{
		"",
		"not-a-mac",
		"GG:HH:II:JJ:KK:LL",
		"AA:BB:CC:DD:EE",
		"AA:BB:CC:DD:EE:FF:00",
	}

	for _, mac := range cases {
		_, err := wol.New(mac)
		if err == nil {
			t.Errorf("expected error for MAC %q, got nil", mac)
		}
	}
}

func TestNewValidMAC(t *testing.T) {
	cases := []string{
		"AA:BB:CC:DD:EE:FF",
		"aa:bb:cc:dd:ee:ff",
		"00:00:00:00:00:00",
		"FF:FF:FF:FF:FF:FF",
		"AA-BB-CC-DD-EE-FF",
		"aa-bb-cc-dd-ee-ff",
	}

	for _, mac := range cases {
		_, err := wol.New(mac)
		if err != nil {
			t.Errorf("unexpected error for valid MAC %q: %v", mac, err)
		}
	}
}

func TestWakeGroupEmpty(t *testing.T) {
	err := wol.WakeGroup([]string{})
	if err != nil {
		t.Errorf("expected nil error for empty group, got %v", err)
	}
}

func TestWakeGroupInvalidMAC(t *testing.T) {
	err := wol.WakeGroup([]string{"not-a-mac"})
	if err == nil {
		t.Error("expected error for invalid MAC in group, got nil")
	}
}
