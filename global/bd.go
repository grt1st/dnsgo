package global

import "github.com/grt1st/dnsgo/backends"

var DB backends.Memory

func initBD() {
	DB.Init()
}