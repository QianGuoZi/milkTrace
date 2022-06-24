package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"server/service"
)

func GetTrace(c *gin.Context) {
	code := c.Query("code")
	data, err := service.GetByCode(code)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "溯源失败",
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "溯源成功",
		"data":    data,
	})
}
