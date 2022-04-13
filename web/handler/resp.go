package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type returnMsg struct {
	c    *gin.Context
	err  error
	data interface{}
}

func JSON(c *gin.Context) *returnMsg {
	return &returnMsg{c: c}
}

func (m *returnMsg) Error(err error) *returnMsg {
	m.err = err
	return m
}

func (m *returnMsg) Data(data interface{}) *returnMsg {
	m.data = data
	return m
}

func (m *returnMsg) Return() {
	m.c.JSON(http.StatusOK, gin.H{
		"success": m.err == nil,
		"data":    m.data,
	})
}
