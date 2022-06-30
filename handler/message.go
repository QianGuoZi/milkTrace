package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"server/dal"
	"server/service"
)

type Input struct {
	Code    string       `form:"code" json:"code"`
	Message FactoryInput `form:"message" json:"message"`
}

type RanchInput struct {
	BatchID string `form:"batchId" json:"batchId"` // 批次号
	Date    string `form:"date" json:"date"`       // 产奶日期
	Weight  int64  `form:"weight" json:"weight"`   // 总净重
}

type FactoryInput struct {
	BatchID     string `form:"batchId" json:"batchId"`         // 批次号
	CheckDate   string `form:"checkDate" json:"checkDate"`     // 抽检日期
	CheckPerson string `form:"checkPerson" json:"checkPerson"` // 抽检人姓名
	Material    string `form:"material" json:"material"`       // 产品成分
	Product     string `form:"product" json:"product"`         // 产品名称
	WorkDate    string `form:"workDate" json:"workDate"`       // 加工日期
	WorkPerson  string `form:"workPerson" json:"workPerson"`   // 加工人姓名
}

type StorageInput struct {
	BatchID string `form:"batchId" json:"batchId"` // 批次号
	Driver  string `form:"driver" json:"driver"`   // 运输负责人
}

type SellerInput struct {
	BatchID string `form:"batchId" json:"batchId"` // 批次号
	Date    string `form:"date" json:"date"`       // 上架时间
	Price   int64  `form:"price" json:"price"`     // 商品售价
}

func GetMessage(c *gin.Context) {
	// 根据 token 获得用户名
	username, err := service.GetUsername(c)
	if err != nil {
		log.Printf("[GetUserInfo] failed err=%+v", err)
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	// 根据用户名获取用户信息
	user, err := dal.GetUserInfoByName(username)
	if err != nil {
		log.Printf("[GetUserInfo] failed err=%+v", err)
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	if user.Role == "0" {
		//牧场获取信息
		resultList, err := service.GetInfoRanch(user.Id)
		if err != nil {
			log.Printf("[GetInfoRanch] failed err=%+v", err)
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "成功获取产品信息",
			"data":    resultList,
		})
	} else if user.Role == "1" {
		//加工厂获取信息
		resultList, err := service.GetInfoFactory(user.Id)
		if err != nil {
			log.Printf("[GetInfoRanch] failed err=%+v", err)
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "成功获取产品信息",
			"data":    resultList,
		})
	} else if user.Role == "2" {
		//储运商获取信息
		resultList, err := service.GetInfoStorage(user.Id)
		if err != nil {
			log.Printf("[GetInfoStorage] failed err=%+v", err)
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "成功获取产品信息",
			"data":    resultList,
		})
	} else if user.Role == "3" {
		//储运商获取信息
		resultList, err := service.GetInfoSeller(user.Id)
		if err != nil {
			log.Printf("[GetInfoSeller] failed err=%+v", err)
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "成功获取产品信息",
			"data":    resultList,
		})
	}

}

func SetMessage(c *gin.Context) {
	// 根据 token 获得用户名
	username, err := service.GetUsername(c)
	if err != nil {
		log.Printf("[GetUserInfo] failed err=%+v", err)
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	// 根据用户名获取用户信息
	user, err := dal.GetUserInfoByName(username)
	if err != nil {
		log.Printf("[GetUserInfo] failed err=%+v", err)
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	obj := struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}{}
	err = c.BindJSON(&obj)
	if err != nil {
		fmt.Println("json数据获取失败", err)
	} else {
		fmt.Println("获取的数据为", obj)
	}
	code := obj.Code

	if user.Role == "0" {
		var ranchInput = RanchInput{}
		fmt.Println("message", obj.Message)
		json.Unmarshal([]byte(obj.Message), &ranchInput)
		fmt.Println("ranchInput", ranchInput)

		err = service.AddInfoRanch(user.Id, ranchInput.BatchID, ranchInput.Date, ranchInput.Weight)
		if err != nil {
			log.Printf("[AddInfoRanch] failed err=%+v", err)
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "成功添加产品信息",
		})
	} else if user.Role == "1" {
		var factoryInput = FactoryInput{}
		fmt.Println("message", obj.Message)
		json.Unmarshal([]byte(obj.Message), &factoryInput)
		fmt.Println("factoryInput", factoryInput)

		err = service.AddInfoFactory(code, user.Id, factoryInput.BatchID,
			factoryInput.CheckDate, factoryInput.CheckPerson, factoryInput.Material,
			factoryInput.Product, factoryInput.WorkDate, factoryInput.WorkPerson)
		if err != nil {
			log.Printf("[AddInfoFactory] failed err=%+v", err)
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "成功添加产品信息",
		})
	} else if user.Role == "2" {
		var storageInput = StorageInput{}
		fmt.Println("message", obj.Message)
		json.Unmarshal([]byte(obj.Message), &storageInput)
		fmt.Println("storageInput", storageInput)

		err = service.AddInfoStorage(code, user.Id, storageInput.BatchID, storageInput.Driver)
		if err != nil {
			log.Printf("[AddInfoStorage] failed err=%+v", err)
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "成功添加产品信息",
		})
	} else if user.Role == "3" {
		var sellerInput = SellerInput{}
		fmt.Println("message", obj.Message)
		json.Unmarshal([]byte(obj.Message), &sellerInput)
		fmt.Println("sellerInput", sellerInput)

		err = service.AddInfoSeller(code, user.Id, sellerInput.BatchID, int(sellerInput.Price), sellerInput.Date)
		if err != nil {
			log.Printf("[AddInfoSeller] failed err=%+v", err)
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "成功添加产品信息",
		})
	}
}

func getRequestBody(context *gin.Context, s interface{}) error { //获取request的body
	body, _ := context.Get("json") //转换成json格式
	reqBody, _ := body.(string)
	decoder := json.NewDecoder(bytes.NewReader([]byte(reqBody)))
	decoder.UseNumber()       //作为数字而不是float64
	err := decoder.Decode(&s) //从body中获取的参数存入s中
	return err
}
