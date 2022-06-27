package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"server/dal"
	"server/service"
)

type RanchInput struct {
	BatchID string `form:"batchId" json:"batchId"` // 批次号
	Date    string `form:"date" json:"date"`       // 产奶日期
	Weight  int64  `form:"weight" json:"weight"`   // 总净重
}

func GetMessage(c *gin.Context) {
	// 根据 token 获得用户名
	username, err := service.GetUsername(c)
	if err != nil {
		log.Printf("[GetUserInfo] failed err=%+v", err)
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	// 根据用户名获取用户信息
	user, err := dal.GetUserInfoByName(username)
	if err != nil {
		log.Printf("[GetUserInfo] failed err=%+v", err)
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	if user.Role == "0" {
		resultList, err := service.GetInfoRanch(user.Id)
		if err != nil {
			log.Printf("[GetInfoRanch] failed err=%+v", err)
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "成功获取产品信息",
			"data":    resultList,
		})
	}

}

func SetMessage(c *gin.Context) {
	// 根据 token 获得用户名
	username, err := service.GetUsername(c)
	if err != nil {
		log.Printf("[GetUserInfo] failed err=%+v", err)
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	// 根据用户名获取用户信息
	user, err := dal.GetUserInfoByName(username)
	if err != nil {
		log.Printf("[GetUserInfo] failed err=%+v", err)
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	code := c.Query("code")

	if user.Role == "0" {
		var ranchInput = RanchInput{}
		c.ShouldBind(&ranchInput)
		fmt.Println("ranchInput", ranchInput)
		err := service.AddInfoRanch(user.Id, ranchInput.BatchID, ranchInput.Date, ranchInput.Weight)
		if err != nil {
			log.Printf("[AddInfoRanch] failed err=%+v", err)
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "成功添加产品信息",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "成功添加产品信息",
			"code":    code,
		})
	}
}
