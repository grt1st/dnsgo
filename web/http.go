package web

import (
	"github.com/gin-gonic/gin"

	"github.com/grt1st/dnsgo/web/handler"
)

func RunWeb(addr string) error {
	r := gin.Default()

	httpRouter := r.Group("http")
	httpRouter.GET("/webhook")
	httpRouter.POST("/webhook")
	httpRecordRouter := httpRouter.Group("record")
	{
		httpRecordRouter.GET("/list", handler.HTTPRecordList)
	}

	dnsRouter := r.Group("dns")
	dnsConfigRouter := dnsRouter.Group("config")
	{
		dnsConfigRouter.GET("/list", handler.DNSConfigList)
		dnsConfigRouter.POST("/set", handler.DNSConfigSet)
		dnsConfigRouter.POST("/delete", handler.DNSConfigDel)
		dnsConfigRouter.GET("/webhook")
		dnsConfigRouter.POST("/webhook")
	}
	dnsRecordRouter := dnsRouter.Group("record")
	{
		dnsRecordRouter.GET("/list", handler.DNSRecordList)
	}

	r.Any("/s/*", handler.AnyRequest)

	return r.Run(addr)
}
