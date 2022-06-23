package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"server/service"
)

type UserInfo struct {
	Username string `form:"username" json:"username"`
	Password string `form:"password" json:"password"`
	Role     string `form:"role"json:"role"`
}

type Data struct {
	Token     string `json:"token"`
	ExpiredAt string `json:"expiredAt"`
}

func Login(c *gin.Context) {
	// 获取用户名、密码和角色过来
	var user UserInfo
	err := c.ShouldBind(&user)
	fmt.Println("Login传入的user信息", user)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "数据格式有误",
		})
		return
	}
	// 校验用户名和密码是否正确
	// 生成Token
	token, expiredAt, err1 := service.Login(user.Username, user.Password, user.Role)
	returnData := Data{token, expiredAt.String()}
	if err1 != nil {
		fmt.Println("handler error:", err1)
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err1,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "登陆成功",
		"data":    returnData,
	})
	return
}

func Register(c *gin.Context) {
	var user UserInfo
	err := c.ShouldBind(&user)
	fmt.Println("user", user)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "数据格式有误",
		})
		return
	}

	_, err1 := service.Register(user.Username, user.Password, user.Role)
	if err1 != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "注册失败",
		})
	} else {
		token, times, err := service.GenerateToken(user.Username, user.Password, user.Role)
		returnData := Data{token, times.String()}
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": "注册成功，自动登陆失败，请手动登陆",
			})
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "注册成功",
			"data":    returnData,
		})
	}
}

func GetUsername(c *gin.Context) {
	result, err := service.GetUsername(c)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"success": false,
		"message": result,
	})
}
