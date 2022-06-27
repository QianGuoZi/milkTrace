package dal

import (
	"errors"
	"fmt"
)

// AddCode 向information表中添加新记录
func AddCode(code string, idA int64) error {
	info := Information{}
	info.Id = idA
	info.Code = code
	result := DB.Model(&Information{}).Create(&info)
	if result.Error != nil {
		return errors.New("数据库添加溯源记录失败")
	}
	fmt.Println("插入的溯源码info", info)
	return nil
}

// GetCode 获取溯源码code
func GetCode(idA int64) (string, error) {
	info := Information{}
	result := DB.Model(&Information{}).Where("id = ?", idA).First(&info)
	if result.Error != nil {
		return "", errors.New("获取溯源码失败")
	}
	fmt.Println("获取的溯源码为：", info.Code)
	return info.Code, nil
}
