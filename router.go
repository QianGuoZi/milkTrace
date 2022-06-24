package main

import (
	"github.com/gin-gonic/gin"
	"server/handler"
)

func initRouter(r *gin.Engine) {
	apiRouter := r.Group("/milkTrace")

	//用户 apis
	{
		apiRouter.POST("/login/", handler.Login)
		apiRouter.POST("/register/", handler.Register)
		apiRouter.POST("/logout/", handler.Logout)
		apiRouter.POST("/getUsername/", handler.GetUsername)
	}
	//溯源信息 apis
	{
		apiRouter.GET("/getTrace/", handler.GetTrace)
		apiRouter.POST("")
	}
}
