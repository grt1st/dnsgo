package service

import (
	"bufio"
	"log"
	"net"
	"os"
	"strings"
)

func (h *Host) InitHosts(filename string) {
	buf, err := os.Open(filename)
	if err != nil {
		log.Fatalln(err)
	}
	defer buf.Close()
	scanner := bufio.NewScanner(buf)
	for scanner.Scan() {
		line := scanner.Text()

		// comments
		if strings.HasPrefix(line, "#") {
			continue
		}

		// pattern： sli[0] domain; sli[1] ip
		var domain, value string
		sli := strings.Split(line, " ")
		if len(sli) == 2 || (len(sli) == 3 && (sli[1] == "A" || sli[1] == "AAAA")) {
			/*
				domain1 1.2.3.4
				domain2 A 1.2.3.4
			*/
			if len(sli) == 2 {
				domain = sli[0]
				value = sli[1]
			} else {
				domain = sli[0]
				value = sli[2]
			}
			// 验证domain、ip
			if IsUsefulIp(value) == false {
				continue
			}
			// dns rebinding
			if len(strings.Split(value, "|")) == 2 {
				R.Create(domain)
			}
		} else {
			/*
				domain3 TXT "hello"
				domain4 NS ns.domain
			*/
			if len(sli) != 3 {
				continue
			}
			switch sli[1] {
			case "TXT":
				domain = sli[0]
				value = sli[2]
			case "NS":
				domain = sli[0]
				value = sli[2]
				NSMap.Store(domain, value)
				continue
			default:
				// not support yet
				continue
			}
		}

		h.Set(domain, value)
	}
}

func (r *Resolver) initNameserver(filename string) {
	buf, err := os.Open(filename)
	if err != nil {
		log.Fatalln(err)
	}
	defer buf.Close()
	scanner := bufio.NewScanner(buf)
	for scanner.Scan() {
		line := scanner.Text()

		if !strings.HasPrefix(line, "nameserver") {
			continue
		}

		sli := strings.Split(line, " ")
		if len(sli) != 2 {
			continue
		}
		r.NameServers = append(r.NameServers, net.JoinHostPort(sli[1], "53"))
	}
}

func (r *Resolver) initResolver(filename string) {
	buf, err := os.Open(filename)
	if err != nil {
		log.Fatalln(err)
	}
	defer buf.Close()
	scanner := bufio.NewScanner(buf)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if !strings.HasPrefix(line, "server") {
			continue
		}

		sli := strings.Split(line, "=")
		if len(sli) != 2 {
			continue
		}

		line = strings.TrimSpace(sli[1])

		tokens := strings.Split(line, "/")
		if len(tokens) != 3 {
			continue
		}

		domain := tokens[1]
		ip := tokens[2]

		r.Forward[domain] = ip
	}
}
