package service

import (
	"bytes"
	"net/http"
	"time"
)

type UpdateFunc func() []string
type ParseFunc func(interface{}) []byte

type Callback struct {
	urls       []string
	updateFunc UpdateFunc
	parseFunc  ParseFunc
}

func NewCallback(updateFunc UpdateFunc, parseFunc ParseFunc) *Callback {
	c := &Callback{
		updateFunc: updateFunc,
		parseFunc:  parseFunc,
	}
	c.urls = updateFunc()
	ticker := time.NewTicker(2 * time.Second)
	go func() {
		for range ticker.C {
			c.urls = c.updateFunc()
		}
	}()
	return c
}

func (c *Callback) Call(data interface{}) {
	d := c.parseFunc(data)
	for _, u := range c.urls {
		go sendRequest(u, d)
	}
}

func sendRequest(url string, data []byte) {
	res, err := http.Post(url,
		"application/json;charset=utf-8", bytes.NewBuffer(data))
	if err != nil {
		return
	}

	defer res.Body.Close()
}
