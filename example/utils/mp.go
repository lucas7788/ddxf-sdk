package utils

import (
	"github.com/ontio/ontology/common"
	"github.com/ont-bizsuite/ddxf-sdk"
	"github.com/ontio/ontology-go-sdk"
	"encoding/hex"
	"io/ioutil"
	"fmt"
)

//wasmFile := "/Users/sss/dev/dockerData/rust_project/ddxf_market/output/marketplace.wasm"
//wasmFile = "/Users/sss/dev/dockerData/rust_project/ddxf_market/output/dtoken.wasm"
//wasmFile = "/Users/sss/dev/dockerData/rust_project/ddxf_market/output/data_id.wasm"
//wasmFile = "/Users/sss/dev/dockerData/rust_project/ddxf_market/output/split_policy.wasm"
//wasmFile = "/Users/sss/dev/rust_project/oep4-rust/output/oep_4.wasm"
//wasmFile = "/Users/sss/dev/dockerData/rust_project/ddxf_market/output/open_kg.wasm"
//wasmFile = "/Users/sss/dev/dockerData/rust_project/ddxf_market/output/accountant.wasm"
//wasmFile = "/Users/sss/dev/dockerData/rust_project/vote/output/vote.wasm"

func DeployOep4Contract(sdk *ddxf_sdk.DdxfSdk, admin *ontology_go_sdk.Account, gasPrice uint64) {
	wasmFile := "/Users/sss/dev/rust_project/oep4-rust/output/oep_4.wasm"
	code, _ := ioutil.ReadFile(wasmFile)

	contractAddr := common.AddressFromVmCode(code)
	fmt.Printf("Oep4Addr:%s, Oep4Addr:%s\n", contractAddr.ToBase58(), contractAddr.ToHexString())

	name := "OpenKG_Token"
	desc := "OEP4 contract for OpenKG, stands for the value of a knowledge, which benefits OpenKG community."
	DeployContract(sdk, admin, hex.EncodeToString(code), name, desc, gasPrice)
}

func DeployMpContract(sdk *ddxf_sdk.DdxfSdk, admin *ontology_go_sdk.Account, gasPrice uint64) {
	wasmFile := "/Users/sss/dev/dockerData/rust_project/ddxf_market/output/marketplace.wasm"
	code, _ := ioutil.ReadFile(wasmFile)

	contractAddr := common.AddressFromVmCode(code)
	fmt.Printf("MpcontractAddr:%s, MpcontractAddr:%s\n", contractAddr.ToBase58(), contractAddr.ToHexString())
    return
	name := "DDXF - Marketplace"
	desc := "Sub-contract of DDXF series, provides a place for token/DToken owner to manage their published items, and enable token exchange between token owner and token acquiers."
	DeployContract(sdk, admin, hex.EncodeToString(code), name, desc, gasPrice)
}


func DeployDTokenContract(sdk *ddxf_sdk.DdxfSdk, admin *ontology_go_sdk.Account, gasPrice uint64) {
	wasmFile := "/Users/sss/dev/dockerData/rust_project/ddxf_market/output/dtoken.wasm"
	code, _ := ioutil.ReadFile(wasmFile)

	contractAddr := common.AddressFromVmCode(code)
	fmt.Printf("DTokenContractAddr:%s, DTokenContractAddr:%s\n", contractAddr.ToBase58(), contractAddr.ToHexString())

	name := "DDXF - DToken (base)"
	desc := "Sub-contract of DDXF series, a standard proposal to combine off-chain access-token (tokenization) with on-chain token (assertization). Provides support for data management, esp., the permission control. This contract is a basic version for DToken, which cannot be retransferred. For the rest DToken contracts, which support OEP4, OEP5, and OEP8."

	DeployContract(sdk, admin, hex.EncodeToString(code), name, desc, gasPrice)

	//migrate
	if false {
		contractAddr, _ := common.AddressFromHexString("45bc078c0664a11cce87d97864c5b3594c7c9f81")
		con := sdk.DefContract(contractAddr)
		txhash, err := con.Invoke("migrate", admin, []interface{}{code, })
		if err != nil {
			fmt.Println(err)
			return
		}
		evt, err := sdk.GetSmartCodeEvent(txhash.ToHexString())
		fmt.Println(evt, err)
		return
	}
}

func DeploySplitPolicyContract(sdk *ddxf_sdk.DdxfSdk, admin *ontology_go_sdk.Account, gasPrice uint64) {
	wasmFile := "/Users/sss/dev/dockerData/rust_project/ddxf_market/output/split_policy.wasm"
	code, _ := ioutil.ReadFile(wasmFile)

	contractAddr := common.AddressFromVmCode(code)
	fmt.Printf("SplitPolicyContractAddr:%s, SplitPolicyContractAddr:%s\n", contractAddr.ToBase58(), contractAddr.ToHexString())

	name := "DDXF - SplitPolicy"
	desc := "Sub-contract of DDXF series, provides fee split services for marketplace."

	DeployContract(sdk, admin, hex.EncodeToString(code), name, desc, gasPrice)
}

func DeployDataIdContract(sdk *ddxf_sdk.DdxfSdk, admin *ontology_go_sdk.Account, gasPrice uint64) {
	wasmFile := "/Users/sss/dev/dockerData/rust_project/ddxf_market/output/data_id.wasm"
	code, _ := ioutil.ReadFile(wasmFile)

	contractAddr := common.AddressFromVmCode(code)
	fmt.Printf("DataIdContractAddr:%s, DataIdContractAddr:%s\n", contractAddr.ToBase58(), contractAddr.ToHexString())

	name := "DDXF - DataID"
	desc := "Sub-contract of DDXF series, provides batch function to register data identifier with attributes by ONT ID, and save transaction cost."

	DeployContract(sdk, admin, hex.EncodeToString(code), name, desc, gasPrice)
}

func DeployAccountantContract(sdk *ddxf_sdk.DdxfSdk, admin *ontology_go_sdk.Account, gasPrice uint64) {
	wasmFile := "/Users/sss/dev/dockerData/rust_project/ddxf_market/output/accountant.wasm"
	code, _ := ioutil.ReadFile(wasmFile)

	contractAddr := common.AddressFromVmCode(code)
	fmt.Printf("AccountantContractAddr:%s, AccountantContractAddr:%s\n", contractAddr.ToBase58(), contractAddr.ToHexString())

	name := "DDXF - Accountant"
	desc := "Sub-contract of DDXF series, provides fee split services between the marketplace and the token owner(s)."

	DeployContract(sdk, admin, hex.EncodeToString(code), name, desc, gasPrice)
}

func DeployOpenkgContract(sdk *ddxf_sdk.DdxfSdk, admin *ontology_go_sdk.Account, gasPrice uint64) {
	wasmFile := "/Users/sss/dev/dockerData/rust_project/ddxf_market/output/open_kg.wasm"
	code, _ := ioutil.ReadFile(wasmFile)

	contractAddr := common.AddressFromVmCode(code)
	fmt.Printf("OpenkgContractAddr:%s, OpenkgContractAddr:%s\n", contractAddr.ToBase58(), contractAddr.ToHexString())

	name := "OpenKG"
	desc := "Combination smart contract for OpenKG, which makes full use of Ontology DDXF and ONT ID. OpenKG (http://openkg.cn/) is an open knowledge graph project sponsor by Chinese Information Processing Society of China (http://www.cipsc.org.cn/sigkg/). The project is targeting to improve the openness and interconnection of knowledge graph in Chinese domain, to encourage the degree of openness of knowledge graph, its algorithms, models and tools, and to promote the implementation of knowledge graph and semantic technology."

	DeployContract(sdk, admin, hex.EncodeToString(code), name, desc, gasPrice)
}