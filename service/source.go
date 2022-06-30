package service

import (
	"errors"
	"fmt"
	"log"
	"math/big"
	"server/dal"
	"strconv"
)

type Data struct {
	Ranch   Ranch   `form:"ranch" json:"ranch"`     // 牧场
	Factory Factory `form:"factory" json:"factory"` // 加工厂
	Storage Storage `form:"storage" json:"storage"` // 储运商
	Seller  Seller  `form:"seller" json:"seller"`   // 销售商
}

// Ranch 牧场
type Ranch struct {
	BatchID string `form:"batchId" json:"batchId"` // 批次号
	Date    string `form:"date" json:"date"`       // 产奶日期
	Weight  int64  `form:"weight" json:"weight"`   // 总净重
	Company string `form:"company" json:"company"` //公司名称
	Phone   string `form:"phone" json:"phone"`     //联系方式
	Address string `form:"address" json:"address"` //地址
}

// Factory 加工厂
type Factory struct {
	BatchID     string `form:"batchId" json:"batchId"`         // 批次号
	CheckDate   string `form:"checkDate" json:"checkDate"`     // 抽检日期
	CheckPerson string `form:"checkPerson" json:"checkPerson"` // 抽检人姓名
	Material    string `form:"material" json:"material"`       // 产品成分
	Product     string `form:"product" json:"product"`         // 产品名称
	WorkDate    string `form:"workDate" json:"workDate"`       // 加工日期
	WorkPerson  string `form:"workPerson" json:"workPerson"`   // 加工人姓名
	Company     string `form:"company" json:"company"`         //公司名称
	Phone       string `form:"phone" json:"phone"`             //联系方式
	Address     string `form:"address" json:"address"`         //地址
}

// Storage 储运商
type Storage struct {
	BatchID string `form:"batchId" json:"batchId"` // 批次号
	Driver  string `form:"driver" json:"driver"`   // 运输负责人
	Company string `form:"company" json:"company"` //公司名称
	Phone   string `form:"phone" json:"phone"`     //联系方式
	Address string `form:"address" json:"address"` //地址
}

// Seller 销售商
type Seller struct {
	BatchID string `form:"batchId" json:"batchId"` // 批次号
	Date    string `form:"date" json:"date"`       // 上架时间
	Price   int64  `form:"price" json:"price"`     // 商品售价
	Company string `form:"company" json:"company"` //公司名称
	Phone   string `form:"phone" json:"phone"`     //联系方式
	Address string `form:"address" json:"address"` //地址
}

// GetRanch 使用idA获取牧场信息
func GetRanch(idA *big.Int) (Ranch, error) {
	ranch := Ranch{}

	//使用合约
	_, batchRow, weightRaw, dateRaw, ranchIdStr, err := dal.TlsApi.GetPasture(idA)
	if err != nil {
		log.Fatal(err)
	}

	ranchId, _ := strconv.Atoi(ranchIdStr)
	fmt.Println("牧场id", ranchId)
	batchUser, err1 := dal.GetUserInfo(int64(ranchId))
	fmt.Println("牧场信息", batchUser)
	if err1 != nil {
		return Ranch{}, errors.New("获取牧场信息失败")
	}

	ranch.BatchID = batchRow
	ranch.Date = dateRaw
	ranch.Weight = weightRaw.Int64()
	ranch.Company = batchUser.Company
	ranch.Phone = batchUser.Phone
	ranch.Address = batchUser.Address
	fmt.Println("牧场返回的信息", ranch)
	return ranch, nil
}

// GetFactory 使用idA获取加工场信息
func GetFactory(idA *big.Int) (Factory, error) {
	factory := Factory{}
	//加工厂信息
	idB, err2 := dal.TlsApi.GetByidAId(idA, "idB")
	if err2 != nil {
		return Factory{}, errors.New("获取加工厂id失败")
	}

	_, batchPro, productName, composition, checker, checkTime,
		processor, processTime, factoryIdStr, err := dal.TlsApi.GetFactory(idB)
	if err != nil {
		log.Fatal(err)
	}

	factoryId, _ := strconv.Atoi(factoryIdStr)
	fmt.Println("加工厂id", factoryId)
	factoryUser, err1 := dal.GetUserInfo(int64(factoryId))
	fmt.Println("加工厂信息", factoryUser)
	if err1 != nil {
		return Factory{}, errors.New("获取加工厂信息失败")
	}

	factory.BatchID = batchPro
	factory.CheckDate = checkTime
	factory.CheckPerson = checker
	factory.Material = composition
	factory.Product = productName
	factory.WorkDate = processTime
	factory.WorkPerson = processor
	factory.Company = factoryUser.Company
	factory.Phone = factoryUser.Phone
	factory.Address = factoryUser.Address
	fmt.Println("加工厂返回的信息", factory)
	return factory, nil
}

// GetStorage 使用idA获取储运商信息
func GetStorage(idA *big.Int) (Storage, error) {
	storage := Storage{}
	//储运商信息
	idC, err := dal.TlsApi.GetByidAId(idA, "idC")
	if err != nil {
		return Storage{}, errors.New("获取储运商id失败")
	}

	_, batchLog, transName, logisticsIdStr, err1 := dal.TlsApi.GetLogistics(idC)
	if err1 != nil {
		log.Fatal(err1)
	}

	logisticsId, _ := strconv.Atoi(logisticsIdStr)
	fmt.Println("储运商id", logisticsId)
	logisticUser, err1 := dal.GetUserInfo(int64(logisticsId))
	fmt.Println("储运商信息", logisticUser)
	if err1 != nil {
		return Storage{}, errors.New("获取储运商信息失败")
	}

	storage.BatchID = batchLog
	storage.Driver = transName
	storage.Company = logisticUser.Company
	storage.Phone = logisticUser.Phone
	storage.Address = logisticUser.Address
	fmt.Println("储运商返回的信息", storage)
	return storage, nil
}

// GetSeller 使用idA获取销售商信息
func GetSeller(idA *big.Int) (Seller, error) {
	seller := Seller{}
	//加工厂信息
	idD, err2 := dal.TlsApi.GetByidAId(idA, "idD")
	if err2 != nil {
		return Seller{}, errors.New("获取销售商id失败")
	}

	_, batchSale, priceStr, salesTime, salesIdStr, err := dal.TlsApi.GetSales(idD)
	if err != nil {
		log.Fatal(err)
	}
	price, _ := strconv.Atoi(priceStr)

	salesId, _ := strconv.Atoi(salesIdStr)
	fmt.Println("销售商id", salesIdStr)
	salesUser, err1 := dal.GetUserInfo(int64(salesId))
	fmt.Println("销售商信息", salesUser)
	if err1 != nil {
		return Seller{}, errors.New("获取销售商信息失败")
	}

	seller.BatchID = batchSale
	seller.Price = int64(price)
	seller.Date = salesTime
	seller.Company = salesUser.Company
	seller.Phone = salesUser.Phone
	seller.Address = salesUser.Address
	fmt.Println("销售商返回的信息", seller)
	return seller, nil
}

// GetByCode 使用溯源码获取四类信息
func GetByCode(code string) (Data, error) {
	ranch := Ranch{}
	factory := Factory{}
	storage := Storage{}
	seller := Seller{}
	_id, err := dal.GetIdA(code)
	if err != nil {
		return Data{}, errors.New("溯源码有误")
	}
	idA := new(big.Int).SetUint64(uint64(_id))
	ranch, err = GetRanch(idA)
	if err != nil {
		return Data{}, errors.New("获取牧场信息有误")
	}
	factory, err = GetFactory(idA)
	if err != nil {
		return Data{}, errors.New("获取加工厂信息有误")
	}
	storage, err = GetStorage(idA)
	if err != nil {
		return Data{}, errors.New("获取储运商信息有误")
	}
	seller, err = GetSeller(idA)
	if err != nil {
		return Data{}, errors.New("获取销售商信息有误")
	}

	data := Data{}
	data.Ranch = ranch
	data.Factory = factory
	data.Storage = storage
	data.Seller = seller
	fmt.Println("data", data)
	return data, nil
}
