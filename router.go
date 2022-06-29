package main

import (
	"server/handler"

	"github.com/gin-gonic/gin"
)

func initRouter(r *gin.Engine) {
	apiRouter := r.Group("/milkTrace")

	//用户 apis
	{
		apiRouter.POST("/login", handler.Login)
		apiRouter.POST("/register", handler.Register)
		apiRouter.POST("/logout", handler.Logout)
		apiRouter.POST("/getUsername", handler.GetUsername)
		apiRouter.GET("/getUserInfo", handler.GetUserInfo)
		apiRouter.POST("/setUserInfo", handler.UpdateUserInfo)
		apiRouter.POST("/setPassword", handler.UpdateUserPwd)
	}
	//溯源信息 apis
	{
		apiRouter.GET("/getTrace", handler.GetTrace)
		apiRouter.GET("/getMessage", handler.GetMessage)
		apiRouter.POST("/setMessage", handler.SetMessage)
	}
}
