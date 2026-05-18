package system

import (
	"runtime"
	"strings"

	gnet "github.com/shirou/gopsutil/v4/net"
)

func isLoopbackInterface(name string, flags []string) bool {
	for _, flag := range flags {
		if strings.EqualFold(flag, "loopback") {
			return true
		}
	}

	lowerName := strings.ToLower(name)

	return lowerName == "lo" ||
		strings.Contains(lowerName, "loopback")
}

func isLikelyVirtualInterface(name string) bool {
	lower := strings.ToLower(name)

	return strings.Contains(lower, "docker") ||
		strings.Contains(lower, "veth") ||
		strings.Contains(lower, "virbr") ||
		strings.Contains(lower, "vbox") ||
		strings.Contains(lower, "vmnet") ||
		strings.Contains(lower, "vethernet") ||
		strings.Contains(lower, "loopback") ||
		strings.HasPrefix(lower, "br-") ||
		strings.HasPrefix(lower, "tun") ||
		strings.HasPrefix(lower, "tap")
}

func isLikelyEthernetName(name string) bool {
	lower := strings.ToLower(name)

	// Windows는 이름 규칙이 너무 제각각이라
	// virtual만 제외하고 대부분 허용
	if runtime.GOOS == "windows" {
		return true
	}

	return strings.HasPrefix(lower, "eth") ||
		strings.HasPrefix(lower, "en")
}

func hasFlag(flags []string, want string) bool {
	for _, flag := range flags {
		if strings.EqualFold(flag, want) {
			return true
		}
	}

	return false
}

func isTrackedEthernetInterface(iface gnet.InterfaceStat) bool {
	if isLoopbackInterface(iface.Name, iface.Flags) {
		return false
	}

	if isLikelyVirtualInterface(iface.Name) {
		return false
	}

	if !isLikelyEthernetName(iface.Name) {
		return false
	}

	if !hasFlag(iface.Flags, "up") {
		return false
	}

	return true
}