package main

import (
	"encoding/csv"
	"fmt"
	"github.com/ont-bizsuite/ddxf-sdk"
	"github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/core/payload"
	"github.com/ontio/ontology/core/types"
	"github.com/ontio/ontology/smartcontract/states"
	"os"
	"strconv"
	"time"
	"github.com/ontio/ontology-crypto/keypair"
)

//mainnet
const THIRD_UNBOUND_TIME = 1599955215
//testnet
//const THIRD_UNBOUND_TIME = 1600029499

const THIRD_UNBOUND_ID = 3

var (
	admin *ontology_go_sdk.Account
)

func main() {
	testNet := "http://106.75.224.136:20336"
	testNet = ddxf_sdk.TestNet
	//testNet = "http://172.168.3.152:20336"
	//testNet = "http://127.0.0.1:20336"

	//testNet = "http://172.168.3.47:20336"
	//testNet = "http://113.31.112.154:20336"
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
	user, _ := wallet.GetAccountByAddress("Aejfo7ZX5PVpenRj23yChnyH64nf8T1zbu", pwd)

	fmt.Printf("user: %s \n", user.Address.ToBase58())
	pri, _ := keypair.WIF2Key([]byte("KySMiNrDDzFyUxfpK2hV9wFivq6hEmgB81D1UynhwjXjgd7xUZ88"))
	pub := pri.Public()
	add := types.AddressFromPubKey(pub)
	admin = &ontology_go_sdk.Account{
		PrivateKey: pri,
		PublicKey:  pub,
		Address:    add,
	}


	rec, err := readCsv("wing-3-new2.csv")
	if err != nil {
		fmt.Println(err)
		return
	}

	// 校验总和对不对
	if false {
		sum := uint64(0)
		for _, r := range rec {
			sum += r.amt
		}
		fmt.Println("sum:", sum)
		fmt.Println(len(rec))
		return
	}

	//d3878cbec82d763d9fd72ec1bff0d22ae703f2d7
	zeroPoolContract, _ := common.AddressFromHexString("47414f5d38fef4ea199502c7423716d208468a2c")

	if true {
		zeroPool := sdk.DefContract(zeroPoolContract)

		if false {
			txhash, err := zeroPool.Invoke("audit_settlement", admin, []interface{}{})
			if err != nil {
				fmt.Println(err)
				return
			}
			showNotify(sdk, "audit_settlement", txhash.ToHexString())
			return
		}
		if true {
			txhash, err := zeroPool.Invoke("staking", user, []interface{}{user.Address, 1})
			if err != nil {
				fmt.Println(err)
				return
			}
			showNotify(sdk, "staking", txhash.ToHexString())
			return
		}
		if true {
			txhash, err := zeroPool.Invoke("unstaking", admin, []interface{}{admin.Address, 1})
			if err != nil {
				fmt.Println(err)
				return
			}
			showNotify(sdk, "unstaking", txhash.ToHexString())
			return
		}
		return
	}

	flag := false
	for i := 0; i < len(rec)/20+1; i++ {
		contract := &states.WasmContractParam{}
		contract.Address = zeroPoolContract
		sink := common.NewZeroCopySink(nil)
		sink.WriteString("audit_user_wv")

		//时间戳
		sink.WriteI128(common.I128FromUint64(THIRD_UNBOUND_TIME))
		//治理结算id
		sink.WriteI128(common.I128FromUint64(THIRD_UNBOUND_ID))

		l := 20
		end := 20*i+20
		if 20*i+20 > len(rec) {
			end = len(rec)
			l = len(rec) - 20*i
		}


		// 地址和数量的数组
		sink.WriteVarUint(uint64(l))

		for j := 20 * i; j < end; j++ {
			sink.WriteAddress(rec[j].addr)
			sink.WriteI128(common.I128FromUint64(rec[j].amt))
			if j == len(rec) - 1 {
				flag = true
			}
		}

		sink.WriteBool(flag)
		contract.Args = sink.Bytes()
		invokePayload := &payload.InvokeCode{
			Code: common.SerializeToBytes(contract),
		}
		tx := &types.MutableTransaction{
			Payer:    admin.Address,
			GasPrice: 2500,
			GasLimit: 300000,
			TxType:   types.InvokeWasm,
			Nonce:    uint32(time.Now().Unix()),
			Payload:  invokePayload,
			Sigs:     nil,
		}
		sdk.GetOntologySdk().SignToTransaction(tx, admin)

		txhash, err := sdk.GetOntologySdk().SendTransaction(tx)
		if err != nil {
			fmt.Println(err)
			return
		}
		showNotify(sdk, "audit_user_wv", txhash.ToHexString())
	}
}

type UserRec struct {
	time uint64
	addr common.Address
	amt  uint64
}

func readCsv(fileName string) ([]UserRec, error) {
	f, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	reader := csv.NewReader(f)

	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	rec := make([]UserRec, len(records))
	for k, record := range records {
		if len(record) != 5 {
			panic("record is wrong")
		}
		addr, err := common.AddressFromBase58(record[0])
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		amt, err := strconv.ParseUint(record[2], 10, 64)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		time, err := strconv.ParseUint(record[3], 10, 64)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		rec[k] = UserRec{
			time: time,
			addr: addr,
			amt:  amt,
		}
	}
	return rec, nil
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
