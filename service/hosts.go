package service

import (
	"os"
	"strings"
	"sync"
	"time"
)

const (
	HOST   = 600
	REBIND = 0
)

var (
	ModTime time.Time
    NSMap sync.Map
)

type Host struct {
	domain *suffixTreeNode
	sync.RWMutex
}

func NewHost(filename string) *Host {
	h := Host{
		domain: NewSuffixTreeRoot(),
	}
	h.InitHosts(filename)
	if fileInfo, _ := os.Stat(filename); fileInfo != nil {
		ModTime = fileInfo.ModTime()
	}
	ticker := time.NewTicker(2 * time.Second)
	go func() {
		for range ticker.C {
			if fileInfo, _ := os.Stat(filename); fileInfo != nil {
				nowTime := fileInfo.ModTime()
				if nowTime != ModTime {
					ModTime = nowTime
					newHost := Host{
						domain: NewSuffixTreeRoot(),
					}
					newHost.InitHosts(filename)
					h.Lock()
					h.domain = newHost.domain
					h.Unlock()
				}
			}
		}
	}()
	return &h
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
	h.set(domain, ip)
}

func (h *Host) set(domain, ip string) {
	domain = UnFqdn(domain)
	h.domain.sinsert(strings.Split(domain, "."), ip)
}

func (h *Host) Delete(qname string) {
	h.Lock()
	defer h.Unlock()
	h.delete(qname)
}

func (h *Host) delete(qname string) {
	queryKeys := strings.Split(qname, ".")
	queryKeys = queryKeys[:len(queryKeys)-1]
	h.domain.delete(queryKeys)
}

func (h *Host) Update(qname, ip string) {
	h.Lock()
	defer h.Unlock()
	h.delete(qname)
	h.set(qname, ip)
}
