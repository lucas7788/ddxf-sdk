package dtoken_contract

import (
	"fmt"

	"github.com/ont-bizsuite/ddxf-sdk/base_contract"
	"github.com/ont-bizsuite/ddxf-sdk/market_place_contract"
	"github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/core/types"
)

type DTokenKit struct {
	bc              *base_contract.BaseContract
	contractAddress common.Address
}

func NewDTokenKit(contractAddress common.Address, bc *base_contract.BaseContract) *DTokenKit {
	return &DTokenKit{
		contractAddress: contractAddress,
		bc:              bc,
	}
}

func (this *DTokenKit) SetContractAddr(addr common.Address) {
	this.contractAddress = addr
}

func (this *DTokenKit) SetMpContractAddr(contractAddress common.Address, admin *ontology_go_sdk.Account, addr common.Address) (common.Uint256, error) {
	if contractAddress == common.ADDRESS_EMPTY {
		contractAddress = this.contractAddress
	}
	return this.bc.Invoke(contractAddress, admin, "setMpContract",
		[]interface{}{addr})
}

func (this *DTokenKit) BuildSetMpContractAddrTx(contractAddress common.Address, mpContractAddress common.Address) (*types.MutableTransaction, error) {
	if contractAddress == common.ADDRESS_EMPTY {
		contractAddress = this.contractAddress
	}
	return this.bc.BuildTx(contractAddress, "setMpContract", []interface{}{mpContractAddress})

}

func (this *DTokenKit) GetMpContractAddr(contractAddress common.Address) (common.Address, error) {
	if contractAddress == common.ADDRESS_EMPTY {
		contractAddress = this.contractAddress
	}
	res, err := this.bc.PreInvoke(contractAddress, "getMpContract", []interface{}{})
	if err != nil {
		return common.ADDRESS_EMPTY, err
	}
	bs, err := res.ToByteArray()
	if err != nil {
		return common.ADDRESS_EMPTY, err
	}
	source := common.NewZeroCopySource(bs)
	addr, eof := source.NextAddress()
	if eof {
		return common.ADDRESS_EMPTY, fmt.Errorf("read address failed, eof: %v", eof)
	}
	return addr, nil
}

func (this *DTokenKit) CreateTokenTemplate(contractAddress common.Address, creator *ontology_go_sdk.Account,
	tt market_place_contract.TokenTemplate) (common.Uint256, error) {
	sink := common.NewZeroCopySink(nil)
	tt.Serialize(sink)
	if contractAddress == common.ADDRESS_EMPTY {
		contractAddress = this.contractAddress
	}
	return this.bc.Invoke(contractAddress, creator, "createTokenTemplate", []interface{}{creator.Address, sink.Bytes()})
}

func (this *DTokenKit) GetTokenTemplateById(contractAddress common.Address, tokenTemplateId []byte) (tt *market_place_contract.TokenTemplate, err error) {
	if contractAddress == common.ADDRESS_EMPTY {
		contractAddress = this.contractAddress
	}
	res, err := this.bc.PreInvoke(contractAddress, "getTokenTemplateById", []interface{}{tokenTemplateId})
	if err != nil {
		return
	}
	bs, err := res.ToByteArray()
	if err != nil {
		return
	}
	source := common.NewZeroCopySource(bs)
	hasVal, ir, eof := source.NextBool()
	if ir || eof {
		err = fmt.Errorf("ir: %v, eof: %v", ir, eof)
		return
	}
	if hasVal {
		tt = &market_place_contract.TokenTemplate{}
		err = tt.Deserialize(source)
	}
	return
}

func (this *DTokenKit) BuildCreateTokenTemplateTx(contractAddress common.Address, creator common.Address,
	tt market_place_contract.TokenTemplate) (*types.MutableTransaction, error) {
	sink := common.NewZeroCopySink(nil)
	tt.Serialize(sink)
	if contractAddress == common.ADDRESS_EMPTY {
		contractAddress = this.contractAddress
	}
	return this.bc.BuildTx(contractAddress, "createTokenTemplate", []interface{}{creator, sink.Bytes()})
}

func (this *DTokenKit) BuildUpdateTokenTemplateTx(contractAddress common.Address, tokenTemplateId []byte,
	tt market_place_contract.TokenTemplate) (*types.MutableTransaction, error) {
	sink := common.NewZeroCopySink(nil)
	tt.Serialize(sink)
	if contractAddress == common.ADDRESS_EMPTY {
		contractAddress = this.contractAddress
	}
	return this.bc.BuildTx(contractAddress, "updateTokenTemplate", []interface{}{tokenTemplateId, sink.Bytes()})
}

func (this *DTokenKit) BuildRemoveTokenTemplateTx(contractAddress common.Address, tokenTemplateId []byte) (*types.MutableTransaction, error) {
	if contractAddress == common.ADDRESS_EMPTY {
		contractAddress = this.contractAddress
	}
	return this.bc.BuildTx(contractAddress, "removeTokenTemplate", []interface{}{tokenTemplateId})
}

func (this *DTokenKit) AuthorizeTokenTemplate(contractAddress common.Address, creator *ontology_go_sdk.Account, tokenTemplateId []byte,
	authorizeAddr []common.Address) (common.Uint256, error) {
	if contractAddress == common.ADDRESS_EMPTY {
		contractAddress = this.contractAddress
	}
	return this.bc.Invoke(contractAddress, creator, "authorizeTokenTemplate",
		[]interface{}{tokenTemplateId, parseAddressArr(authorizeAddr)})
}

func (this *DTokenKit) BuildAuthorizeAddrTx(contractAddress common.Address, tokenTemplateId []byte,
	authorizeAddr []common.Address) (*types.MutableTransaction, error) {
	if contractAddress == common.ADDRESS_EMPTY {
		contractAddress = this.contractAddress
	}
	return this.bc.BuildTx(contractAddress, "authorizeTokenTemplate",
		[]interface{}{tokenTemplateId, parseAddressArr(authorizeAddr)})
}

func (this *DTokenKit) BuildRemAuthorizeAddrTx(contractAddress common.Address, tokenTemplateId []byte,
	authorizeAddr []common.Address) (*types.MutableTransaction, error) {
	if contractAddress == common.ADDRESS_EMPTY {
		contractAddress = this.contractAddress
	}
	return this.bc.BuildTx(contractAddress, "removeAuthorizeAddr",
		[]interface{}{tokenTemplateId, parseAddressArr(authorizeAddr)})
}

func (this *DTokenKit) GetAuthorizedAddr(contractAddress common.Address, tokenTemplateId []byte) ([]common.Address, error) {
	if contractAddress == common.ADDRESS_EMPTY {
		contractAddress = this.contractAddress
	}
	res, err := this.bc.PreInvoke(contractAddress, "getAuthorizedAddr", []interface{}{tokenTemplateId})
	if err != nil {
		return nil, err
	}
	arr, err := res.ToArray()
	if err != nil {
		return nil, err
	}
	addrs := make([]common.Address, 0)
	for _, item := range arr {
		bs, err := item.ToByteArray()
		if err != nil {
			return nil, err
		}
		addr, err := common.AddressParseFromBytes(bs)
		if err != nil {
			return nil, err
		}
		addrs = append(addrs, addr)
	}
	return addrs, nil
}

func (this *DTokenKit) GetTokenIdByTemplateId(contractAddress common.Address, templateId []byte) ([]byte, error) {
	if contractAddress == common.ADDRESS_EMPTY {
		contractAddress = this.contractAddress
	}
	res, err := this.bc.PreInvoke(contractAddress, "getTokenIdByTemplateId", []interface{}{templateId})
	if err != nil {
		return nil, err
	}
	return res.ToByteArray()
}
func (this *DTokenKit) GetTemplateIdByTokenId(contractAddress common.Address, tokenId []byte) ([]byte, error) {
	if contractAddress == common.ADDRESS_EMPTY {
		contractAddress = this.contractAddress
	}
	res, err := this.bc.PreInvoke(contractAddress, "getTemplateIdByTokenId", []interface{}{tokenId})
	if err != nil {
		return nil, err
	}
	return res.ToByteArray()
}

func (this *DTokenKit) GenerateDToken(contractAddress common.Address, acc *ontology_go_sdk.Account, tokenTemplateId []byte, n int) (common.Uint256, error) {
	if contractAddress == common.ADDRESS_EMPTY {
		contractAddress = this.contractAddress
	}
	return this.bc.Invoke(contractAddress, acc, "generateDToken", []interface{}{acc.Address, tokenTemplateId, n})
}

func (this *DTokenKit) BuildGenerateDTokenTx(contractAddress common.Address, acc common.Address, tokenTemplateId []byte,
	n int) (*types.MutableTransaction, error) {
	if contractAddress == common.ADDRESS_EMPTY {
		contractAddress = this.contractAddress
	}
	return this.bc.BuildTx(contractAddress, "generateDToken", []interface{}{acc, tokenTemplateId, n})
}

func (this *DTokenKit) GenerateDTokenForOther(contractAddress common.Address, acc *ontology_go_sdk.Account, to common.Address, tokenTemplateId []byte, n int) (common.Uint256, error) {
	if contractAddress == common.ADDRESS_EMPTY {
		contractAddress = this.contractAddress
	}
	return this.bc.Invoke(contractAddress, acc, "generateDTokenForOther", []interface{}{acc.Address, to, tokenTemplateId, n})
}

func (this *DTokenKit) BuildGenerateDTokenForOtherTx(contractAddress common.Address, acc, to common.Address, tokenTemplateId []byte,
	n int) (*types.MutableTransaction, error) {
	if contractAddress == common.ADDRESS_EMPTY {
		contractAddress = this.contractAddress
	}
	return this.bc.BuildTx(contractAddress, "generateDTokenForOther", []interface{}{acc, to, tokenTemplateId, n})
}

func (this *DTokenKit) BalanceOf(contractAddress common.Address, addr common.Address, tokenId []byte) (uint64, error) {
	if contractAddress == common.ADDRESS_EMPTY {
		contractAddress = this.contractAddress
	}
	res, err := this.bc.PreInvoke(contractAddress, "balanceOf", []interface{}{addr, tokenId})
	if err != nil {
		return 0, err
	}
	in, err := res.ToInteger()
	if err != nil {
		return 0, err
	}
	return in.Uint64(), nil
}

func (this *DTokenKit) DeleteToken(contractAddress common.Address, acc *ontology_go_sdk.Account, tokenId []byte) (common.Uint256, error) {
	if contractAddress == common.ADDRESS_EMPTY {
		contractAddress = this.contractAddress
	}
	return this.bc.Invoke(contractAddress, acc, "deleteToken", []interface{}{acc.Address, tokenId})
}

func (this *DTokenKit) GetAgentBalance(contractAddress common.Address, agent common.Address) (uint64, error) {
	if contractAddress == common.ADDRESS_EMPTY {
		contractAddress = this.contractAddress
	}
	res, err := this.bc.PreInvoke(contractAddress, "getAgentBalance", []interface{}{agent})
	if err != nil {
		return 0, err
	}
	in, err := res.ToInteger()
	if err != nil {
		return 0, err
	}
	return in.Uint64(), nil
}

func (this *DTokenKit) UseToken(contractAddress common.Address, buyer *ontology_go_sdk.Account,
	tokenId []byte, n int) (common.Uint256, error) {
	if contractAddress == common.ADDRESS_EMPTY {
		contractAddress = this.contractAddress
	}
	return this.bc.Invoke(contractAddress, buyer, "useToken",
		[]interface{}{buyer.Address, tokenId, n})
}

func (this *DTokenKit) TransferToken(contractAddress common.Address, from, to *ontology_go_sdk.Account,
	tokenId []byte, n int) (common.Uint256, error) {
	if contractAddress == common.ADDRESS_EMPTY {
		contractAddress = this.contractAddress
	}
	return this.bc.Invoke(contractAddress, from, "transfer",
		[]interface{}{from.Address, to.Address, tokenId, n})
}

func (this *DTokenKit) BuildTransferTokenTx(contractAddress common.Address, from, to common.Address,
	tokenId []byte, n int) (*types.MutableTransaction, error) {
	if contractAddress == common.ADDRESS_EMPTY {
		contractAddress = this.contractAddress
	}
	return this.bc.BuildTx(contractAddress, "transfer",
		[]interface{}{from, to, tokenId, n})
}

func (this *DTokenKit) BuildUseTokenTx(contractAddress common.Address, buyer common.Address, tokenId []byte, n int) (*types.MutableTransaction, error) {
	if contractAddress == common.ADDRESS_EMPTY {
		contractAddress = this.contractAddress
	}
	return this.bc.BuildTx(contractAddress, "useToken", []interface{}{buyer, tokenId, n})
}

func (this *DTokenKit) UseTokenByAgents(contractAddress common.Address, tokenOwner common.Address,
	agent *ontology_go_sdk.Account, tokenId []byte, n int) (common.Uint256, error) {
	if contractAddress == common.ADDRESS_EMPTY {
		contractAddress = this.contractAddress
	}
	return this.bc.Invoke(contractAddress, agent, "useTokenByAgent",
		[]interface{}{tokenOwner, agent.Address, tokenId, n})
}

func (this *DTokenKit) SetAgents(contractAddress common.Address, account *ontology_go_sdk.Account,
	agents []common.Address, n []int) (common.Uint256, error) {
	if contractAddress == common.ADDRESS_EMPTY {
		contractAddress = this.contractAddress
	}
	if len(agents) != len(n) {
		return common.UINT256_EMPTY, fmt.Errorf("the length of agents: %d not equal n length: %d", len(agents), len(n))
	}
	return this.bc.Invoke(contractAddress, account, "setAgents",
		[]interface{}{account.Address, parseAddressArr(agents), parseNArr(n)})
}

func (this *DTokenKit) SetTokenAgents(contractAddress common.Address, account *ontology_go_sdk.Account,
	agents []common.Address, tokenId []byte, n []int) (common.Uint256, error) {
	if len(agents) != len(n) {
		return common.UINT256_EMPTY, fmt.Errorf("the length of agents: %d not equal n length: %d", len(agents), len(n))
	}
	if contractAddress == common.ADDRESS_EMPTY {
		contractAddress = this.contractAddress
	}
	return this.bc.Invoke(this.contractAddress, account, "setTokenAgents",
		[]interface{}{account.Address, tokenId, parseAddressArr(agents), parseNArr(n)})
}

func (this *DTokenKit) BuildSetTokenAgentsTx(contractAddress common.Address, account common.Address,
	agents []common.Address, tokenId []byte, n []int) (*types.MutableTransaction, error) {
	if len(agents) != len(n) {
		return nil, fmt.Errorf("the length of agents: %d not equal n length: %d", len(agents), len(n))
	}
	if contractAddress == common.ADDRESS_EMPTY {
		contractAddress = this.contractAddress
	}
	return this.bc.BuildTx(this.contractAddress, "setTokenAgents",
		[]interface{}{account, tokenId, parseAddressArr(agents), parseNArr(n)})
}

func (this *DTokenKit) AddAgents(contractAddress common.Address, account *ontology_go_sdk.Account,
	agents []common.Address, n []int, tokenIds [][]byte) (common.Uint256, error) {
	if len(agents) != len(n) {
		return common.UINT256_EMPTY, fmt.Errorf("the length of agents: %d not equal n length: %d", len(agents), len(n))
	}
	tokenIdParam := make([]interface{}, 0)
	for _, item := range tokenIds {
		tokenIdParam = append(tokenIdParam, item)
	}
	if contractAddress == common.ADDRESS_EMPTY {
		contractAddress = this.contractAddress
	}
	return this.bc.Invoke(contractAddress, account, "addAgents",
		[]interface{}{account.Address, parseAddressArr(agents), parseNArr(n), tokenIdParam})
}

func parseAddressArr(addrs []common.Address) []interface{} {
	agentArr := make([]interface{}, len(addrs))
	for i := 0; i < len(addrs); i++ {
		agentArr[i] = addrs[i]
	}
	return agentArr
}
func parseNArr(ns []int) []interface{} {
	agentArr := make([]interface{}, len(ns))
	for i := 0; i < len(ns); i++ {
		agentArr[i] = ns[i]
	}
	return agentArr
}

func (this *DTokenKit) AddTokenAgents(contractAddress common.Address, account *ontology_go_sdk.Account,
	agents []common.Address, tokenId []byte, n []int) (common.Uint256, error) {
	if len(agents) != len(n) {
		return common.UINT256_EMPTY, fmt.Errorf("the length of agents: %d not equal n length: %d", len(agents), len(n))
	}
	if contractAddress == common.ADDRESS_EMPTY {
		contractAddress = this.contractAddress
	}
	return this.bc.Invoke(contractAddress, account, "addTokenAgents",
		[]interface{}{account.Address, tokenId, parseAddressArr(agents), parseNArr(n)})
}

func (this *DTokenKit) BuildAddTokenAgentsTx(contractAddress common.Address, account common.Address,
	agents []common.Address, tokenId []byte, n []int) (*types.MutableTransaction, error) {
	if len(agents) != len(n) {
		return nil, fmt.Errorf("the length of agents: %d not equal n length: %d", len(agents), len(n))
	}
	if contractAddress == common.ADDRESS_EMPTY {
		contractAddress = this.contractAddress
	}
	return this.bc.BuildTx(contractAddress, "addTokenAgents",
		[]interface{}{account, tokenId, parseAddressArr(agents), parseNArr(n)})
}

func (this *DTokenKit) RemoveAgents(contractAddress common.Address, account *ontology_go_sdk.Account,
	agents []common.Address, tokenIds [][]byte) (common.Uint256, error) {
	tokenIdParam := make([]interface{}, 0)
	for _, item := range tokenIds {
		tokenIdParam = append(tokenIdParam, item)
	}
	if contractAddress == common.ADDRESS_EMPTY {
		contractAddress = this.contractAddress
	}
	return this.bc.Invoke(contractAddress, account, "removeAgents",
		[]interface{}{account.Address, parseAddressArr(agents), tokenIdParam})
}

func (this *DTokenKit) RemoveTokenAgents(contractAddress common.Address, tokenId []byte, account *ontology_go_sdk.Account,
	agents []common.Address) (common.Uint256, error) {
	if contractAddress == common.ADDRESS_EMPTY {
		contractAddress = this.contractAddress
	}
	return this.bc.Invoke(contractAddress, account, "removeTokenAgents",
		[]interface{}{account.Address, tokenId, parseAddressArr(agents)})
}

func (this *DTokenKit) BuildRemoveTokenAgentsTx(contractAddress common.Address, tokenId []byte, account common.Address,
	agents []common.Address) (*types.MutableTransaction, error) {
	if contractAddress == common.ADDRESS_EMPTY {
		contractAddress = this.contractAddress
	}
	return this.bc.BuildTx(contractAddress, "removeTokenAgents",
		[]interface{}{account, tokenId, parseAddressArr(agents)})
}
