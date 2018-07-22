package handles

import (
	"bufio"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"
	"github.com/miekg/dns"
)

const (
	HOST = 600
	REBIND = 0
)

var ModTime time.Time

type Host struct {
	domain *suffixTreeNode
	sync.RWMutex
}

func NewHost(filename string) *Host {
	h := Host{
		domain: newSuffixTreeRoot(),
	}
	h.InitHosts(filename)
	fileInfo, _ := os.Stat("./conf/" + filename)
	ModTime = fileInfo.ModTime()
	go h.CheckUpdate(filename)
	return &h
}

func (h *Host) InitHosts(filename string) {
	buf, err := os.Open("./conf/" + filename)
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

		// pattern： sli[0] domain; sli[1 ip
		sli := strings.Split(line, " ")
		if len(sli) != 2 {
			continue
		}

		// 验证domain、ip
		if IsUsefulIp(sli[1]) == false {
			continue
		}

		if len(strings.Split(sli[1], "|")) == 2 {
			R.Create(sli[0])
		}

		h.Set(sli[0], sli[1])
	}
}

func (h *Host) Get(qname string) (string, bool) {
	h.Lock()
	defer h.Unlock()
	queryKeys := strings.Split(qname, ".")
	queryKeys = queryKeys[:len(queryKeys)-1]
	key, ok := h.domain.searchWidcard(queryKeys)
	return key, ok
}

func (h *Host) Set(domain, ip string) {
	h.Lock()
	defer h.Unlock()
	domain = UnFqdn(domain)
	h.domain.sinsert(strings.Split(domain, "."), ip)
}

func (h *Host) Delete(qname string) {
	h.Lock()
	h.Unlock()
	queryKeys := strings.Split(qname, ".")
	queryKeys = queryKeys[:len(queryKeys)-1]
	h.domain.delete(queryKeys)
}

func (h *Host) Refresh(filename string) {
	fileInfo, _ := os.Stat("./conf/" + filename)
	NowTime := fileInfo.ModTime()

	if NowTime.After(ModTime){
		newH := NewHost(filename)
		h.domain = newH.domain
	}
}

func ParserUrl(domain, ip string) ([]string, bool) {
	// 新建时，已经验证过了
	// 这里不再做验证了
	var ips []string
	seg := strings.Split(ip, "|")
	if len(seg) > 2 {
		return ips, false
	}
	if R.Get(domain)%2 == 0 && len(seg) == 2 {
		return append(ips, seg[1]), true
	}
	nor := strings.Split(seg[0], "&")
	for _, i := range nor {
		ips = append(ips, i)
	}
	return ips, true
}

func isIP(ip string) bool {
	return net.ParseIP(ip) != nil
}

func IsUsefulIp(ip string) bool {
	seg := strings.Split(ip, "|")
	if len(seg) > 2 {
		return false
	}
	nor := strings.Split(seg[0], "&")
	for _, i := range nor {
		if isIP(i) == false {
			return false
		}
	}
	if len(seg) == 2 {
		if isIP(seg[1]) == false {
			return false
		}
	}
	return true
}

type Counter struct {
	mapCount map[string]int
	sync.RWMutex
}

func (c *Counter) Create(domain string) {
	c.Lock()
	defer c.Unlock()
	if _, ok := c.mapCount[domain]; ok == false {
		c.mapCount[domain] = 1
	}
}

func (c *Counter) Get(domain string) int {
	c.Lock()
	defer c.Unlock()
	count, ok := c.mapCount[domain]
	if ok {
		c.mapCount[domain]++
		return count
	}
	return 1
}

func (c *Counter) Has(domain string) bool {
	c.Lock()
	defer c.Unlock()
	_, ok := c.mapCount[domain]
	return ok
}

func UnFqdn(s string) string {
	if dns.IsFqdn(s) {
		return s[:len(s)-1]
	}
	return s
}

func (node *suffixTreeNode) searchWidcard(keys []string) (string, bool) {

	if len(keys) == 0 {
		return "", false
	}

	key := keys[len(keys)-1]
	n, ok := node.children[key]
	if ok == false {
		n, ok = node.children["*"]
	}
	if ok {
		if nextValue, found := n.searchWidcard(keys[:len(keys)-1]); found {
			return nextValue, found
		}
		return n.value, (n.value != "")
	}

	return "", false
}

func (h *Host) CheckUpdate(filename string) {
	for {
		time.Sleep(2 * time.Second)
		fileInfo, err := os.Stat("./conf/" + filename)
		if err != nil {
			log.Println(err, ModTime)
			continue
		}
		//fmt.Println(fileInfo)
		NowTime := fileInfo.ModTime()
		if NowTime != ModTime {
			ModTime = NowTime
			nh := Host{
				domain: newSuffixTreeRoot(),
			}
			nh.InitHosts(filename)
			h.Lock()
			h.domain = nh.domain
			h.Unlock()
		}
	}
}