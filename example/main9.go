package main

import (
	"encoding/json"
	"fmt"
	"github.com/ont-bizsuite/ddxf-sdk"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/common/password"
	"github.com/urfave/cli"
	"io/ioutil"
	"os"
	"runtime"
)

var (
	NetworkIdFlag = cli.UintFlag{
		Name:  "networkid",
		Usage: "networkid `<id>` (0~1). 0:testnet 1:mainnet, default 0",
		Value: uint(0),
	}
	FTokenFlag = cli.BoolFlag{
		Name:  "ftoken",
		Usage: "--ftoken.",
	}
	ComptrollerFlag = cli.BoolFlag{
		Name:  "comptroller",
		Usage: "--comptroller",
	}
)

func setupAPP() *cli.App {
	app := cli.NewApp()
	app.Usage = "Bonus CLI"
	app.Action = start
	app.Copyright = "Copyright in 2018 The Ontology Authors"
	app.Flags = []cli.Flag{
		NetworkIdFlag,
		FTokenFlag,
		ComptrollerFlag,
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
	testNet := ddxf_sdk.TestNet
	if networkId == 1 {
		testNet = ddxf_sdk.MainNet
	}
	sdk := ddxf_sdk.NewDdxfSdk(testNet)
	ontSdk := sdk.GetOntologySdk()

	var configFile string
	if networkId == 1 {
		configFile = "./config.json"
	} else {
		configFile = "./config-test.json"
	}
	configBs, err := ioutil.ReadFile(configFile)
	if err != nil {
		fmt.Println("read config.json failed", err)
		return
	}
	configMap := make(map[string]interface{})

	err = json.Unmarshal(configBs, &configMap)
	if err != nil {
		fmt.Println("read config.json failed", err)
		return
	}

	walletFile := configMap["wallet"].(string)
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

	adminAddress := configMap["admin"].(string)

	admin, err := wallet.GetAccountByAddress(adminAddress, passwd)
	if err != nil {
		fmt.Printf("GetDefaultAccount error: %s\n", err)
		return
	}

	wingUtilsAddres := configMap["wing_utils_address"].(string)
	wingUtilsAddress, err := common.AddressFromHexString(wingUtilsAddres)
	if err != nil {
		fmt.Println("parse address failed", err)
		return
	}
	wingUtils := sdk.DefContract(wingUtilsAddress)
	ftoken := ctx.GlobalBool(FTokenFlag.Name)
	comptroller := ctx.GlobalBool(ComptrollerFlag.Name)
	if ftoken {
		ftoken_address := configMap["ftoken_address"].(string)
		ftoken, err := common.AddressFromHexString(ftoken_address)
		if err != nil {
			fmt.Println("parse ftoken address failed: ", err)
			return
		}
		txhash, err := wingUtils.Invoke("set_flash_pool", admin, []interface{}{ftoken})
		if err != nil {
			fmt.Println("invoke contract failed: ", err)
			return
		}
		showNotify(sdk, "set_flash_pool", txhash.ToHexString())
	}
	if comptroller {
		comptroller_address := configMap["comptroller_address"].(string)
		comptroller, err := common.AddressFromHexString(comptroller_address)
		if err != nil {
			fmt.Println("parse ftoken address failed: ", err)
			return
		}
		txhash, err := wingUtils.Invoke("set_comptroller", admin, []interface{}{comptroller})
		if err != nil {
			fmt.Println("invoke contract failed: ", err)
			return
		}
		showNotify(sdk, "set_comptroller", txhash.ToHexString())
	}
	fmt.Println("over")
}

func showNotify(sdk *ddxf_sdk.DdxfSdk, method, txHash string) error {
	fmt.Printf("method: %s, txHash: %s\n", method, txHash)
	evt, err := sdk.GetSmartCodeEvent(txHash)
	if err != nil {
		return err
	}
	for _, notify := range evt.Notify {
		fmt.Printf("method: %s,evt: %v\n", method, notify)
	}
	return nil
}
