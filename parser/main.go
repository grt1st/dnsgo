package parser

import (
	"github.com/miekg/dns"
	"log"
	"net"
)

func ParseQuery(msg *dns.Msg) {
	for _, que := range msg.Question {

		//switch q.Qtype {}

		readR, err := getRecord(que.Name, que.Qtype)
		if err != nil {
			log.Fatalln(err)
		}

		RR := readR.(dns.RR)
		if RR.Header().Name == que.Name {
			msg.Answer = append(msg.Answer, RR)
		}
	}
}

func ParseUpdate(r dns.RR, g *dns.Question) {
	var (
		rr    dns.RR
		name  string
		rtype uint16
		ttl   uint32
		ip    net.IP
	)

	header := r.Header()
	name = header.Name
	rtype = header.Rrtype
	ttl = header.Ttl

	if _, ok := dns.IsDomainName(name); ok {
		if header.Class == dns.ClassANY && header.Rdlength == 0 {
			deleteRecord(name, rtype)
		} else { // Add record
			rheader := dns.RR_Header{
				Name:   name,
				Rrtype: rtype,
				Class:  dns.ClassINET,
				Ttl:    ttl,
			}

			if a, ok := r.(*dns.A); ok {
				rrr, err := getRecord(name, rtype)
				if err == nil {
					rr = rrr.(*dns.A)
				} else {
					rr = new(dns.A)
				}

				ip = a.A
				rr.(*dns.A).Hdr = rheader
				rr.(*dns.A).A = ip
			} else if a, ok := r.(*dns.AAAA); ok {
				rrr, err := getRecord(name, rtype)
				if err == nil {
					rr = rrr.(*dns.AAAA)
				} else {
					rr = new(dns.AAAA)
				}

				ip = a.AAAA
				rr.(*dns.AAAA).Hdr = rheader
				rr.(*dns.AAAA).AAAA = ip
			}

			storeRecord(rr)
		}
	}
}
