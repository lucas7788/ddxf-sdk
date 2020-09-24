package main

import (
	"fmt"
	"github.com/ont-bizsuite/ddxf-sdk"
	"github.com/ontio/ontology/common/password"
	"github.com/urfave/cli"
	"os"
	"runtime"
	"time"
)

var (
	AddressFlag = cli.StringFlag{
		Name:  "address",
		Usage: "address `<address base58>`",
		Value: "",
	}
	WalletFileFlag = cli.StringFlag{
		Name:  "wallet",
		Usage: "wallet `<file>`, default ./wallet.dat",
		Value: "./wallet.dat",
	}
	NetworkIdFlag = cli.UintFlag{
		Name:  "networkid",
		Usage: "networkid `<id>` (0~1). 0:testnet 1:mainnet, default 1",
		Value: uint(1),
	}
	SlotFlag = cli.UintFlag{
		Name:  "slot",
		Usage: "send tx slot `<slot>` (1~) millisecond. default 3000 millisecond",
		Value: uint(3000),
	}
)

func setupAPP() *cli.App {
	app := cli.NewApp()
	app.Usage = "Bonus CLI"
	app.Action = start
	app.Copyright = "Copyright in 2018 The Ontology Authors"
	app.Flags = []cli.Flag{
		AddressFlag,
		WalletFileFlag,
		NetworkIdFlag,
		SlotFlag,
	}
	app.Before = func(context *cli.Context) error {
		runtime.GOMAXPROCS(runtime.NumCPU())
		return nil
	}
	return app
}

func main() {
	if err := setupAPP().Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func start(ctx *cli.Context) {
	networkId := ctx.GlobalInt(NetworkIdFlag.Name)
	//默认主网
	testNet := ddxf_sdk.MainNet
	if networkId == 0 {
		testNet = ddxf_sdk.TestNet
	}
	sdk := ddxf_sdk.NewDdxfSdk(testNet)
	ontSdk := sdk.GetOntologySdk()
	walletFile := ctx.GlobalString(WalletFileFlag.Name)
	fmt.Printf("walletFile: %s\n", walletFile)
	wallet, err := ontSdk.OpenWallet(walletFile)
	if err != nil {
		fmt.Printf("error in ReadFile:%s\n", err)
		return
	}
	passwd, err := password.GetAccountPassword()
	if err != nil {
		fmt.Printf("input password error: %s\n", err)
		return
	}
	address := ctx.GlobalString(AddressFlag.Name)
	if address == "" {
		fmt.Printf("please input address encode in base58\n")
		return
	}
	acc, err := wallet.GetAccountByAddress(address, passwd)
	if err != nil {
		fmt.Printf("GetDefaultAccount error: %s\n", err)
		return
	}
	slot := ctx.GlobalInt(SlotFlag.Name)
	for {
		txhash, err := sdk.GetOntologySdk().Native.Ong.Transfer(2500, 20000, acc, acc, acc.Address, 1)
		if err != nil {
			fmt.Printf("ong transfer fail: %s\n", err)
			return
		}
		fmt.Printf("txhash: %s\n", txhash.ToHexString())
		time.Sleep(time.Duration(slot) * time.Millisecond)
	}
}
