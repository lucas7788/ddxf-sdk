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

func (this *DTokenKit) SetMpContractAddr(admin *ontology_go_sdk.Account, addr common.Address) (common.Uint256, error) {
	return this.bc.Invoke(this.contractAddress, admin, "setMpContract",
		[]interface{}{addr})
}

func (this *DTokenKit) GetMpContractAddr() (common.Address, error) {
	res, err := this.bc.PreInvoke(this.contractAddress, "getMpContract", []interface{}{})
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

func (this *DTokenKit) CreateTokenTemplate(creator *ontology_go_sdk.Account,
	tt market_place_contract.TokenTemplate) (common.Uint256, error) {
	sink := common.NewZeroCopySink(nil)
	tt.Serialize(sink)
	return this.bc.Invoke(this.contractAddress, creator, "createTokenTemplate", []interface{}{creator.Address, sink.Bytes()})
}

func (this *DTokenKit) GetTokenTemplateById(tokenTemplateId []byte) (tt *market_place_contract.TokenTemplate, err error) {
	res, err := this.bc.PreInvoke(this.contractAddress, "getTokenTemplateById", []interface{}{tokenTemplateId})
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

func (this *DTokenKit) BuildCreateTokenTemplateTx(creator common.Address,
	tt market_place_contract.TokenTemplate) (*types.MutableTransaction, error) {
	sink := common.NewZeroCopySink(nil)
	tt.Serialize(sink)
	return this.bc.BuildTx(this.contractAddress, "createTokenTemplate", []interface{}{creator, sink.Bytes()})
}

func (this *DTokenKit) BuildUpdateTokenTemplateTx(tokenTemplateId []byte,
	tt market_place_contract.TokenTemplate) (*types.MutableTransaction, error) {
	sink := common.NewZeroCopySink(nil)
	tt.Serialize(sink)
	return this.bc.BuildTx(this.contractAddress, "updateTokenTemplate", []interface{}{tokenTemplateId, sink.Bytes()})
}

func (this *DTokenKit) BuildRemoveTokenTemplateTx(tokenTemplateId []byte) (*types.MutableTransaction, error) {
	return this.bc.BuildTx(this.contractAddress, "removeTokenTemplate", []interface{}{tokenTemplateId})
}

func (this *DTokenKit) AuthorizeTokenTemplate(creator *ontology_go_sdk.Account, tokenTemplateId []byte,
	authorizeAddr []common.Address) (common.Uint256, error) {
	return this.bc.Invoke(this.contractAddress, creator, "authorizeTokenTemplate",
		[]interface{}{tokenTemplateId, parseAddressArr(authorizeAddr)})
}

func (this *DTokenKit) BuildAuthorizeAddrTx(tokenTemplateId []byte,
	authorizeAddr []common.Address) (*types.MutableTransaction, error) {
	return this.bc.BuildTx(this.contractAddress, "authorizeTokenTemplate",
		[]interface{}{tokenTemplateId, parseAddressArr(authorizeAddr)})
}

func (this *DTokenKit) BuildRemAuthorizeAddrTx(tokenTemplateId []byte,
	authorizeAddr []common.Address) (*types.MutableTransaction, error) {
	return this.bc.BuildTx(this.contractAddress, "removeAuthorizeAddr",
		[]interface{}{tokenTemplateId, parseAddressArr(authorizeAddr)})
}

func (this *DTokenKit) GetAuthorizedAddr(tokenTemplateId []byte) ([]common.Address, error) {
	res, err := this.bc.PreInvoke(this.contractAddress, "getAuthorizedAddr", []interface{}{tokenTemplateId})
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

func (this *DTokenKit) GetTokenIdByTemplateId(templateId []byte) ([]byte, error) {
	res, err := this.bc.PreInvoke(this.contractAddress, "getTokenIdByTemplateId", []interface{}{templateId})
	if err != nil {
		return nil, err
	}
	return res.ToByteArray()
}
func (this *DTokenKit) GetTemplateIdByTokenId(tokenId []byte) ([]byte, error) {
	res, err := this.bc.PreInvoke(this.contractAddress, "getTemplateIdByTokenId", []interface{}{tokenId})
	if err != nil {
		return nil, err
	}
	return res.ToByteArray()
}

func (this *DTokenKit) GenerateDToken(acc *ontology_go_sdk.Account, tokenTemplateId []byte, n int) (common.Uint256, error) {
	return this.bc.Invoke(this.contractAddress, acc, "generateDToken", []interface{}{acc.Address, tokenTemplateId, n})
}

func (this *DTokenKit) BuildGenerateDTokenTx(acc common.Address, tokenTemplateId []byte,
	n int) (*types.MutableTransaction, error) {
	return this.bc.BuildTx(this.contractAddress, "generateDToken", []interface{}{acc, tokenTemplateId, n})
}

func (this *DTokenKit) GenerateDTokenForOther(acc *ontology_go_sdk.Account, to common.Address, tokenTemplateId []byte, n int) (common.Uint256, error) {
	return this.bc.Invoke(this.contractAddress, acc, "generateDTokenForOther", []interface{}{acc.Address, to, tokenTemplateId, n})
}

func (this *DTokenKit) BuildGenerateDTokenForOtherTx(acc, to common.Address, tokenTemplateId []byte,
	n int) (*types.MutableTransaction, error) {
	return this.bc.BuildTx(this.contractAddress, "generateDTokenForOther", []interface{}{acc, to, tokenTemplateId, n})
}

func (this *DTokenKit) BalanceOf(addr common.Address, tokenId []byte) (uint64, error) {
	res, err := this.bc.PreInvoke(this.contractAddress, "balanceOf", []interface{}{addr, tokenId})
	if err != nil {
		return 0, err
	}
	in, err := res.ToInteger()
	if err != nil {
		return 0, err
	}
	return in.Uint64(), nil
}

func (this *DTokenKit) DeleteToken(acc *ontology_go_sdk.Account, tokenId []byte) (common.Uint256, error) {
	return this.bc.Invoke(this.contractAddress, acc, "deleteToken", []interface{}{acc.Address, tokenId})
}

func (this *DTokenKit) GetAgentBalance(agent common.Address) (uint64, error) {
	res, err := this.bc.PreInvoke(this.contractAddress, "getAgentBalance", []interface{}{agent})
	if err != nil {
		return 0, err
	}
	in, err := res.ToInteger()
	if err != nil {
		return 0, err
	}
	return in.Uint64(), nil
}

func (this *DTokenKit) UseToken(buyer *ontology_go_sdk.Account,
	tokenId []byte, n int) (common.Uint256, error) {
	return this.bc.Invoke(this.contractAddress, buyer, "useToken",
		[]interface{}{buyer.Address, tokenId, n})
}

func (this *DTokenKit) TransferToken(from, to *ontology_go_sdk.Account,
	tokenId []byte, n int) (common.Uint256, error) {
	return this.bc.Invoke(this.contractAddress, from, "transfer",
		[]interface{}{from.Address, to.Address, tokenId, n})
}

func (this *DTokenKit) BuildTransferTokenTx(from, to common.Address,
	tokenId []byte, n int) (*types.MutableTransaction, error) {
	return this.bc.BuildTx(this.contractAddress, "transfer",
		[]interface{}{from, to, tokenId, n})
}

func (this *DTokenKit) BuildUseTokenTx(buyer common.Address, tokenId []byte, n int) (*types.MutableTransaction, error) {
	return this.bc.BuildTx(this.contractAddress, "useToken", []interface{}{buyer, tokenId, n})
}

func (this *DTokenKit) UseTokenByAgents(tokenOwner common.Address,
	agent *ontology_go_sdk.Account, tokenId []byte, n int) (common.Uint256, error) {
	return this.bc.Invoke(this.contractAddress, agent, "useTokenByAgent",
		[]interface{}{tokenOwner, agent.Address, tokenId, n})
}

func (this *DTokenKit) SetAgents(account *ontology_go_sdk.Account,
	agents []common.Address, n []int) (common.Uint256, error) {
	if len(agents) != len(n) {
		return common.UINT256_EMPTY, fmt.Errorf("the length of agents: %d not equal n length: %d", len(agents), len(n))
	}
	return this.bc.Invoke(this.contractAddress, account, "setAgents",
		[]interface{}{account.Address, parseAddressArr(agents), parseNArr(n)})
}

func (this *DTokenKit) SetTokenAgents(account *ontology_go_sdk.Account,
	agents []common.Address, tokenId []byte, n []int) (common.Uint256, error) {
	if len(agents) != len(n) {
		return common.UINT256_EMPTY, fmt.Errorf("the length of agents: %d not equal n length: %d", len(agents), len(n))
	}
	return this.bc.Invoke(this.contractAddress, account, "setTokenAgents",
		[]interface{}{account.Address, tokenId, parseAddressArr(agents), parseNArr(n)})
}

func (this *DTokenKit) BuildSetTokenAgentsTx(account common.Address,
	agents []common.Address, tokenId []byte, n []int) (*types.MutableTransaction, error) {
	if len(agents) != len(n) {
		return nil, fmt.Errorf("the length of agents: %d not equal n length: %d", len(agents), len(n))
	}
	return this.bc.BuildTx(this.contractAddress, "setTokenAgents",
		[]interface{}{account, tokenId, parseAddressArr(agents), parseNArr(n)})
}

func (this *DTokenKit) AddAgents(account *ontology_go_sdk.Account,
	agents []common.Address, n []int, tokenIds [][]byte) (common.Uint256, error) {
	if len(agents) != len(n) {
		return common.UINT256_EMPTY, fmt.Errorf("the length of agents: %d not equal n length: %d", len(agents), len(n))
	}
	tokenIdParam := make([]interface{}, 0)
	for _, item := range tokenIds {
		tokenIdParam = append(tokenIdParam, item)
	}
	return this.bc.Invoke(this.contractAddress, account, "addAgents",
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

func (this *DTokenKit) AddTokenAgents(account *ontology_go_sdk.Account,
	agents []common.Address, tokenId []byte, n []int) (common.Uint256, error) {
	if len(agents) != len(n) {
		return common.UINT256_EMPTY, fmt.Errorf("the length of agents: %d not equal n length: %d", len(agents), len(n))
	}
	return this.bc.Invoke(this.contractAddress, account, "addTokenAgents",
		[]interface{}{account.Address, tokenId, parseAddressArr(agents), parseNArr(n)})
}

func (this *DTokenKit) BuildAddTokenAgentsTx(account common.Address,
	agents []common.Address, tokenId []byte, n []int) (*types.MutableTransaction, error) {
	if len(agents) != len(n) {
		return nil, fmt.Errorf("the length of agents: %d not equal n length: %d", len(agents), len(n))
	}
	return this.bc.BuildTx(this.contractAddress, "addTokenAgents",
		[]interface{}{account, tokenId, parseAddressArr(agents), parseNArr(n)})
}

func (this *DTokenKit) RemoveAgents(account *ontology_go_sdk.Account,
	agents []common.Address, tokenIds [][]byte) (common.Uint256, error) {
	tokenIdParam := make([]interface{}, 0)
	for _, item := range tokenIds {
		tokenIdParam = append(tokenIdParam, item)
	}
	return this.bc.Invoke(this.contractAddress, account, "removeAgents",
		[]interface{}{account.Address, parseAddressArr(agents), tokenIdParam})
}

func (this *DTokenKit) RemoveTokenAgents(tokenId []byte, account *ontology_go_sdk.Account,
	agents []common.Address) (common.Uint256, error) {
	return this.bc.Invoke(this.contractAddress, account, "removeTokenAgents",
		[]interface{}{account.Address, tokenId, parseAddressArr(agents)})
}

func (this *DTokenKit) BuildRemoveTokenAgentsTx(tokenId []byte, account common.Address,
	agents []common.Address) (*types.MutableTransaction, error) {
	return this.bc.BuildTx(this.contractAddress, "removeTokenAgents",
		[]interface{}{account, tokenId, parseAddressArr(agents)})
}
