package dal

import (
	"errors"
	"fmt"
)

// AddCode 向information表中添加新记录
func AddCode(code string, idA int64) (bool, error) {
	info := Information{}
	info.Id = idA
	info.Code = code
	result := DB.Model(&Information{}).Create(&info)
	if result.Error != nil {
		return false, errors.New("数据库添加溯源记录失败")
	}
	return true, nil
}

// GetIdA 通过溯源码获取idA
func GetIdA(code string) (int64, error) {
	info := Information{}
	fmt.Println("code", code)
	result := DB.Model(&Information{}).Where("code = ?", code).First(&info)
	fmt.Println("info", info)
	if result.Error != nil {
		return -1, errors.New("溯源码有误")
	}
	return info.Id, nil
}
