package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"server/dal"
	"strings"
	"time"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	rand.Seed(time.Now().UnixNano())

	dal.InitDB()

	dal.InitTls()

	r := gin.Default()
	r.Use(Cors())
	//r.Use(cors.Default())

	initRouter(r)

	err := r.Run(":8080") // http端口
	if err != nil {
		panic(err)
	}

}

//func Cors() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		method := c.Request.Method
//		origin := c.Request.Header.Get("Origin")
//		if origin != "" {
//			c.Header("Access-Control-Allow-Origin", "*")
//			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
//			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
//			c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization")
//			//c.Header("Access-Control-Allow-Credentials", "true")
//			c.Set("content-type", "application/json")
//		}
//		//放行所有OPTIONS方法
//		if method == "OPTIONS" {
//			c.AbortWithStatus(http.StatusNoContent)
//		}
//		c.Next()
//	}
//}

func Cors() gin.HandlerFunc {
	return func(context *gin.Context) {
		method := context.Request.Method
		origin := context.Request.Header.Get("Origin")
		var headerKeys []string
		for k, _ := range context.Request.Header {
			headerKeys = append(headerKeys, k)
		}
		headerStr := strings.Join(headerKeys, ",")
		if headerStr != "" {
			headerStr = fmt.Sprintf("access-control-allow-origin, access-control-allow-headers, %s", headerStr)
		} else {
			headerStr = "access-control-allow-origin, access-control-allow-headers"
		}

		if origin != "" {
			context.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			context.Header("Access-Control-Allow-Origin", "*") // 设置允许访问所有域
			context.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
			context.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")
			context.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar")
			context.Header("Access-Control-Max-Age", "172800")
			context.Header("Access-Control-Allow-Credentials", "false")
			context.Set("content-type", "application/json") //// 设置返回格式是json
		}

		if method == "OPTIONS" {
			context.JSON(http.StatusOK, "Options Request!")
		}
		//处理请求
		context.Next()
	}
}

//
//func Cors() gin.HandlerFunc {
//	return func(context *gin.Context) {
//		method := context.Request.Method
//		// 必须，接受指定域的请求，可以使用*不加以限制，但不安全
//		//context.Header("Access-Control-Allow-Origin", "*")
//		context.Header("Access-Control-Allow-Origin", context.GetHeader("Origin"))
//		//context.Writer.Header().Set("Access-Control-Allow-Origin", context.GetHeader("Origin"))
//		fmt.Println(context.GetHeader("Origin"))
//		// 必须，设置服务器支持的所有跨域请求的方法
//		context.Header("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
//		//context.Header("Access-Control-Allow-Methods", "*")
//		// 服务器支持的所有头信息字段，不限于浏览器在"预检"中请求的字段
//		context.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Token")
//		//context.Header("Access-Control-Allow-Headers", "*")
//		// 可选，设置XMLHttpRequest的响应对象能拿到的额外字段
//		context.Header("Access-Control-Expose-Headers", "Access-Control-Allow-Headers, Token")
//		// 可选，是否允许后续请求携带认证信息Cookie，该值只能是true，不需要则不设置
//		context.Header("Access-Control-Allow-Credentials", "true")
//		// 放行所有OPTIONS方法
//		if method == "OPTIONS" {
//			context.AbortWithStatus(http.StatusNoContent)
//			return
//		}
//		context.Next()
//	}
//}
