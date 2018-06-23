package backends

import (
	"github.com/miekg/dns"
	"time"
)

type Record struct {
	Name string
	Ttl  time.Time
	Mesg *dns.Msg
}

type Backend interface {
	Init() error
	getValue() error
	DeleteRecord(record Record) error
	SaveRecord(record Record) error
	GetRecord(name string) (Record, bool)
	UpdateRecord(record Record) error
	Close() error
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
