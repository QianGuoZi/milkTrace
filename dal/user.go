package dal

import (
	"errors"
	"fmt"
	"log"

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
		fmt.Print("用户名错误")
		return false, errors.New("用户名错误")
	}
	err2 := bcrypt.CompareHashAndPassword([]byte(user.Pwd), []byte(userName+password+user.Salt))
	if err2 != nil {
		fmt.Print("密码错误")
		return false, errors.New("密码错误")
	}
	return true, nil
}

// GetUserInfo 获取用户信息
func GetUserInfo(id int64) (User, error) {
	user := User{}
	DB.Model(&User{}).Where("id = ?", id).First(&user)
	if user.Id == 0 {
		return User{}, errors.New("无法找到该用户")
	}
	//屏蔽掉密码等信息
	user.Pwd = ""
	user.Salt = ""
	fmt.Println("getUserInfo", user)
	return user, nil
}

// GetUserInfoByName 获取用户信息
func GetUserInfoByName(name string) (User, error) {
	user, err := GetUserByRedis(name)

	if err != nil {
		return User{}, err
	}

	if user == nil {
		DB.Model(&User{}).Where("user_name = ?", name).First(user)
		SetUserInRedis(user)
	}

	if user.Id == 0 {
		return User{}, errors.New("无法找到该用户")
	}
	//屏蔽掉密码等信息
	user.Pwd = ""
	user.Salt = ""
	log.Printf("GetUserInfoByName usr=%+v", *user)
	return *user, nil
}

// UpdateUser 更改用户信息
func UpdateUser(user *User) error {
	info := make(map[string]interface{})
	if user.Company != "" {
		info["company"] = user.Company
	}
	if user.Phone != "" {
		info["phone"] = user.Phone
	}
	if user.Address != "" {
		info["address"] = user.Address
	}
	if user.Pwd != "" {
		info["pwd"] = user.Pwd
	}
	log.Printf("[mysql UpdateUser] user=%+v data=%+v", user, info)
	result := DB.Model(&User{}).Where("user_name = ?", user.UserName).Updates(info)

	if result.Error != nil {
		log.Printf("mysql update user failed err=%+v", result.Error)
		return result.Error
	}
	SetUserInRedis(user)

	return nil
}

func GetUserByRedis(username string) (*User, error) {
	userKey := "user:" + username
	company, err1 := rdb.HGet(userKey, "company").Result()
	address, err2 := rdb.HGet(userKey, "address").Result()
	phone, err3 := rdb.HGet(userKey, "phone").Result()

	if err1 != nil || err2 != nil || err3 != nil {
		return nil, err1
	}

	user := &User{
		UserName: username,
		Company:  company,
		Address:  address,
		Phone:    phone,
	}

	return user, nil
}

func SetUserInRedis(user *User) error {
	userKey := "user:" + user.UserName
	err := rdb.HSet(userKey, "company", user.Company).Err()
	if err != nil {
		return err
	}
	err = rdb.HSet(userKey, "address", user.Address).Err()
	if err != nil {
		return err
	}
	err = rdb.HSet(userKey, "phone", user.Phone).Err()
	if err != nil {
		return err
	}
	return nil
}
