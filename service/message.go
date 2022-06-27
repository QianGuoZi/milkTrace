package service

import (
	"errors"
	"fmt"
	"math/big"
	"server/dal"
	"strconv"
)

type MessageData struct {
	Code    string   `form:"code" json:"code"`       //溯源码
	Ranch   *Ranch   `form:"ranch" json:"ranch"`     // 牧场
	Factory *Factory `form:"factory" json:"factory"` // 加工厂
	Storage *Storage `form:"storage" json:"storage"` // 储运商
	Seller  *Seller  `form:"seller" json:"seller"`   // 销售商
}

// GetInfoRanch 牧场获取数据
func GetInfoRanch(ranchId int64) ([]MessageData, error) {
	//将牧场的userId转成string
	ranchIdStr := strconv.FormatInt(ranchId, 10)
	//
	idAList, batchRawList, weightRawList, dateRawList, err := dal.TlsApi.GetByPasIdPasture(ranchIdStr)
	if err != nil {
		return []MessageData{}, errors.New("合约获取牧场已发布的数据失败")
	}
	fmt.Println("牧场id", ranchId)
	batchUser, err1 := dal.GetUserInfo(int64(ranchId))
	fmt.Println("牧场信息", batchUser)
	if err1 != nil {
		return []MessageData{}, errors.New("获取牧场信息失败")
	}

	arrayLength := len(idAList)
	dataArray := make([]MessageData, arrayLength)

	for i := 0; i < arrayLength; i++ {
		//获取溯源码
		code, err := dal.GetCode(idAList[i].Int64())
		if err != nil {
			return []MessageData{}, errors.New("获取溯源码失败")
		}
		dataArray[i].Code = code

		//获取牧场信息
		ranch := Ranch{}
		ranch.BatchID = batchRawList[i]
		ranch.Date = dateRawList[i]
		ranch.Weight = weightRawList[i].Int64()
		ranch.Company = batchUser.Company
		ranch.Phone = batchUser.Phone
		ranch.Address = batchUser.Address
		dataArray[i].Ranch = &ranch
		dataArray[i].Factory = nil
		dataArray[i].Storage = nil
		dataArray[i].Seller = nil
	}
	fmt.Println("牧场返回的信息", dataArray)
	return dataArray, nil
}

// AddInfoRanch 牧场添加信息
func AddInfoRanch(userId int64, batchId string, date string, weight int64) error {
	//合约添加信息
	userIdStr := strconv.FormatInt(userId, 10)
	weightBInt := new(big.Int).SetUint64(uint64(int(weight)))
	address, _, err := dal.TlsApi.SetPasture(batchId, weightBInt, date, userIdStr)
	if err != nil {
		return errors.New("合约添加牧场信息失败")
	}

	//获取溯源码
	code := address.Hash().Hex()
	fmt.Println("添加牧场信息的溯源码为", code)

	//获取idA
	_idA, err := dal.TlsApi.GetId("idA")
	if err != nil {
		return errors.New("获取idA失败")
	}
	idA := _idA.Int64()

	//information 插入信息
	err = dal.AddCode(code, idA)
	if err != nil {
		return errors.New("数据库information插入信息失败")
	}
	return nil
}
