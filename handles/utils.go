package handles

import (
	"strings"

	"github.com/grt1st/dnsgo/service"
)

func ParserUrl(domain, ip string) ([]string, bool) {
	// 新建时，已经验证过了
	// 这里不再做验证了
	var ips []string
	seg := strings.Split(ip, "|")
	if len(seg) > 2 {
		return ips, false
	}
	if service.R.Get(domain)%2 == 0 && len(seg) == 2 {
		return append(ips, seg[1]), true
	}
	nor := strings.Split(seg[0], "&")
	for _, i := range nor {
		ips = append(ips, i)
	}
	return ips, true
}
