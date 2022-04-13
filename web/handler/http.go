package handler

import (
	"encoding/base64"
	"encoding/json"

	"github.com/gin-gonic/gin"

	"github.com/grt1st/dnsgo/backends"
)

func AnyRequest(c *gin.Context) {
	var reqBody string
	if c.Request.Body != nil {
		var data []byte
		c.Request.Body.Read(data)
		reqBody = base64.StdEncoding.EncodeToString(data)
	}
	var reqHeader []byte
	reqHeader, _ = json.Marshal(c.Request.Header)
	backends.DB.Create(&backends.HTTPRecord{
		ClientIP:   c.ClientIP(),
		ReqHeaders: string(reqHeader),
		ReqBody:    reqBody,
		URL:        c.Request.URL.String(),
	})
}

func HTTPRecordList(c *gin.Context) {
	// parse
	var req ListParams
	err := c.ShouldBind(&req)
	if err != nil {
		JSON(c).Error(err).Return()
		return
	}
	// handle
	var results []backends.HTTPRecord
	var count int64
	err = backends.DB.Find(&results).Count(&count).Error
	if err != nil {
		JSON(c).Error(err).Return()
		return
	}
	err = backends.DB.Find(&results).Offset(req.Offset).Limit(req.Limit).Error
	if err != nil {
		JSON(c).Error(err).Return()
		return
	}
	// return
	var resp []HTTPRecord
	for _, v := range results {
		resp = append(resp, HTTPRecord{
			ID:         v.ID,
			CreatedAt:  v.CreatedAt.Format("2006-01-02 15:04:05"),
			IP:         v.ClientIP,
			ReqHeaders: nil,
			ReqBody:    v.ReqBody,
			URL:        v.URL,
		})
	}
	JSON(c).Data(gin.H{
		"data":  resp,
		"count": count,
	}).Return()
	return
}