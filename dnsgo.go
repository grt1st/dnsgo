package main

import (
	"github.com/miekg/dns"
	"log"
	"fmt"
	"github.com/grt1st/dnsgo/parser"
)

func main() {
	server := &dns.Server{Addr: ":53", Net: "udp"}

	server.Handler = &handler{}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
	defer server.Shutdown()
}

type handler struct{}

func (this *handler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	msg := new(dns.Msg)
	msg.SetReply(r)

	switch r.Opcode {
	case dns.OpcodeQuery:
		fmt.Println("query")
		parser.ParseQuery(msg)
	case dns.OpcodeUpdate:
		fmt.Println("update")
		/*for _, question := range r.Question {
			for _, rr := range r.Ns {
				parser.ParseUpdate(rr, &question)
			}
		}*/
	}

	w.WriteMsg(msg)
}
