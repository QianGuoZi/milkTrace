package main

import (
	"fmt"
	"github.com/FISCO-BCOS/go-sdk/client"
	"github.com/FISCO-BCOS/go-sdk/conf"
	"log"
	"server/store"
)

func main() {
	//配置文件
	configs, err := conf.ParseConfigFile("config.toml")
	if err != nil {
		log.Fatalf("ParseConfigFile failed, err: %v", err)
	}

	//用于访问节点
	client, err := client.Dial(&configs[0])
	if err != nil {
		log.Fatal(err)
	}

	//部署合约
	input := "Store deployment 1.0"
	address, tx, instance, err := store.DeployStore(client.GetTransactOpts(), client, input)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("contract address: ", address.Hex()) // the address should be saved, will use in next example
	fmt.Println("transaction hash: ", tx.Hash().Hex())

	// load the contract
	// contractAddress := common.HexToAddress("contract address in hex String")
	// instance, err := store.NewStore(contractAddress, client)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	//调用合约
	fmt.Println("================================")
	storeSession := &store.StoreSession{Contract: instance, CallOpts: *client.GetCallOpts(), TransactOpts: *client.GetTransactOpts()}

	//获取合约中的参数version
	version, err := storeSession.Version()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("version :", version) // "Store deployment 1.0"

	// contract write interface demo
	fmt.Println("================================")
	key := [32]byte{}
	value := [32]byte{}
	copy(key[:], []byte("foo"))
	copy(value[:], []byte("bar"))

	//插入key value
	tx, receipt, err := storeSession.SetItem(key, value)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("tx sent: %s\n", tx.Hash().Hex())
	fmt.Printf("transaction hash of receipt: %s\n", receipt.GetTransactionHash())

	// read the result
	result, err := storeSession.Items(key)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("get item: " + string(result[:])) // "bar"
}
