package main

import (
	"fmt"
	"github.com/ont-bizsuite/ddxf-sdk"
	"github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology/common/password"
	"github.com/urfave/cli"
	"os"
	"runtime"
	"sync"
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
		Value: uint(2500),
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

var mainNet = []string{
	"http://dappnode1.ont.io:20336",
	"http://dappnode2.ont.io:20336",
	"http://dappnode3.ont.io:20336"}

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
	index := 0
	oldHeight := getCurBlockHeight(ontSdk, index)
	var curHeight uint32

	ticker := time.NewTimer(time.Duration(slot) * time.Millisecond)

	wait := new(sync.WaitGroup)

	wait.Add(1)
	go func() {
		defer wait.Done()
		for {
			time.Sleep(1 * time.Second)
			curHeight = getCurBlockHeight(ontSdk, index)
			fmt.Printf("oldHeight:%d, curHeight:%d\n", oldHeight, curHeight)
			if curHeight == oldHeight {
				continue
			}
			oldHeight = curHeight
			if !ticker.Stop() {
				<-ticker.C
			}
			ticker.Reset(time.Duration(slot) * time.Millisecond)
		}
	}()

	wait.Add(1)
	go func() {
		defer wait.Done()
		for {
			select {
			case <-ticker.C:
				txhash, err := ontSdk.Native.Ong.Transfer(2500, 20000, acc, acc, acc.Address, 1)
				if err != nil {
					fmt.Printf("ong transfer fail: %s\n", err)
					ind := (index + 1) % len(mainNet)
					ontSdk.NewRpcClient().SetAddress(mainNet[ind])
				}
				fmt.Printf("txhash: %s\n", txhash.ToHexString())
			}
		}
	}()

	wait.Wait()

	//for {
	//	time.Sleep(time.Duration(slot) * time.Millisecond)
	//	curHeight = getCurBlockHeight(ontSdk, index)
	//	if curHeight != oldHeight {
	//		oldHeight = curHeight
	//		continue
	//	}
	//	txhash, err := ontSdk.Native.Ong.Transfer(2500, 20000, acc, acc, acc.Address, 1)
	//	if err != nil {
	//		fmt.Printf("ong transfer fail: %s\n", err)
	//		ind := (index + 1) % len(mainNet)
	//		ontSdk.NewRpcClient().SetAddress(mainNet[ind])
	//		continue
	//	}
	//	oldHeight = curHeight
	//	fmt.Printf("curBlockHeight: %d, txhash: %s\n", curHeight, txhash.ToHexString())
	//}
}

func getCurBlockHeight(ontSdk *ontology_go_sdk.OntologySdk, index int) uint32 {
	curHeight, err := ontSdk.GetCurrentBlockHeight()
	if err != nil {
		ind := (index + 1) % len(mainNet)
		ontSdk.NewRpcClient().SetAddress(mainNet[ind])
		return getCurBlockHeight(ontSdk, ind)
	}
	return curHeight
}
