package utils

import (
	"github.com/ont-bizsuite/ddxf-sdk/any_contract"
	"github.com/ontio/ontology-go-sdk"
	"github.com/ont-bizsuite/ddxf-sdk"
	"fmt"
	"github.com/ontio/ontology/common"
)

func Init(sdk *ddxf_sdk.DdxfSdk,con *any_contract.ContractKit,admin *ontology_go_sdk.Account,mp,dtoken common.Address) {
	txHash, err := con.Invoke("init",admin, []interface{}{mp,dtoken})
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	showNotify(sdk,"init", txHash.ToHexString())
}

func SetDtokenContractAddr(sdk *ddxf_sdk.DdxfSdk,con *any_contract.ContractKit,admin *ontology_go_sdk.Account,dtoken common.Address) {
	txHash, err := con.Invoke("setDtokenContractAddr",admin, []interface{}{dtoken})
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	showNotify(sdk,"setDtokenContractAddr", txHash.ToHexString())
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

func SetMpContractAddr(sdk *ddxf_sdk.DdxfSdk,con *any_contract.ContractKit,admin *ontology_go_sdk.Account,mp common.Address) {
	txHash, err := con.Invoke("setMpContractAddr",admin, []interface{}{mp})
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	showNotify(sdk,"SetMpContractAddr", txHash.ToHexString())
}

func BuyAndUseToken(sdk *ddxf_sdk.DdxfSdk,con *any_contract.ContractKit,
	resource_id []byte, n int, buyer_account *ontology_go_sdk.Account, payer *ontology_go_sdk.Account,
	token_template_id []byte) {
	tx, err := con.BuildTx("buyAndUseToken",
		[]interface{}{resource_id, n,buyer_account.Address,payer.Address, token_template_id})
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	sdk.SignTx(tx, buyer_account)
	sdk.SignTx(tx, payer)
	txHash, err := sdk.SendTx(tx)
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	showNotify(sdk,"BuyAndUseToken", txHash.ToHexString())
}