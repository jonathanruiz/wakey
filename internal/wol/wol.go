package wol

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"regexp"
	"runtime"
	"time"

	probing "github.com/prometheus-community/pro-bing"
)

var (
	delims = ":-"
	reMAC  = regexp.MustCompile(`^([0-9a-fA-F]{2}[` + delims + `]){5}([0-9a-fA-F]{2})$`)
)

// MACAddress represents a 6 byte network mac address.
type MACAddress [6]byte

// MagicPacket is constituted of 6 bytes of 0xFF followed by 16-groups of the
// destination MAC address.
type MagicPacket struct {
	header  [6]byte
	payload [16]MACAddress
}

// New returns a magic packet based on a mac address string.
func New(mac string) (*MagicPacket, error) {
	var packet MagicPacket
	var macAddr MACAddress

	hwAddr, err := net.ParseMAC(mac)
	if err != nil {
		return nil, err
	}

	// We only support 6 byte MAC addresses since it is much harder to use the
	// binary.Write(...) interface when the size of the MagicPacket is dynamic.
	if !reMAC.MatchString(mac) {
		return nil, fmt.Errorf("%s is not a IEEE 802 MAC-48 address", mac)
	}

	// Copy bytes from the returned HardwareAddr -> a fixed size MACAddress.
	for idx := range macAddr {
		macAddr[idx] = hwAddr[idx]
	}

	// Setup the header which is 6 repetitions of 0xFF.
	for idx := range packet.header {
		packet.header[idx] = 0xFF
	}

	// Setup the payload which is 16 repetitions of the MAC addr.
	for idx := range packet.payload {
		packet.payload[idx] = macAddr
	}

	return &packet, nil
}

// Marshal serializes the magic packet structure into a 102 byte slice.
func (mp *MagicPacket) Marshal() ([]byte, error) {
	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.BigEndian, mp); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Wake the device
func WakeDevice(mac string) error {
	// Create a new magic packet
	packet, err := New(mac)

	// Check for errors
	if err != nil {
		return err
	}

	// Open a UDP connection to the broadcast address
	conn, err := net.Dial("udp", "255.255.255.255:9")

	// Check for errors
	if err != nil {
		return err
	}

	// Close the connection when the function returns
	defer conn.Close()

	// Marshal the magic packet
	packetBytes, err := packet.Marshal()
	if err != nil {
		return err
	}

	// Send the magic packet to the broadcast address
	_, err = conn.Write(packetBytes)
	if err != nil {
		return err
	}

	// Return nil if everything was successful
	return nil
}

// WakeGroup sends a Wake-on-LAN packet to each MAC address in the list.
func WakeGroup(macAddresses []string) error {
	for _, mac := range macAddresses {
		err := WakeDevice(mac)
		if err != nil {
			return err
		}
	}
	return nil
}

func checkOS() string {
	return runtime.GOOS
}

// Ping the device
func IsOnline(ip string) bool {

	// Get OS
	userOS := checkOS()

	pinger, err := probing.NewPinger(ip)
	if err != nil {
		panic(err)
	}

	// Check if the user is using Windows
	// If the user is using Windows, use then set the pinger to use the Windows implementation
	// If the user is not using Windows, use the default implementation
	if userOS == "windows" {
		pinger.SetPrivileged(true)
	}

	pinger.Count = 1
	pinger.Timeout = time.Second * 1 // 1 seconds

	pinger.Run() // blocks until finished

	stats := pinger.Statistics() // get send/receive/rtt stats

	if stats.PacketsRecv > 0 {
		return true
	} else {
		return false
	}
}
