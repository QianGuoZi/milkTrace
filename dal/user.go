package dal

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

// CheckUser 查看用户是否已存在
func CheckUser(userName string) (User, error) {
	user := User{}
	DB.Model(&User{}).Where("user_name = ?", userName).First(&user)
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

// SearchUser 判断userName password role是否正确
func SearchUser(userName, password, role string) (bool, error) {
	user := User{}
	DB.Model(&User{}).Where("user_name = ? && role = ?", userName, role).First(&user)
	fmt.Println("查询的user", user)
	if user.UserName != userName {
		return false, errors.New("用户名错误")
	}
	err2 := bcrypt.CompareHashAndPassword([]byte(user.Pwd), []byte(userName+password+user.Salt))
	if err2 != nil {
		return false, errors.New("密码错误")
	}
	return true, nil
}

func GetUserInfo(id int64) (User, error) {
	user := User{}
	DB.Model(&User{}).Where("id = ?", id).First(&user)
	fmt.Println("getUserInfo", user)
	if user.Id == 0 {
		return User{}, errors.New("无法找到该用户")
	}
	return user, nil
}
