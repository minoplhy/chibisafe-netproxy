package handler

import (
	"net"
	"strings"
)

func IsInternalIP(ip string) bool {

	parts := strings.Split(ip, ":")
	ip = parts[0]

	// Parse the IP address
	ipAddress := net.ParseIP(ip)
	if ipAddress == nil {
		return false
	}

	internalRanges := []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
	}

	for _, internalRange := range internalRanges {
		_, ipNet, err := net.ParseCIDR(internalRange)
		if err != nil {
			continue
		}

		if ipNet.Contains(ipAddress) {
			return true
		}
	}

	return false
}
