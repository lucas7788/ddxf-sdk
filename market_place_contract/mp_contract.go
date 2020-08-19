package market_place_contract

import (
	"github.com/ont-bizsuite/ddxf-sdk/base_contract"
	"github.com/ont-bizsuite/ddxf-sdk/split_policy_contract"
	ontology_go_sdk "github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/core/types"
)

type MpKit struct {
	bc              *base_contract.BaseContract
	contractAddress common.Address
}

func NewDDXFContractKit(contractAddress common.Address, bc *base_contract.BaseContract) *MpKit {
	return &MpKit{
		contractAddress: contractAddress,
		bc:              bc,
	}
}

func (this *MpKit) SetContractAddress(addr common.Address) {
	this.contractAddress = addr
}

func (this *MpKit) Init(contractAddress common.Address, admin *ontology_go_sdk.Account, dtoken, splitPolicy common.Address) (common.Uint256, error) {
	if contractAddress == common.ADDRESS_EMPTY {
		contractAddress = this.contractAddress
	}
	return this.bc.Invoke(contractAddress, admin, "init",
		[]interface{}{dtoken, splitPolicy})
}

func (this *MpKit) setDTokenContractAddress(contractAddress common.Address, admin *ontology_go_sdk.Account, dtoken common.Address) (common.Uint256, error) {
	if contractAddress == common.ADDRESS_EMPTY {
		contractAddress = this.contractAddress
	}
	return this.bc.Invoke(contractAddress, admin, "setDTokenContract", []interface{}{dtoken})
}

func (this *MpKit) getDTokenContractAddress(contractAddress common.Address) (common.Address, error) {
	if contractAddress == common.ADDRESS_EMPTY {
		contractAddress = this.contractAddress
	}
	res, err := this.bc.PreInvoke(contractAddress, "getDTokenContract", []interface{}{})
	if err != nil {
		return common.ADDRESS_EMPTY, err
	}
	data, err := res.ToByteArray()
	if err != nil {
		return common.ADDRESS_EMPTY, err
	}
	return common.AddressParseFromBytes(data)
}

func (this *MpKit) setSplitPolicyContractAddress(contractAddress common.Address, admin *ontology_go_sdk.Account, dtoken common.Address) (common.Uint256, error) {
	if contractAddress == common.ADDRESS_EMPTY {
		contractAddress = this.contractAddress
	}
	return this.bc.Invoke(contractAddress, admin, "setSplitPolicyContract", []interface{}{dtoken})
}

func (this *MpKit) getSplitPolicyContractAddress(contractAddress common.Address) (common.Address, error) {
	if contractAddress == common.ADDRESS_EMPTY {
		contractAddress = this.contractAddress
	}
	res, err := this.bc.PreInvoke(contractAddress, "getSplitPolicyContract", []interface{}{})
	if err != nil {
		return common.ADDRESS_EMPTY, err
	}
	data, err := res.ToByteArray()
	if err != nil {
		return common.ADDRESS_EMPTY, err
	}
	return common.AddressParseFromBytes(data)
}

//publish product on block chain,
func (this *MpKit) Publish(contractAddress common.Address, seller *ontology_go_sdk.Account, resourceId []byte, ddo ResourceDDO, item DTokenItem,
	splitPolicyParam split_policy_contract.SplitPolicyRegisterParam) (common.Uint256, error) {
	if contractAddress == common.ADDRESS_EMPTY {
		contractAddress = this.contractAddress
	}
	return this.bc.Invoke(contractAddress, seller, "dtokenSellerPublish",
		[]interface{}{resourceId, ddo.ToBytes(), item.ToBytes(), splitPolicyParam.ToBytes()})
}

func (this *MpKit) Update(contractAddress common.Address, seller *ontology_go_sdk.Account, resourceId []byte, ddo ResourceDDO, item DTokenItem,
	splitPolicyParam split_policy_contract.SplitPolicyRegisterParam) (common.Uint256, error) {
	if contractAddress == common.ADDRESS_EMPTY {
		contractAddress = this.contractAddress
	}
	return this.bc.Invoke(contractAddress, seller, "update",
		[]interface{}{resourceId, ddo.ToBytes(), item.ToBytes(), splitPolicyParam.ToBytes()})
}

func (this *MpKit) BuildUpdateTx(contractAddress common.Address, resourceId []byte, ddo ResourceDDO, item DTokenItem,
	splitPolicyParam split_policy_contract.SplitPolicyRegisterParam) (*types.MutableTransaction, error) {
	if contractAddress == common.ADDRESS_EMPTY {
		contractAddress = this.contractAddress
	}
	return this.bc.BuildTx(contractAddress, "update",
		[]interface{}{resourceId, ddo.ToBytes(), item.ToBytes(), splitPolicyParam.ToBytes()})
}

func (this *MpKit) Delete(contractAddress common.Address, seller *ontology_go_sdk.Account, resourceId []byte) (common.Uint256, error) {
	if contractAddress == common.ADDRESS_EMPTY {
		contractAddress = this.contractAddress
	}
	return this.bc.Invoke(contractAddress, seller, "delete",
		[]interface{}{resourceId})
}

func (this *MpKit) BuildDeleteTx(contractAddress common.Address, resourceId []byte) (*types.MutableTransaction, error) {
	if contractAddress == common.ADDRESS_EMPTY {
		contractAddress = this.contractAddress
	}
	return this.bc.BuildTx(contractAddress, "delete",
		[]interface{}{resourceId})
}

func (this *MpKit) FreezeAndPublish(contractAddress common.Address, seller *ontology_go_sdk.Account, resourceIdOld, resourceIdNew []byte,
	ddo ResourceDDO, item DTokenItem,
	splitPolicyParam split_policy_contract.SplitPolicyRegisterParam) (common.Uint256, error) {
	if contractAddress == common.ADDRESS_EMPTY {
		contractAddress = this.contractAddress
	}
	return this.bc.Invoke(contractAddress, seller, "freezeAndPublish",
		[]interface{}{resourceIdOld, resourceIdNew, ddo.ToBytes(), item.ToBytes(), splitPolicyParam.ToBytes()})
}

func (this *MpKit) BuildFreezeAndPublishTx(contractAddress common.Address, resourceIdOld, resourceIdNew []byte,
	ddo ResourceDDO, item DTokenItem,
	splitPolicyParam split_policy_contract.SplitPolicyRegisterParam) (*types.MutableTransaction, error) {
	if contractAddress == common.ADDRESS_EMPTY {
		contractAddress = this.contractAddress
	}
	return this.bc.BuildTx(contractAddress, "freezeAndPublish", []interface{}{
		resourceIdOld, resourceIdNew, ddo.ToBytes(), item.ToBytes(), splitPolicyParam.ToBytes(),
	})
}

func (this *MpKit) BuildPublishTx(contractAddress common.Address, resourceId []byte,
	ddo ResourceDDO, item DTokenItem,
	splitPolicyParam split_policy_contract.SplitPolicyRegisterParam) (*types.MutableTransaction, error) {
	if contractAddress == common.ADDRESS_EMPTY {
		contractAddress = this.contractAddress
	}
	tx, err := this.bc.BuildTx(contractAddress, "dtokenSellerPublish",
		[]interface{}{resourceId, ddo.ToBytes(), item.ToBytes(), splitPolicyParam.ToBytes()})
	return tx, err
}

func (this *MpKit) GetPublishProductInfo(contractAddress common.Address, resourceId []byte) (*ProductInfoOnChain, error) {
	if contractAddress == common.ADDRESS_EMPTY {
		contractAddress = this.contractAddress
	}
	res, err := this.bc.PreInvoke(contractAddress, "getSellerItemInfo",
		[]interface{}{resourceId})
	if err != nil {
		return nil, err
	}
	data, err := res.ToByteArray()
	if err != nil {
		return nil, err
	}
	p := &ProductInfoOnChain{}
	err = p.FromBytes(data)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (this *MpKit) Freeze(contractAddress common.Address, manager *ontology_go_sdk.Account,
	resourceId []byte) (common.Uint256, error) {
	if contractAddress == common.ADDRESS_EMPTY {
		contractAddress = this.contractAddress
	}
	tx, err := this.bc.BuildTx(contractAddress, "freeze",
		[]interface{}{resourceId})
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	err = this.bc.GetOntologySdk().SignToTransaction(tx, manager)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	return this.bc.GetOntologySdk().SendTransaction(tx)
}

func (this *MpKit) BuildFreezeTx(contractAddress common.Address, resourceId []byte) (*types.MutableTransaction, error) {
	if contractAddress == common.ADDRESS_EMPTY {
		contractAddress = this.contractAddress
	}
	tx, err := this.bc.BuildTx(contractAddress, "freeze",
		[]interface{}{resourceId})
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (this *MpKit) BuyDToken(contractAddress common.Address, buyer, payer *ontology_go_sdk.Account, resourceId []byte,
	n int) (common.Uint256, error) {
	if payer == nil {
		payer = buyer
	}
	if contractAddress == common.ADDRESS_EMPTY {
		contractAddress = this.contractAddress
	}
	tx, err := this.bc.BuildTx(contractAddress, "buyDToken",
		[]interface{}{resourceId, n, buyer.Address, payer.Address})
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	err = this.bc.GetOntologySdk().SignToTransaction(tx, buyer)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	if buyer.Address != payer.Address {
		err = this.bc.GetOntologySdk().SignToTransaction(tx, payer)
		if err != nil {
			return common.UINT256_EMPTY, err
		}
	}
	return this.bc.GetOntologySdk().SendTransaction(tx)
}

func (this *MpKit) BuildBuyDTokenTx(contractAddress common.Address, buyer, payer common.Address, resourceId []byte,
	n int) (*types.MutableTransaction, error) {
	if contractAddress == common.ADDRESS_EMPTY {
		contractAddress = this.contractAddress
	}
	return this.bc.BuildTx(contractAddress, "buyDToken",
		[]interface{}{resourceId, n, buyer, payer})
}

func (this *MpKit) BuyDtokenReward(contractAddress common.Address, buyer, payer *ontology_go_sdk.Account, resourceId []byte,
	n int, unitPrice int) (common.Uint256, error) {
	if payer == nil {
		payer = buyer
	}
	if contractAddress == common.ADDRESS_EMPTY {
		contractAddress = this.contractAddress
	}
	tx, err := this.bc.BuildTx(contractAddress, "buyDToken",
		[]interface{}{resourceId, n, buyer.Address, payer.Address, unitPrice})
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	if buyer.Address != payer.Address {
		err = this.bc.GetOntologySdk().SignToTransaction(tx, payer)
		if err != nil {
			return common.UINT256_EMPTY, err
		}
	}
	return this.bc.GetOntologySdk().SendTransaction(tx)
}

//func (this *MpKit) BuyAndUseToken(contractAddress common.Address, buyer, payer *ontology_go_sdk.Account, resourceId []byte,
//	n int, tokenTemplate TokenTemplate) (common.Uint256, error) {
//	if payer == nil {
//		payer = buyer
//	}
// if contractAddress == common.ADDRESS_EMPTY {
// 	contractAddress = this.contractAddress
// }
//	tx, err := this.bc.BuildTx(contractAddress, "buyAndUseToken",
//		[]interface{}{resourceId, n, buyer.Address, payer.Address, tokenTemplate.ToBytes()})
//	if err != nil {
//		return common.UINT256_EMPTY, err
//	}
//	if buyer.Address != payer.Address {
//		err = this.bc.GetOntologySdk().SignToTransaction(tx, payer)
//		if err != nil {
//			return common.UINT256_EMPTY, err
//		}
//	}
//	return this.bc.GetOntologySdk().SendTransaction(tx)
//}

func (this *MpKit) BuildBuyAndUseTokenTx(contractAddress common.Address, buyer, payer common.Address, resourceId []byte,
	n int, tokenTemplate TokenTemplate) (*types.MutableTransaction, error) {
	if contractAddress == common.ADDRESS_EMPTY {
		contractAddress = this.contractAddress
	}
	tx, err := this.bc.BuildTx(contractAddress, "buyAndUseToken",
		[]interface{}{resourceId, n, buyer, payer, tokenTemplate.ToBytes()})
	if err != nil {
		return nil, err
	}
	return tx, nil
}
