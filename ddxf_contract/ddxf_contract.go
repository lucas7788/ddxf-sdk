package ddxf_contract

import (
	"github.com/ont-bizsuite/ddxf-sdk/base_contract"
	"github.com/ont-bizsuite/ddxf-sdk/split_policy_contract"
	"github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/core/types"
)

type DDXFKit struct {
	bc              *base_contract.BaseContract
	contractAddress common.Address
}

func NewDDXFContractKit(contractAddress common.Address, bc *base_contract.BaseContract) *DDXFKit {
	return &DDXFKit{
		contractAddress: contractAddress,
		bc:              bc,
	}
}

func (this *DDXFKit) SetContractAddress(addr common.Address) {
	this.contractAddress = addr
}

func (this *DDXFKit) Init(admin *ontology_go_sdk.Account, dtoken, splitPolicy common.Address) (common.Uint256, error) {
	return this.bc.Invoke(this.contractAddress, admin, "init",
		[]interface{}{dtoken, splitPolicy})
}

func (this *DDXFKit) setDTokenContractAddress(admin *ontology_go_sdk.Account, dtoken common.Address) (common.Uint256, error) {
	return this.bc.Invoke(this.contractAddress, admin, "setDTokenContract", []interface{}{dtoken})
}

func (this *DDXFKit) getDTokenContractAddress() (common.Address, error) {
	res, err := this.bc.PreInvoke(this.contractAddress, "getDTokenContract", []interface{}{})
	if err != nil {
		return common.ADDRESS_EMPTY, err
	}
	data, err := res.ToByteArray()
	if err != nil {
		return common.ADDRESS_EMPTY, err
	}
	return common.AddressParseFromBytes(data)
}

func (this *DDXFKit) setSplitPolicyContractAddress(admin *ontology_go_sdk.Account, dtoken common.Address) (common.Uint256, error) {
	return this.bc.Invoke(this.contractAddress, admin, "setSplitPolicyContract", []interface{}{dtoken})
}

func (this *DDXFKit) getSplitPolicyContractAddress() (common.Address, error) {
	res, err := this.bc.PreInvoke(this.contractAddress, "getSplitPolicyContract", []interface{}{})
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
func (this *DDXFKit) Publish(seller *ontology_go_sdk.Account, resourceId []byte, ddo ResourceDDO, item DTokenItem,
	splitPolicyParam split_policy_contract.SplitPolicyRegisterParam) (common.Uint256, error) {
	return this.bc.Invoke(this.contractAddress, seller, "dtokenSellerPublish",
		[]interface{}{resourceId, ddo.ToBytes(), item.ToBytes(), splitPolicyParam.ToBytes()})
}

func (this *DDXFKit) BuildPublishTx(resourceId []byte,
	ddo ResourceDDO, item DTokenItem,
	splitPolicyParam split_policy_contract.SplitPolicyRegisterParam) (*types.MutableTransaction, error) {
	tx, err := this.bc.BuildTx(this.contractAddress,"dtokenSellerPublish",
		[]interface{}{resourceId, ddo.ToBytes(), item.ToBytes(), splitPolicyParam.ToBytes()})
	return tx, err
}

func (this *DDXFKit) getPublishProductInfo(resourceId []byte) (*ProductInfoOnChain, error) {
	res, err := this.bc.PreInvoke(this.contractAddress, "getSellerItemInfo",
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

func (this *DDXFKit) Freeze(manager *ontology_go_sdk.Account,
	resourceId []byte) (common.Uint256, error) {
	tx, err := this.bc.BuildTx(this.contractAddress, manager, "buyDtoken",
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

func (this *DDXFKit) BuyDtoken(buyer, payer *ontology_go_sdk.Account, resourceId []byte,
	n int) (common.Uint256, error) {
	if payer == nil {
		payer = buyer
	}
	tx, err := this.bc.BuildTx(this.contractAddress, buyer, "buyDtoken",
		[]interface{}{resourceId, n, buyer.Address, payer.Address})
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

func (this *DDXFKit) BuyDtokenReward(buyer, payer *ontology_go_sdk.Account, resourceId []byte,
	n int, unitPrice int) (common.Uint256, error) {
	if payer == nil {
		payer = buyer
	}
	tx, err := this.bc.BuildTx(this.contractAddress, buyer, "buyDtoken",
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

func (this *DDXFKit) BuyAndUseToken(buyer, payer *ontology_go_sdk.Account, resourceId []byte,
	n int, tokenTemplate TokenTemplate) (common.Uint256, error) {
	if payer == nil {
		payer = buyer
	}
	tx, err := this.bc.BuildTx(this.contractAddress, buyer, "buyAndUseToken",
		[]interface{}{resourceId, n, buyer.Address, payer.Address,tokenTemplate.ToBytes()})
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

func (this *DDXFKit) UseToken(resourceId []byte, buyer *ontology_go_sdk.Account,
	tokenTemplate TokenTemplate, n int) (common.Uint256, error) {
	return this.bc.Invoke(this.contractAddress, buyer, "useToken",
		[]interface{}{resourceId, buyer.Address, tokenTemplate.ToBytes(), n})
}

func (this *DDXFKit) UseTokenByAgents(resource_id []byte, tokenOwner common.Address,
	agent *ontology_go_sdk.Account, tokenTemplate TokenTemplate, n int) (common.Uint256, error) {
	return this.bc.Invoke(this.contractAddress, agent, "useTokenByAgent",
		[]interface{}{resource_id, tokenOwner, agent.Address, tokenTemplate.ToBytes(), n})
}

func (this *DDXFKit) BuyDtokens(buyer *ontology_go_sdk.Account,
	param []ResourceIdAndN) (common.Uint256, error) {
	rids := make([]interface{}, len(param))
	for i := 0; i < len(param); i++ {
		rids[i] = param[i].ResourceId
	}
	ns := make([]interface{}, len(param))
	for i := 0; i < len(param); i++ {
		ns[i] = param[i].N
	}
	return this.bc.Invoke(this.contractAddress, buyer, "buyDtokens",
		[]interface{}{rids, ns, buyer.Address})
}

func (this *DDXFKit) SetAgents(resourceId []byte, account *ontology_go_sdk.Account,
	agents []common.Address, n int) (common.Uint256, error) {
	return this.bc.Invoke(this.contractAddress, account, "setAgents",
		[]interface{}{resourceId, account.Address, agents, n})
}

func (this *DDXFKit) SetTokenAgents(resourceId []byte, account *ontology_go_sdk.Account,
	agents []common.Address, tokenTemplate TokenTemplate, n int) (common.Uint256, error) {
	return this.bc.Invoke(this.contractAddress, account, "setTokenAgents",
		[]interface{}{resourceId, account.Address, agents, tokenTemplate.ToBytes(), n})
}

func (this *DDXFKit) AddAgents(resourceId []byte, account *ontology_go_sdk.Account,
	agents []common.Address, n int) (common.Uint256, error) {
	return this.bc.Invoke(this.contractAddress, account, "addAgents",
		[]interface{}{resourceId, account.Address, parseAddressArr(agents), n})
}

func parseAddressArr(addrs []common.Address) []interface{} {
	agentArr := make([]interface{}, len(addrs))
	for i := 0; i < len(addrs); i++ {
		agentArr[i] = addrs[i]
	}
	return agentArr
}

func (this *DDXFKit) AddTokenAgents(resourceId []byte, account *ontology_go_sdk.Account,
	agents []common.Address, tokenTemplate TokenTemplate, n int) (common.Uint256, error) {
	return this.bc.Invoke(this.contractAddress, account, "addTokenAgents",
		[]interface{}{resourceId, account.Address, tokenTemplate.ToBytes(), parseAddressArr(agents), n})
}

func (this *DDXFKit) RemoveAgents(resourceId []byte, account *ontology_go_sdk.Account,
	agents []common.Address) (common.Uint256, error) {
	return this.bc.Invoke(this.contractAddress, account, "removeAgents",
		[]interface{}{resourceId, account.Address, parseAddressArr(agents)})
}

func (this *DDXFKit) RemoveTokenAgents(resourceId []byte, tokenTemplate TokenTemplate, account *ontology_go_sdk.Account,
	agents []common.Address) (common.Uint256, error) {
	return this.bc.Invoke(this.contractAddress, account, "removeTokenAgents",
		[]interface{}{resourceId, tokenTemplate.ToBytes(), account.Address, parseAddressArr(agents)})
}
