package backends

import (
	"time"

	"github.com/miekg/dns"
)

type Record struct {
	Name string    `json:"name"`
	Ttl  time.Time `json:"ttl"`
	Mesg *dns.Msg  `json:"mesg"`
}

func (r *Record) Valid() bool {
	t := time.Now()
	if t.After(r.Ttl) {
		return false
	}
	return true
}

func GetTtl(mesg dns.Msg) time.Time {
	t := time.Now()
	var space uint32
	space = 86400
	for _, rr := range mesg.Answer {
		s := rr.Header().Ttl
		if s < space {
			space = s
		}
	}
	return t.Add(time.Second * time.Duration(int(space)))
}
