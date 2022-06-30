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
	batchUser, err1 := dal.GetUserInfo(ranchId)
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
	factoryUser, err1 := dal.GetUserInfo(factoryId)
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
		ranchUser, err1 := dal.GetUserInfo(int64(ranchId))
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
		ranchUser, err1 := dal.GetUserInfo(int64(ranchId))
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
	//dataArray[arrayLength-1].Ranch = nil
	//dataArray[arrayLength-1].Factory = nil
	//dataArray[arrayLength-1].Storage = nil
	//dataArray[arrayLength-1].Seller = nil
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

// GetInfoStorage 储运商获取信息
func GetInfoStorage(storageId int64) ([]MessageData, error) {
	//牧场+加工厂+__
	//牧场idList和加工厂idList
	idAListEmpty, idBListEmpty, err := dal.TlsApi.GetEmptyFactory()
	if err != nil {
		return []MessageData{}, errors.New("获取未添加储运商信息的idA和idB失败")
	}
	//获取牧场信息列表
	batchRawListEmpty, weightRawListEmpty, dateRawListEmpty, ranchIdListEmpty, err1 := dal.TlsApi.GetByIdAPasture(idAListEmpty)
	if err1 != nil {
		return []MessageData{}, errors.New("获取牧场信息列表失败")
	}
	//获取加工厂信息列表
	factoryIdListEmpty, batchProListEmpty, productNameListEmpty, err := dal.TlsApi.GetByIdBFactory1(idBListEmpty)
	if err != nil {
		return []MessageData{}, errors.New("获取加工厂信息1失败")
	}
	checkerListEmpty, checkTimeListEmpty, err := dal.TlsApi.GetByIdBFactory2(idBListEmpty)
	if err != nil {
		return []MessageData{}, errors.New("获取加工厂信息2失败")
	}
	processorListEmpty, processTimeListEmpty, compositionListEmpty, err := dal.TlsApi.GetByIdBFactory3(idBListEmpty)
	if err != nil {
		return []MessageData{}, errors.New("获取加工厂信息3失败")
	}

	//牧场+加工厂+储运商
	//储运商id字符串
	storageIdStr := strconv.FormatInt(storageId, 10)
	_, _, driverList, batchStoList, err := dal.TlsApi.GetByLogidLogistics(storageIdStr)
	if err != nil {
		return []MessageData{}, errors.New("获取储运商已添加的信息失败")
	}
	fmt.Println("储运商id", storageId)
	storageUser, err1 := dal.GetUserInfo(storageId)
	fmt.Println("储运商信息", storageUser)
	if err1 != nil {
		return []MessageData{}, errors.New("获取储运商信息失败")
	}

	//牧场idList
	idAList, err := dal.TlsApi.GetByLogIdTrace1(storageIdStr)
	if err != nil {
		return []MessageData{}, errors.New("获取储运商已添加的信息中idAList失败")
	}
	//获取牧场信息列表
	batchRawList, weightRawList, dateRawList, ranchIdList, err := dal.TlsApi.GetByIdAPasture(idAList)

	//加工厂idList
	idBList, err := dal.TlsApi.GetByLogIdTrace2(storageIdStr)
	if err != nil {
		return []MessageData{}, errors.New("获取储运商已添加的信息中idBList失败")
	}
	//获取加工厂信息列表
	factoryIdList, batchProList, productNameList, err := dal.TlsApi.GetByIdBFactory1(idBList)
	if err != nil {
		return []MessageData{}, errors.New("获取储运商已添加的信息中Factory1失败")
	}
	checkerList, checkTimeList, err := dal.TlsApi.GetByIdBFactory2(idBList)
	if err != nil {
		return []MessageData{}, errors.New("获取储运商已添加的信息中Factory2失败")
	}
	processorList, processTimeList, compositionList, err := dal.TlsApi.GetByIdBFactory3(idBList)
	if err != nil {
		return []MessageData{}, errors.New("获取储运商已添加的信息中Factory3失败")
	}

	arrayLength := len(idAList) + len(idAListEmpty)
	fmt.Println("arrayLength", arrayLength)
	dataArray := make([]MessageData, arrayLength)

	//牧场+加工厂+__
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
		ranchUser, err1 := dal.GetUserInfo(int64(ranchId))
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

		//获取加工厂信息
		factoryId, _ := strconv.Atoi(factoryIdListEmpty[i])
		fmt.Println("加工厂id", factoryId)
		//！！！！！！！！！！！！！！！！！！！！！！！！！！1是临时值要改回factoryId
		factoryUser, err1 := dal.GetUserInfo(int64(factoryId))
		fmt.Println("加工厂信息", factoryUser)
		if err1 != nil {
			return []MessageData{}, errors.New("获取加工厂信息失败")
		}
		factory := Factory{}
		factory.BatchID = batchProListEmpty[i]
		factory.CheckDate = checkTimeListEmpty[i]
		factory.CheckPerson = checkerListEmpty[i]
		factory.Material = compositionListEmpty[i]
		factory.Product = productNameListEmpty[i]
		factory.WorkDate = processTimeListEmpty[i]
		factory.WorkPerson = processorListEmpty[i]
		factory.Company = factoryUser.Company
		factory.Phone = factoryUser.Phone
		factory.Address = factoryUser.Address
		dataArray[i].Factory = &factory

		dataArray[i].Storage = nil
		dataArray[i].Seller = nil
	}

	//牧场+加工厂+储运商
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
		ranchUser, err1 := dal.GetUserInfo(int64(ranchId))
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

		//获取加工厂信息
		factoryId, _ := strconv.Atoi(factoryIdList[ii])
		fmt.Println("加工厂id", factoryId)
		//！！！！！！！！！！！！！！！！！！！！！！！！！！1是临时值要改回factoryId
		factoryUser, err1 := dal.GetUserInfo(int64(factoryId))
		fmt.Println("加工厂信息", factoryUser)
		if err1 != nil {
			return []MessageData{}, errors.New("获取加工厂信息失败")
		}
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

		//储运商信息
		storage := Storage{}
		storage.BatchID = batchStoList[ii]
		storage.Driver = driverList[ii]
		storage.Company = storageUser.Company
		storage.Phone = storageUser.Phone
		storage.Address = storageUser.Address
		dataArray[i].Storage = &storage

		dataArray[i].Seller = nil
	}
	//dataArray[arrayLength-1].Ranch = nil
	//dataArray[arrayLength-1].Factory = nil
	//dataArray[arrayLength-1].Storage = nil
	//dataArray[arrayLength-1].Seller = nil
	fmt.Println("储运商返回的信息", dataArray)
	return dataArray, nil
}

// AddInfoStorage 储运商添加信息
func AddInfoStorage(code string, storageId int64, batchId string, driver string) error {
	//获取idA
	_idA, err := dal.GetIdA(code)
	if err != nil {
		return errors.New("溯源码获取idA失败")
	}
	idA := new(big.Int).SetUint64(uint64(int(_idA)))
	_, _, err = dal.TlsApi.SetLogistics(idA, batchId, driver, strconv.Itoa(int(storageId)))
	if err != nil {
		return errors.New("储运商发布数据失败")
	}
	return nil
}

// GetInfoSeller 销售商获取信息
func GetInfoSeller(sellerId int64) ([]MessageData, error) {
	//牧场+加工厂+储运商+__
	//牧场idList和加工厂idList和储运商idList
	idAListEmpty, idBListEmpty, idCListEmpty, err := dal.TlsApi.GetEmptyLogistics()
	if err != nil {
		return []MessageData{}, errors.New("获取未添加储运商信息的idA和idB失败")
	}
	//获取牧场信息列表
	batchRawListEmpty, weightRawListEmpty, dateRawListEmpty, ranchIdListEmpty, err1 := dal.TlsApi.GetByIdAPasture(idAListEmpty)
	if err1 != nil {
		return []MessageData{}, errors.New("获取牧场信息列表失败")
	}
	//获取加工厂信息列表
	factoryIdListEmpty, batchProListEmpty, productNameListEmpty, err := dal.TlsApi.GetByIdBFactory1(idBListEmpty)
	if err != nil {
		return []MessageData{}, errors.New("获取加工厂信息1失败")
	}
	checkerListEmpty, checkTimeListEmpty, err := dal.TlsApi.GetByIdBFactory2(idBListEmpty)
	if err != nil {
		return []MessageData{}, errors.New("获取加工厂信息2失败")
	}
	processorListEmpty, processTimeListEmpty, compositionListEmpty, err := dal.TlsApi.GetByIdBFactory3(idBListEmpty)
	if err != nil {
		return []MessageData{}, errors.New("获取加工厂信息3失败")
	}
	//获取储运商信息列表
	_, batchStoListEmpty, driverListEmpty, storageIdListEmpty, err := dal.TlsApi.GetByIdCLogistics(idCListEmpty)
	if err != nil {
		return []MessageData{}, errors.New("获取储运商失败")
	}

	//牧场+加工厂+储运商+销售商
	sellerIdStr := strconv.FormatInt(sellerId, 10)
	_, priceList, salesTimeList, batchSaleList, err := dal.TlsApi.GetBySaleidSales(sellerIdStr)
	if err != nil {
		return []MessageData{}, errors.New("获取销售商已添加的信息失败")
	}
	fmt.Println("销售商id", sellerId)
	sellerUser, err1 := dal.GetUserInfo(sellerId)
	fmt.Println("销售商信息", sellerUser)
	if err1 != nil {
		return []MessageData{}, errors.New("获取销售商信息失败")
	}

	//牧场idList
	idAList, err := dal.TlsApi.GetBySaleIdTrace1(sellerIdStr)
	if err != nil {
		return []MessageData{}, errors.New("获取销售商已添加的信息中idAList失败")
	}
	//获取牧场信息列表
	batchRawList, weightRawList, dateRawList, ranchIdList, err := dal.TlsApi.GetByIdAPasture(idAList)
	fmt.Print("获取牧场信息结束")
	//加工厂idList
	idBList, err := dal.TlsApi.GetBySaleIdTrace2(sellerIdStr)
	if err != nil {
		return []MessageData{}, errors.New("获取销售商已添加的信息中idBList失败")
	}
	//获取加工厂信息列表
	factoryIdList, batchProList, productNameList, err := dal.TlsApi.GetByIdBFactory1(idBList)
	if err != nil {
		return []MessageData{}, errors.New("获取销售商已添加的信息中Factory1失败")
	}
	checkerList, checkTimeList, err := dal.TlsApi.GetByIdBFactory2(idBList)
	if err != nil {
		return []MessageData{}, errors.New("获取销售商已添加的信息中Factory2失败")
	}
	processorList, processTimeList, compositionList, err := dal.TlsApi.GetByIdBFactory3(idBList)
	if err != nil {
		return []MessageData{}, errors.New("获取销售商已添加的信息中Factory3失败")
	}
	fmt.Print("获取加工厂信息结束")

	idCList, err := dal.TlsApi.GetBySaleIdTrace3(sellerIdStr)
	if err != nil {
		return []MessageData{}, errors.New("获取销售商已添加的信息中idCList失败")
	}
	_, batchStoList, driverList, storageIdList, err := dal.TlsApi.GetByIdCLogistics(idCList)
	if err != nil {
		return []MessageData{}, errors.New("获取销售商已添加的信息中Storage失败")
	}
	fmt.Print("获取储运商信息结束")

	arrayLength := len(idAList) + len(idAListEmpty)
	fmt.Println("arrayLength", arrayLength)
	dataArray := make([]MessageData, arrayLength)

	//牧场+加工厂+储运商+__
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
		ranchUser, err1 := dal.GetUserInfo(int64(ranchId))
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

		//获取加工厂信息
		factoryId, _ := strconv.Atoi(factoryIdListEmpty[i])
		fmt.Println("加工厂id", factoryId)
		//！！！！！！！！！！！！！！！！！！！！！！！！！！1是临时值要改回factoryId
		factoryUser, err1 := dal.GetUserInfo(int64(factoryId))
		fmt.Println("加工厂信息", factoryUser)
		if err1 != nil {
			return []MessageData{}, errors.New("获取加工厂信息失败")
		}
		factory := Factory{}
		factory.BatchID = batchProListEmpty[i]
		factory.CheckDate = checkTimeListEmpty[i]
		factory.CheckPerson = checkerListEmpty[i]
		factory.Material = compositionListEmpty[i]
		factory.Product = productNameListEmpty[i]
		factory.WorkDate = processTimeListEmpty[i]
		factory.WorkPerson = processorListEmpty[i]
		factory.Company = factoryUser.Company
		factory.Phone = factoryUser.Phone
		factory.Address = factoryUser.Address
		dataArray[i].Factory = &factory

		//获取储运商信息
		storageId, _ := strconv.Atoi(storageIdListEmpty[i])
		fmt.Println("储运商id", storageId)
		//！！！！！！！！！！！！！！！！！！！！！！！！！！1是临时值要改回storageId
		storageUser, err1 := dal.GetUserInfo(int64(storageId))
		fmt.Println("储运商信息", storageUser)
		if err1 != nil {
			return []MessageData{}, errors.New("获取储运商信息失败")
		}
		storage := Storage{}
		storage.BatchID = batchStoListEmpty[i]
		storage.Driver = driverListEmpty[i]
		storage.Company = storageUser.Company
		storage.Phone = storageUser.Phone
		storage.Address = storageUser.Address
		dataArray[i].Storage = &storage

		dataArray[i].Seller = nil
	}

	//牧场+加工厂+储运商
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
		ranchUser, err1 := dal.GetUserInfo(int64(ranchId))
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

		//获取加工厂信息
		factoryId, _ := strconv.Atoi(factoryIdList[ii])
		fmt.Println("加工厂id", factoryId)
		//！！！！！！！！！！！！！！！！！！！！！！！！！！1是临时值要改回factoryId
		factoryUser, err1 := dal.GetUserInfo(int64(factoryId))
		fmt.Println("加工厂信息", factoryUser)
		if err1 != nil {
			return []MessageData{}, errors.New("获取加工厂信息失败")
		}
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

		//储运商信息
		storageId, _ := strconv.Atoi(storageIdList[ii])
		fmt.Println("储运商id", storageId)
		//！！！！！！！！！！！！！！！！！！！！！！！！！！1是临时值要改回storageId
		storageUser, err1 := dal.GetUserInfo(int64(storageId))
		fmt.Println("储运商信息", storageUser)
		if err1 != nil {
			return []MessageData{}, errors.New("获取储运商信息失败")
		}
		storage := Storage{}
		storage.BatchID = batchStoList[ii]
		storage.Driver = driverList[ii]
		storage.Company = storageUser.Company
		storage.Phone = storageUser.Phone
		storage.Address = storageUser.Address
		dataArray[i].Storage = &storage

		//销售商信息
		seller := Seller{}
		seller.BatchID = batchSaleList[ii]
		seller.Price = priceList[ii].Int64()
		seller.Date = salesTimeList[ii]
		seller.Company = sellerUser.Company
		seller.Phone = sellerUser.Phone
		seller.Address = sellerUser.Address

		dataArray[i].Seller = &seller
	}
	//dataArray[arrayLength-1].Ranch = nil
	//dataArray[arrayLength-1].Factory = nil
	//dataArray[arrayLength-1].Storage = nil
	//dataArray[arrayLength-1].Seller = nil
	fmt.Println("销售商返回的信息", dataArray)
	return dataArray, nil
}

// AddInfoSeller 销售商添加信息
func AddInfoSeller(code string, sellerId int64, batchId string, price int, salesTime string) error {
	//获取idA
	_idA, err := dal.GetIdA(code)
	if err != nil {
		return errors.New("溯源码获取idA失败")
	}
	idA := new(big.Int).SetUint64(uint64(int(_idA)))
	priceBigInt := new(big.Int).SetUint64(uint64(price))
	_, _, err = dal.TlsApi.SetSales(idA, batchId, priceBigInt, salesTime, strconv.Itoa(int(sellerId)))
	if err != nil {
		return errors.New("销售商添加数据失败")
	}

	return nil
}
