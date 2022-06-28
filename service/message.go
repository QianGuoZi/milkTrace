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

// GetInfoRanch 牧场获取信息
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
		code, err1 := dal.GetCode(idAList[i].Int64())
		if err1 != nil {
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
	_idA, err1 := dal.TlsApi.GetId("idA")
	if err1 != nil {
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

// GetInfoFactory 加工厂获取信息
func GetInfoFactory(factoryId int64) ([]MessageData, error) {

	//未添加加工厂信息的牧场信息
	idAListEmpty, err2 := dal.TlsApi.GetEmptyPasture()
	if err2 != nil {
		return []MessageData{}, errors.New("获取未补充加工厂信息的牧场信息失败")
	}
	//获取牧场信息列表
	batchRawListEmpty, weightRawListEmpty, dateRawListEmpty, ranchIdListEmpty, err := dal.TlsApi.GetByIdAPasture(idAListEmpty)

	//加工厂发布的信息
	//加工厂id字符串
	factoryIdStr := strconv.FormatInt(factoryId, 10)
	//牧场信息idA列表
	idAList, err := dal.TlsApi.GetByFacIdTrace1(factoryIdStr)
	if err != nil {
		return []MessageData{}, errors.New("合约获取牧场idA列表失败")
	}
	fmt.Println("加工厂发布的信息中牧场信息idA列表：", idAList)

	//获取牧场信息列表
	batchRawList, weightRawList, dateRawList, ranchIdList, err := dal.TlsApi.GetByIdAPasture(idAList)

	//加工厂信息列表
	_, batchProList, productNameList, err := dal.TlsApi.GetByFacIdFactory1(factoryIdStr)
	if err != nil {
		return []MessageData{}, errors.New("合约获取加工厂已发布的数据1失败")
	}
	checkerList, checkTimeList, err := dal.TlsApi.GetByFacIdFactory2(factoryIdStr)
	if err != nil {
		return []MessageData{}, errors.New("合约获取加工厂已发布的数据2失败")
	}
	processorList, processTimeList, compositionList, err := dal.TlsApi.GetByFacIdFactory3(factoryIdStr)
	if err != nil {
		return []MessageData{}, errors.New("合约获取加工厂已发布的数据3失败")
	}

	fmt.Println("加工厂id", factoryId)
	factoryUser, err1 := dal.GetUserInfo(int64(factoryId))
	fmt.Println("加工厂信息", factoryUser)
	if err1 != nil {
		return []MessageData{}, errors.New("获取加工厂信息失败")
	}

	arrayLength := len(idAList) + len(idAListEmpty)
	fmt.Println("arrayLength", arrayLength)
	dataArray := make([]MessageData, arrayLength)

	//牧场+__
	for i := 0; i < len(idAListEmpty); i++ {
		//获取溯源码
		code, err := dal.GetCode(idAListEmpty[i].Int64())
		if err != nil {
			return []MessageData{}, errors.New("获取溯源码失败")
		}
		dataArray[i].Code = code

		//获取牧场信息
		ranchId, _ := strconv.Atoi(ranchIdListEmpty[i])
		fmt.Println("牧场id", ranchId)
		//！！！！！！！！！！！！！！！！！！！！！！！！！！1是临时值要改回ranchId
		ranchUser, err1 := dal.GetUserInfo(int64(1))
		fmt.Println("牧场信息", ranchUser)
		if err1 != nil {
			return []MessageData{}, errors.New("获取牧场信息失败")
		}
		ranch := Ranch{}
		ranch.BatchID = batchRawListEmpty[i]
		ranch.Date = dateRawListEmpty[i]
		ranch.Weight = weightRawListEmpty[i].Int64()
		ranch.Company = ranchUser.Company
		ranch.Phone = ranchUser.Phone
		ranch.Address = ranchUser.Address
		dataArray[i].Ranch = &ranch

		dataArray[i].Factory = nil
		dataArray[i].Storage = nil
		dataArray[i].Seller = nil
	}

	//牧场+加工厂
	for i := len(idAListEmpty); i < arrayLength; i++ {
		var ii = i - len(idAListEmpty)
		//获取溯源码
		code, err := dal.GetCode(idAList[ii].Int64())
		if err != nil {
			return []MessageData{}, errors.New("获取溯源码失败")
		}
		dataArray[i].Code = code

		//获取牧场信息
		ranchId, _ := strconv.Atoi(ranchIdList[ii])
		fmt.Println("牧场id", ranchId)
		////！！！！！！！！！！！！！！！！！！！！！！！！！！1是临时值要改回ranchId
		ranchUser, err1 := dal.GetUserInfo(int64(1))
		fmt.Println("牧场信息", ranchUser)
		if err1 != nil {
			return []MessageData{}, errors.New("获取牧场信息失败")
		}
		ranch := Ranch{}
		ranch.BatchID = batchRawList[ii]
		ranch.Date = dateRawList[ii]
		ranch.Weight = weightRawList[ii].Int64()
		ranch.Company = ranchUser.Company
		ranch.Phone = ranchUser.Phone
		ranch.Address = ranchUser.Address
		dataArray[i].Ranch = &ranch

		factory := Factory{}
		factory.BatchID = batchProList[ii]
		factory.CheckDate = checkTimeList[ii]
		factory.CheckPerson = checkerList[ii]
		factory.Material = compositionList[ii]
		factory.Product = productNameList[ii]
		factory.WorkDate = processTimeList[ii]
		factory.WorkPerson = processorList[ii]
		factory.Company = factoryUser.Company
		factory.Phone = factoryUser.Phone
		factory.Address = factoryUser.Address
		dataArray[i].Factory = &factory

		dataArray[i].Storage = nil
		dataArray[i].Seller = nil
	}

	fmt.Println("加工厂返回的信息", dataArray)
	return dataArray, nil
}

// AddInfoFactory 加工厂添加信息
func AddInfoFactory(code string, factoryId int64, batchId string, checkDate string, checkPerson string,
	material string, product string, workDate string, workPerson string) error {
	//获取idA
	_idA, err := dal.GetIdA(code)
	if err != nil {
		return errors.New("溯源码获取idA失败")
	}
	idA := new(big.Int).SetUint64(uint64(int(_idA)))

	_, _, err = dal.TlsApi.SetFactory1(idA, batchId, product, material)
	if err != nil {
		return errors.New("加工厂发布数据1失败")
	}
	_, _, err = dal.TlsApi.SetFactory2(idA, checkPerson, checkDate)
	if err != nil {
		return errors.New("加工厂发布数据2失败")
	}
	_, _, err = dal.TlsApi.SetFactory3(idA, workPerson, workDate, strconv.Itoa(int(factoryId)))
	if err != nil {
		return errors.New("加工厂发布数据3失败")
	}
	return nil
}
