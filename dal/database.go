package dal

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

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
