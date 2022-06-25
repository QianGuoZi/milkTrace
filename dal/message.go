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
	fmt.Println("插入的溯源码info", info)
	return true, nil
}
