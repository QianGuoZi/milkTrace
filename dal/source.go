package dal

import (
	"errors"
	"fmt"
)

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
