package service

import (
	"net"
	"strings"

	"github.com/miekg/dns"
)

func isIP(ip string) bool {
	return net.ParseIP(ip) != nil
}

func IsUsefulIp(ip string) bool {
	seg := strings.Split(ip, "|")
	if len(seg) > 2 {
		return false
	}
	nor := strings.Split(seg[0], "&")
	for _, i := range nor {
		if isIP(i) == false {
			return false
		}
	}
	if len(seg) == 2 {
		if isIP(seg[1]) == false {
			return false
		}
	}
	return true
}

func UnFqdn(s string) string {
	if dns.IsFqdn(s) {
		return s[:len(s)-1]
	}
	return s
}
