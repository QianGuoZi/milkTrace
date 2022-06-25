package service

import (
	"errors"
	"fmt"
	"server/dal"
	"strconv"
)

type MessageData struct {
	Code    string  `form:"code" json:"code"`       //溯源码
	Ranch   Ranch   `form:"ranch" json:"ranch"`     // 牧场
	Factory Factory `form:"factory" json:"factory"` // 加工厂
	Storage Storage `form:"storage" json:"storage"` // 储运商
	Seller  Seller  `form:"seller" json:"seller"`   // 销售商
}

// GetInfoRanch 牧场获取数据
func GetInfoRanch(ranchId int64) ([]Data, error) {
	ranchIdStr := strconv.FormatInt(ranchId, 10)
	idAList, batchRawList, weightRawList, dateRawList, err := dal.TlsApi.GetByPasIdPasture(ranchIdStr)
	if err != nil {
		return []Data{}, errors.New("合约获取牧场已发布的数据失败")
	}
	fmt.Println("牧场id", ranchId)
	batchUser, err1 := dal.GetUserInfo(int64(ranchId))
	fmt.Println("牧场信息", batchUser)
	if err1 != nil {
		return []Data{}, errors.New("获取牧场信息失败")
	}

	arrayLength := len(idAList)
	dataArray := make([]Data, arrayLength)

	for i := 0; i < arrayLength; i++ {
		ranch := Ranch{}
		ranch.BatchID = batchRawList[i]
		ranch.Date = dateRawList[i]
		ranch.Weight = weightRawList[i].Int64()
		ranch.Company = batchUser.Company
		ranch.Phone = batchUser.Phone
		ranch.Address = batchUser.Address
		dataArray[i].Ranch = ranch
		dataArray[i].Factory = Factory{}
		dataArray[i].Storage = Storage{}
		dataArray[i].Seller = Seller{}
	}
	fmt.Println("牧场返回的信息", dataArray)
	return dataArray, nil
}

func AddInfoRanch() {

}
