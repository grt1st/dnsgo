package handles

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/miekg/dns"

	"github.com/grt1st/dnsgo/backends"
	"github.com/grt1st/dnsgo/logger"
	"github.com/grt1st/dnsgo/service"
)

const (
	notIPQuery = iota
	_IP4Query
	_IP6Query
	_TXTQuery
	_NSQuery
)

type Handler struct {
	Hosts     *service.Host
	Cache     *backends.Memory
	Resolv    *service.Resolver
	Logger    *logger.Logger
	QueryFlag bool
}

func NewHandler(queryFlag bool, logfile string) *Handler {
	h := Handler{}
	h.Cache, _ = backends.NewMemory()
	h.Hosts = service.NewHost(service.Config.FilePath + service.Config.Hosts)
	h.Resolv = service.NewResolver(service.Config.FilePath+service.Config.Nameserver, service.Config.FilePath+service.Config.Resolv)
	h.Logger = logger.NewLogger()
	err := h.Logger.SetLogger("console", nil)
	if err != nil {
		log.Fatalln(err)
	}
	if logfile != "" {
		config := map[string]interface{}{"file": logfile}
		err = h.Logger.SetLogger("file", config)
		if err != nil {
			log.Fatalln(err)
		}
	}
	h.Logger.SetLevel(-1)
	h.QueryFlag = queryFlag
	return &h
}

// 处理 dns 请求
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

		var ip string
		var ok bool
		if IPQuery == _NSQuery {
			queryKeys := strings.Split(q.Name, ".")
			queryKeys = queryKeys[:len(queryKeys)-1]
			var value interface{}
			value, ok = service.NSMap.Load(strings.Join(queryKeys, "."))
			if ok {
				ip = value.(string)
			}
		} else {
			ip, ok = h.Hosts.Get(q.Name)
		}

		var ips []string
		var flag bool
		if ok {
			ips, flag = ParserUrl(q.Name, ip)
		}
		if IPQuery == _TXTQuery || IPQuery == _NSQuery {
			flag = true
			ips = []string{ip}
		}

		if ok && flag {

			m := new(dns.Msg)
			m.SetReply(req)

			ttl := 600
			if service.R.Has(q.Name) {
				ttl = 0
			}

			switch IPQuery {
			case _IP4Query:
				rr_header := dns.RR_Header{
					Name:   q.Name,
					Rrtype: dns.TypeA,
					Class:  dns.ClassINET,
					Ttl:    uint32(ttl),
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
					Ttl:    uint32(ttl),
				}
				for _, ip := range ips {
					aaaa := &dns.AAAA{rr_header, net.ParseIP(ip).To16()}
					m.Answer = append(m.Answer, aaaa)
				}
			case _TXTQuery:
				rr_header := dns.RR_Header{
					Name:   q.Name,
					Rrtype: dns.TypeTXT,
					Class:  dns.ClassINET,
					Ttl:    uint32(ttl),
				}
				txt := &dns.TXT{rr_header, ips}
				m.Answer = append(m.Answer, txt)
			case _NSQuery:
				rr_header := dns.RR_Header{
					Name:   q.Name,
					Rrtype: dns.TypeNS,
					Class:  dns.ClassINET,
					Ttl:    uint32(38609),
				}
				for _, ip := range ips {
					fmt.Printf("%+v, %v\n", rr_header, ip)
					ns := &dns.NS{rr_header, ip}
					m.Answer = append(m.Answer, ns)
				}
			}

			/*m.Extra = append(m.Extra, &dns.A{
				dns.RR_Header{

				},
				net.ParseIP("127.0.0.1").To4(),
			})*/

			fmt.Printf("%+v\n", m.Answer)

			h.Logger.Info("response, kind=%v, remote=%v, qname=%v, result=%v, len=%v",
				IPQuery, remoteIp, q.Name, ips, len(m.Answer))

			err := w.WriteMsg(m)
			if err != nil {
				h.Logger.Error("response err: %v", err)
			}
			return
		}
	} else {
		h.Logger.Warn("unknown query type: %+v", q.Qtype)
	}

	// 如果不支持查询其他dns服务器，退出
	if h.QueryFlag == false {
		dns.HandleFailed(w, req)
		return
	}

	// 取出cache
	if rec, ok := h.Cache.GetRecord(q.Name); ok {
		m := rec.Mesg
		m.Id = req.Id
		h.Logger.Info("cache", remoteIp, q.Name, len(m.Answer))
		w.WriteMsg(m)
		return
	}

	mesg, err := h.Resolv.Lookup(Net, req)
	if err != nil {
		h.Logger.Error(err.Error())
		dns.HandleFailed(w, req)
		return
	}
	h.Logger.Info("lookup", remoteIp, q.Name, len(mesg.Answer))

	// 保存cache
	h.Cache.SaveRecord(backends.Record{
		Name: q.Name,
		Ttl:  backends.GetTtl(*mesg),
		Mesg: mesg,
	})
	w.WriteMsg(mesg)
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
	case dns.TypeNS:
		return _NSQuery
	case dns.TypeTXT:
		return _TXTQuery
	default:
		return notIPQuery
	}
}
