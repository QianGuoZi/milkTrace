package dal

import (
	"fmt"
	"github.com/FISCO-BCOS/go-sdk/client"
	"github.com/FISCO-BCOS/go-sdk/conf"
	"log"
	tls "server/Tls"
)

var TlsApi *tls.TlsSession

func InitTls() {
	configs, err := conf.ParseConfigFile("config.toml")
	if err != nil {
		log.Fatalf("ParseConfigFile failed, err: %v", err)
	}

	clients, err := client.Dial(&configs[0])
	if err != nil {
		log.Fatal(err)
	}

	address, tx, instance, err := tls.DeployTls(clients.GetTransactOpts(), clients)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("contract address: ", address.Hex()) // the address should be saved, will use in next example
	fmt.Println("transaction hash: ", tx.Hash().Hex())
	TlsApi = &tls.TlsSession{Contract: instance, CallOpts: *clients.GetCallOpts(), TransactOpts: *clients.GetTransactOpts()}
}
