package handles

import (
	"log"
	"net"
	"strconv"
	"time"

	"github.com/miekg/dns"
)

type Server struct {
	Host     string
	Port     int
	RTimeout time.Duration
	WTimeout time.Duration
}

func (s *Server) Addr() string {
	return net.JoinHostPort(s.Host, strconv.Itoa(s.Port))
}

func (s *Server) Run(queryFlag bool, logfile string) {
	handler := NewHandler(queryFlag, logfile)

	tcpHandler := dns.NewServeMux()
	tcpHandler.HandleFunc(".", handler.DoTCP)

	udpHandler := dns.NewServeMux()
	udpHandler.HandleFunc(".", handler.DoUDP)

	tcpServer := &dns.Server{Addr: s.Addr(),
		Net:          "tcp",
		Handler:      tcpHandler,
		ReadTimeout:  s.RTimeout,
		WriteTimeout: s.WTimeout}

	udpServer := &dns.Server{Addr: s.Addr(),
		Net:          "udp",
		Handler:      udpHandler,
		UDPSize:      65535,
		ReadTimeout:  s.RTimeout,
		WriteTimeout: s.WTimeout}

	log.Println("[+] dns server start listening at", s.Addr())
	s.start(udpServer)
	s.start(tcpServer)

}

func (s *Server) start(dnsServer *dns.Server) {

	err := dnsServer.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
	defer dnsServer.Shutdown()

}
