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

	err := r.Run(":8080") // http默认端口
	if err != nil {
		panic(err)
	}

	//log.Fatal(err.Error())
	//dal.InitDB()
	//register, err := service.Register("a", "123", "1")
	//if err != nil {
	//	fmt.Println("err", err)
	//	return
	//}
	//fmt.Println("id", register)
}
