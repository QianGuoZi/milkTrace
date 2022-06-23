package main

import (
	"github.com/gin-gonic/gin"
	"math/rand"
	"server/dal"
	"time"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	rand.Seed(time.Now().UnixNano())

	dal.InitDB()

	r := gin.Default()

	initRouter(r)

	err := r.Run(":8080") // http端口
	if err != nil {
		panic(err)
	}

}
