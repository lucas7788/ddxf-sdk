package utils

import (
	"fmt"
	"github.com/ont-bizsuite/ddxf-sdk"
	"github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/smartcontract/states"
)

type Token struct {
	TokenName       string
	TokenType       uint8
	ContractAddress common.Address
}

func NewToken(TokenName string,
	TokenType uint8,
	ContractAddress common.Address) Token {
	return Token{
		TokenName:       TokenName,
		TokenType:       TokenType,
		ContractAddress: ContractAddress,
	}
}

func (token *Token) Serialize(sink *common.ZeroCopySink) {
	sink.WriteString(token.TokenName)
	sink.WriteByte(token.TokenType)
	sink.WriteAddress(token.ContractAddress)
}

func BuildWasmParam(contractAddress common.Address, argbytes []byte) {
	contract := &states.WasmContractParam{}
	contract.Address = contractAddress
	contract.Args = argbytes
}

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

	txHash, err := sdk.DeployContract(admin, codeHex, name, "1.0", "Wing Team", "support@wing.finance", desc)
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
