package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type BaseController struct{}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type SuccessResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (b *BaseController) Success(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code": 1,
		"msg":  message,
		"data": data,
	})
}

func (b *BaseController) Error(c *gin.Context, message string) {
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  message,
		"data": nil,
	})
}
