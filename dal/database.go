package dal

import (
	"log"

	"github.com/go-redis/redis"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var rdb *redis.Client

// InitDB 初始化数据库
func InitDB() {
	var err error
	dsn := "root:12345678@tcp(127.0.0.1:3306)/" +
		"milkTrace?charset=utf8mb4&interpolateParams=true&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		PrepareStmt: true,
	})
	if err != nil {
		panic(err)
	}
	err = DB.AutoMigrate(&User{}, &Information{})
	log.Println(err)
}

func InitRedis() {
	rdb = redis.NewClient(&redis.Options{Addr: "localhost:6379", Password: "", DB: 1})
	_, err := rdb.Ping().Result()
	if err != nil {
		log.Print("[InitRedis] failed")
	}
}
