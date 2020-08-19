package dtoken_contract

import (
	"testing"

	"encoding/hex"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/ont-bizsuite/ddxf-sdk/base_contract"
	ontology_go_sdk "github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology/common"
	"github.com/stretchr/testify/assert"
)

var (
	dTokenKit *DTokenKit
	ontSdk    *ontology_go_sdk.OntologySdk
	wallet    *ontology_go_sdk.Wallet
	admin     *ontology_go_sdk.Account

	testNet   = "http://polaris1.ont.io:20336"
	localHost = "http://127.0.0.1:20336"
	pwd       = []byte("123456")
	gasPrice  = uint64(0)
	gasLimit  = uint64(28400000)
)

func TestMain(m *testing.M) {
	ontSdk = ontology_go_sdk.NewOntologySdk()
	ontSdk.NewRpcClient().SetAddress(localHost)
	var err error
	wallet, err = ontSdk.OpenWallet("../wallet.dat")
	if err != nil {
		fmt.Printf("error in ReadFile:%s\n", err)
		return
	}
	admin, _ = wallet.GetAccountByAddress("AYnhakv7kC9R5ppw65JoE2rt6xDzCjCTvD", pwd)
	wasmFile := "/Users/sss/dev/dockerData/rust_project/ddxf_market/output/dtoken.wasm"
	code, err := ioutil.ReadFile(wasmFile)
	if err != nil {
		fmt.Printf("error in ReadFile:%s\n", err)
		return
	}
	addr := common.AddressFromVmCode(code)
	fmt.Printf("contract address: %s\n", addr.ToHexString())
	//only need execute once
	if false {
		txHash, err := ontSdk.WasmVM.DeployWasmVMSmartContract(gasPrice, gasLimit, admin,
			hex.EncodeToString(code), "", "", "", "", "")
		if err != nil {
			fmt.Println("err:", err)
			return
		}
		time.Sleep(10 * time.Second)
		evt, err := ontSdk.GetSmartContractEvent(txHash.ToHexString())
		if err != nil {
			fmt.Println("err:", err)
			return
		}
		fmt.Println("evts: ", evt)
	}
	contractAddress := common.AddressFromVmCode(code)
	bc := base_contract.NewBaseContract(ontSdk, 20000000, gasPrice, admin)
	dTokenKit = NewDTokenKit(contractAddress, bc)
	m.Run()
}

func TestDTokenKit_GetDDXFContractAddr(t *testing.T) {
	addr, err := dTokenKit.GetMpContractAddr(common.ADDRESS_EMPTY)
	assert.Nil(t, err)
	fmt.Println("addr: ", addr.ToHexString())
}

//49c2dc97ee58b2292e55499e1122c579fc0690e3
func TestDTokenKit_SetDDXFContractAddr(t *testing.T) {
	addr, _ := common.AddressFromHexString("f0020843718912d5f25977ffd8ea7e4eb00601a1")
	txHash, err := dTokenKit.SetMpContractAddr(common.ADDRESS_EMPTY, admin, addr)
	assert.Nil(t, err)

	time.Sleep(10 * time.Second)
	evt, err := ontSdk.GetSmartContractEvent(txHash.ToHexString())
	assert.Nil(t, err)
	fmt.Println("evt:", evt)
}
