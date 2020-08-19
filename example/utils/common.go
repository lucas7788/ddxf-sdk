package utils

import (
	"github.com/ont-bizsuite/ddxf-sdk"
	"github.com/ontio/ontology-go-sdk"
	"fmt"
)

func DeployContract(sdk *ddxf_sdk.DdxfSdk, admin *ontology_go_sdk.Account, codeHex, name, desc string, gasPrice uint64) {
	sdk.SetGasPrice(gasPrice)
	if false {
		name = "ontology-vote"
		desc = "smart contract for ontology vote"
	}
	if false {
		name = "dataid-batch"
		desc = "smart contract for ontology dataid-batch-action"
	}

	txHash, err := sdk.DeployContract(admin, codeHex, name, "0.1.1", "lucas", "", desc)
	if err != nil {
		fmt.Printf("DeployContract error:%s\n", err)
		return
	}
	evt, err := sdk.GetSmartCodeEvent(txHash.ToHexString())
	if err != nil {
		fmt.Printf("DeployContract GetSmartCodeEvent error:%s, txHash: %s\n", err, txHash.ToHexString())
		return
	}
	fmt.Println("evt:", evt)
}

