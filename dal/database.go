package dal

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
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

// CheckUser 查看用户是否已存在
func CheckUser(name string) (User, error) {
	user := User{}
	DB.Model(&User{}).Where("user_name = ?", name).First(&user)
	if user.Id != 0 {
		return user, errors.New("用户已存在")
	}
	return user, nil
}

// AddUser 向users表中添加新记录
func AddUser(user User) (int64, error) {
	result := DB.Model(&User{}).Create(&user)
	if result.Error != nil {
		return 0, errors.New("数据库创建用户记录失败")
	}
	return user.Id, nil
}

func SearchUser(username, password, role string) (bool, error) {
	user := User{}
	DB.Model(&User{}).Where("user_name = ? && role = ?", username, role).First(&user)
	fmt.Println("查询的user", user)
	if user.UserName != username {
		fmt.Println("用户名错误", user)
		return false, errors.New("用户名错误")
	}
	err2 := bcrypt.CompareHashAndPassword([]byte(user.Pwd), []byte(username+password+user.Salt))
	if err2 != nil {
		fmt.Println("密码错误")
		return false, errors.New("密码错误")
	}
	return true, nil
}