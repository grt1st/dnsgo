package handles

import (
	"github.com/grt1st/dnsgo/backends"
	"github.com/grt1st/dnsgo/logger"
	"github.com/miekg/dns"
	"log"
	"net"
	"strings"
)

const (
	notIPQuery = 0
	_IP4Query  = 4
	_IP6Query  = 6
)

type Handler struct {
	Hosts  Host
	Cache  backends.Backend
	Resolv Resolver
	Logger logger.ConsoleHandler
}

func NewHandler() *Handler {
	h := Handler{}
	h.Cache = &backends.Memory{
		Saver: map[string]backends.Record{},
	}
	h.Cache.Init()
	h.Hosts = *NewHost("hosts.conf")
	h.Resolv = *NewResolver()

	return &h
}

func (h *Handler) do(Net string, w dns.ResponseWriter, req *dns.Msg) {
	q := req.Question[0]

	var remoteIp net.IP
	if Net == "tcp" {
		remoteIp = w.RemoteAddr().(*net.TCPAddr).IP
	} else {
		remoteIp = w.RemoteAddr().(*net.UDPAddr).IP
	}

	IPQuery := h.isIPQuery(q)

	if IPQuery > 0 {

		ip, ok := h.Hosts.Get(q.Name)
		if ok == false {
			// 尝试匹配通配符
			slt := strings.Split(q.Name, ".")
			ip, ok = h.Hosts.Get("*." + strings.Join(slt[1:], "."))
		}

		var ips []string
		var flag bool
		if ok {
			ips, flag = ParserUrl(q.Name, ip)
		}

		if ok && flag {

			m := new(dns.Msg)
			m.SetReply(req)

			switch IPQuery {
			case _IP4Query:
				rr_header := dns.RR_Header{
					Name:   q.Name,
					Rrtype: dns.TypeA,
					Class:  dns.ClassINET,
					Ttl:    600,
				}
				for _, ip := range ips {
					a := &dns.A{rr_header, net.ParseIP(ip).To4()}
					m.Answer = append(m.Answer, a)
				}
			case _IP6Query:
				rr_header := dns.RR_Header{
					Name:   q.Name,
					Rrtype: dns.TypeAAAA,
					Class:  dns.ClassINET,
					Ttl:    600,
				}
				for _, ip := range ips {
					aaaa := &dns.AAAA{rr_header, net.ParseIP(ip).To16()}
					m.Answer = append(m.Answer, aaaa)
				}
			}

			log.Println("hosts", remoteIp, q.Name, ips, len(m.Answer))

			w.WriteMsg(m)
			return
		}
	}

	// 取出cache
	if rec, ok := h.Cache.GetRecord(q.Name); ok {
		m := rec.Mesg
		m.Id = req.Id
		log.Println("cache", remoteIp, q.Name, len(m.Answer))
		w.WriteMsg(m)
		return
	}

	mesg, err := h.Resolv.Lookup(Net, req)

	if err != nil {
		dns.HandleFailed(w, req)
		return
	}

	log.Println("lookup", remoteIp, q.Name, len(mesg.Answer))

	// 保存cache
	h.Cache.SaveRecord(backends.Record{
		Name: q.Name,
		Ttl:  backends.GetTtl(*mesg),
		Mesg: mesg,
	})
	w.WriteMsg(mesg)

	/*
		if IPQuery > 0 && len(mesg.Answer) > 0 {
			err = h.Cache.SaveRecord(Net, mesg.Answer[0].Header(A))
		}*/

}

func (h *Handler) DoTCP(w dns.ResponseWriter, req *dns.Msg) {
	h.do("tcp", w, req)
}

func (h *Handler) DoUDP(w dns.ResponseWriter, req *dns.Msg) {
	h.do("udp", w, req)
}

func (h *Handler) isIPQuery(q dns.Question) int {
	if q.Qclass != dns.ClassINET {
		return notIPQuery
	}

	switch q.Qtype {
	case dns.TypeA:
		return _IP4Query
	case dns.TypeAAAA:
		return _IP6Query
	default:
		return notIPQuery
	}
}
