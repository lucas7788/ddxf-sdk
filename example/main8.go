package main

import (
	"fmt"
	"github.com/ont-bizsuite/ddxf-sdk"
	"github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology/common"
)

var (
	admin *ontology_go_sdk.Account
)

func main() {
	gasPrice := uint64(2500)
	testNet := "http://172.168.3.226:20336"
	testNet = ddxf_sdk.TestNet

	sdk := ddxf_sdk.NewDdxfSdk(testNet)
	sdk.SetGasPrice(gasPrice)

	pwd := []byte("123456")
	ontSdk := sdk.GetOntologySdk()
	wallet, err := ontSdk.OpenWallet("./wallet.dat")
	if err != nil {
		fmt.Printf("error in ReadFile:%s\n", err)
		return
	}

	admin, _ = wallet.GetAccountByAddress("Aejfo7ZX5PVpenRj23yChnyH64nf8T1zbu", pwd)

	wingUtilsContract, _ := common.AddressFromHexString("403de1ed53be768960599e0996d487e55d0376f1")

	fToken, _ := common.AddressFromHexString("414ec1faffae6a6f9c55e7bc7e5aa91b3a4a0b5a")
	comptroller, _ := common.AddressFromHexString("4d05ce669024b14a2e4594c077435756fb1d9606")

	wingUtils := sdk.DefContract(wingUtilsContract)

	if false {
		txhash, err := wingUtils.Invoke("set_flash_pool", admin, []interface{}{fToken})
		checkerr(err)
		showNotify(sdk, "set_flash_pool", txhash.ToHexString())
		return
	}
	if false {
		txhash, err := wingUtils.Invoke("set_comptroller", admin, []interface{}{comptroller})
		checkerr(err)
		showNotify(sdk, "set_comptroller", txhash.ToHexString())
		return
	}
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

func checkerr(err error) {
	if err != nil {
		panic(err)
	}
}
