package service

import (
	"log"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/miekg/dns"
)

type Resolver struct {
	NameServers []string
	Forward     map[string]string
}

func NewResolver(nameserverFileName, resolverFilename  string) *Resolver {
	r := Resolver{
		NameServers: []string{},
		Forward:     map[string]string{},
	}
	r.initNameserver(nameserverFileName)
	r.initResolver(resolverFilename)
	return &r
}

func (r *Resolver) Lookup(net string, req *dns.Msg) (message *dns.Msg, err error) {
	c := &dns.Client{
		Net:          net,
		ReadTimeout:  r.Timeout(),
		WriteTimeout: r.Timeout(),
	}

	qname := req.Question[0].Name

	res := make(chan *dns.Msg, 1)
	var wg sync.WaitGroup
	L := func(nameserver string) {
		defer wg.Done()
		r, _, err := c.Exchange(req, nameserver)
		if err != nil {
			log.Println(err)
			return
		}
		if r != nil && r.Rcode == dns.RcodeServerFailure {
			return
		}
		select {
		case res <- r:
		default:
		}
	}

	ticker := time.NewTicker(time.Duration(200) * time.Millisecond)
	defer ticker.Stop()

	// Start lookup on each nameserver top-down, in every second
	NameServers := r.GetNameServers(qname)
	for _, nameserver := range NameServers {
		wg.Add(1)
		go L(nameserver)
		// but exit early, if we have an answer
		select {
		case re := <-res:
			return re, nil
		case <-ticker.C:
			continue
		}
	}
	// wait for all the namservers to finish
	wg.Wait()
	select {
	case re := <-res:
		return re, nil
	default:
		return nil, nil
	}
}

func (r *Resolver) GetNameServers(qname string) []string {
	queryKeys := strings.Split(qname, ".")
	queryKeys = queryKeys[:len(queryKeys)-1] // ignore last '.'

	var ns []string
	if v, found := r.Forward[strings.Join(queryKeys, ".")]; found {
		server := v
		nameserver := net.JoinHostPort(server, "53")
		ns = append(ns, nameserver)
		return ns
	}

	for _, nameserver := range r.NameServers {
		ns = append(ns, nameserver)
	}
	return ns
}

func (r *Resolver) Timeout() time.Duration {
	return time.Duration(5) * time.Second
}
