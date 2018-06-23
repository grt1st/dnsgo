package handles

import (
	"bufio"
	"github.com/BurntSushi/toml"
	"log"
	"net"
	"os"
	"strings"
)

type FileChecker struct {
	Filename      string
	FileTimestamp int64
}

type Settings struct {
	Host       string
	Port       int
	Nameserver string
	Hosts      string
	Resolv     string
}

var Config Settings
var R Counter

func (r *Resolver) initNameserver() {
	buf, err := os.Open("./conf/" + Config.Nameserver)
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

func (r *Resolver) initResolver() {
	buf, err := os.Open("./conf/" + Config.Resolv)
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

func init() {
	_, err := toml.DecodeFile("./conf/default.conf", &Config)
	if err != nil {
		log.Fatalln(err)
	}
	R = Counter{
		mapCount: map[string]int{},
	}
}
