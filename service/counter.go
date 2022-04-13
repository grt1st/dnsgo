package service

import "sync"

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
