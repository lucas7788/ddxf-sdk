package main

import (
	"fmt"
	"github.com/ont-bizsuite/ddxf-sdk"
	"github.com/ont-bizsuite/ddxf-sdk/any_contract"
	"github.com/ont-bizsuite/ddxf-sdk/example/utils"
	"github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology/common"
)

var (
	admin *ontology_go_sdk.Account
)

func main() {
	gasPrice := uint64(2500)
	testNet := "http://106.75.224.136:20336"
	testNet = ddxf_sdk.TestNet
	//testNet = "http://172.168.3.152:20336"
	//testNet = "http://127.0.0.1:20336"

	//testNet = ddxf_sdk.MainNet
	sdk := ddxf_sdk.NewDdxfSdk(testNet)
	//106.75.224.136

	pwd := []byte("123456")
	ontSdk := sdk.GetOntologySdk()
	wallet, err := ontSdk.OpenWallet("./wallet.dat")
	if err != nil {
		fmt.Printf("error in ReadFile:%s\n", err)
		return
	}
	//Aejfo7ZX5PVpenRj23yChnyH64nf8T1zbu
	admin, _ = wallet.GetAccountByAddress("Aejfo7ZX5PVpenRj23yChnyH64nf8T1zbu", pwd)

	//pri, _ := keypair.WIF2Key([]byte("KySMiNrDDzFyUxfpK2hV9wFivq6hEmgB81D1UynhwjXjgd7xUZ88"))
	//pub := pri.Public()
	//add := types.AddressFromPubKey(pub)
	//admin = &ontology_go_sdk.Account{
	//	PrivateKey: pri,
	//	PublicKey:  pub,
	//	Address:    add,
	//}

	//dc528ea5bab011f6a069880c37ae1ae499dfb58a
	if false {
		utils.DeployWingUtilsContract(sdk, admin, gasPrice)
		return
	}
	//d3878cbec82d763d9fd72ec1bff0d22ae703f2d7
	wingUtilsContract, _ := common.AddressFromHexString("bbb0cd417e905ddb9a84f84775f9b496ea6344a6")
	zeroContract, _ := common.AddressFromHexString("47414f5d38fef4ea199502c7423716d208468a2c")
	flashContract, _ := common.AddressFromHexString("fcd78ad03e37b21f9d93796857ca7b0e3f4bbfae")
	gov, _ := common.AddressFromHexString("7b5553509adb056d6422827442bc9d9a675abf75")
	comptroller, _ := common.AddressFromHexString("5a86c02cd0a1fd55fc9d3b58954a3d041648b96b")
	ontdContract, _ := common.AddressFromHexString("304454f9166d78901137e203116f9d759d075ece")
	wingUtils := sdk.DefContract(wingUtilsContract)
	if false {
		txhash, err := wingUtils.Invoke("init", admin, []interface{}{zeroContract, flashContract, gov, comptroller, ontdContract})
		if err != nil {
			fmt.Println(err)
			return
		}
		showNotify(sdk, "init", txhash.ToHexString())
		return
	}
	if false {
		txhash, err := wingUtils.Invoke("set_zero", admin, []interface{}{zeroContract})
		if err != nil {
			fmt.Println(err)
			return
		}
		showNotify(sdk, "set_zero", txhash.ToHexString())
		return
	}
	if false {
		txhash, err := sdk.GetOntologySdk().NeoVM.InvokeNeoVMContract(gasPrice, 20000, admin, admin, ontdContract, []interface{}{"ont2ontd", []interface{}{admin.Address, 1}})
		if err != nil {
			fmt.Println(err)
			return
		}
		showNotify(sdk, "ont2ontd", txhash.ToHexString())
		return
	}
	if true {
		txhash, err := wingUtils.Invoke("swap_to_flash_pool", admin, []interface{}{admin.Address})
		if err != nil {
			fmt.Println(err)
			return
		}
		showNotify(sdk, "swap_to_flash_pool", txhash.ToHexString())
		return
	}
	if false {
		getContractAddress(wingUtils, "get_zero")
		getContractAddress(wingUtils, "get_flash_pool")
		getContractAddress(wingUtils, "get_ontd")
		return
	}

	if false {
		data, err := wingUtils.PreInvoke("unbound_wing", []interface{}{admin.Address})
		if err != nil {
			fmt.Println(err)
			return
		}
		bs, _ := data.ToByteArray()
		source := common.NewZeroCopySource(bs)
		res, _ := source.NextI128()
		fmt.Println(res.ToNumString())
		return
	}

}

func getContractAddress(wingUtils *any_contract.ContractKit, method string) {
	res, err := wingUtils.PreInvoke(method, []interface{}{})
	if err != nil {
		fmt.Println(err)
		return
	}
	bs, err := res.ToByteArray()
	if err != nil {
		fmt.Println(err)
		return
	}
	source := common.NewZeroCopySource(bs)
	addr, eof := source.NextAddress()
	if eof {
		fmt.Println("eof is true")
		return
	}
	fmt.Printf("method: %s, addr:%s\n", method, addr.ToHexString())
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
