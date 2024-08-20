package wol

import (
	"net"
)

// Func to check if the device is online
func IsOnline(ip string) bool {
	_, err := net.Dial("tcp", ip)
	return err == nil
}
