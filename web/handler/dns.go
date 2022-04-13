package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/grt1st/dnsgo/backends"
)

func DNSConfigList(c *gin.Context) {
	// parse
	var req ListParams
	err := c.ShouldBind(&req)
	if err != nil {
		JSON(c).Error(err).Return()
		return
	}
	// handle
	var results []backends.DNSConfig
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
	var resp []DNSConfigRecord
	for _, v := range results {
		resp = append(resp, DNSConfigRecord{
			ID:    v.ID,
			Key:   v.Name,
			Value: v.Value,
			Kind:  v.Kind,
		})
	}
	JSON(c).Data(gin.H{
		"data":  resp,
		"count": count,
	}).Return()
	return
}

func DNSConfigSet(c *gin.Context) {
	// parse
	var req DNSConfigRecord
	err := c.ShouldBind(&req)
	if err != nil {
		JSON(c).Error(err).Return()
		return
	}
	// handle
	if req.ID == 0 {
		err = backends.DB.Create(req.ToDBStruct()).Error
	} else {
		err = backends.DB.Where("id = ?", req.ID).Updates(req.ToDBStruct()).Error
	}
	// return
	JSON(c).Error(err).Data(gin.H{"id": req.ID}).Return()
	return
}

func DNSConfigDel(c *gin.Context) {
	// parse
	var req DNSConfigRecord
	err := c.ShouldBind(&req)
	if err != nil {
		JSON(c).Error(err).Return()
		return
	}
	// handle
	err = backends.DB.Where("id = ?", req.ID).Delete(&backends.DNSConfig{}).Error
	// return
	JSON(c).Error(err).Return()
	return
}

func DNSRecordList(c *gin.Context) {
	// parse
	var req ListParams
	err := c.ShouldBind(&req)
	if err != nil {
		JSON(c).Error(err).Return()
		return
	}
	// handle
	var results []backends.DNSRecord
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
	var resp []DNSLookupRecord
	for _, v := range results {
		resp = append(resp, DNSLookupRecord{
			ID:        v.ID,
			CreatedAt: v.CreatedAt.Format("2006-01-02 15:04:05"),
			IP:        v.ClientIP,
			Key:       v.Name,
			Value:     v.Value,
			Kind:      v.Kind,
		})
	}
	JSON(c).Data(gin.H{
		"data":  resp,
		"count": count,
	}).Return()
	return
}
