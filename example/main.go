package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ont-bizsuite/ddxf-sdk"
	"github.com/ont-bizsuite/ddxf-sdk/example/utils"
	"github.com/ont-bizsuite/ddxf-sdk/market_place_contract"
	"github.com/ont-bizsuite/ddxf-sdk/split_policy_contract"

	"github.com/ontio/ontology-crypto/signature"
	"github.com/ontio/ontology-go-sdk"
	utils2 "github.com/ontio/ontology-go-sdk/utils"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/core/payload"
	"github.com/ontio/ontology/core/types"
	"github.com/ontio/ontology/smartcontract/states"
	"github.com/zhiqiangxu/ddxf"
	"io/ioutil"
	"math/rand"
	"strconv"
	"time"
	"github.com/ontio/ontology-crypto/keypair"
)

var (
	admin         *ontology_go_sdk.Account
	seller        *ontology_go_sdk.Account
	buyer         *ontology_go_sdk.Account
	agent         *ontology_go_sdk.Account
	payer         *ontology_go_sdk.Account
	gasPrice      = uint64(2500)
	tokenTemplate *market_place_contract.TokenTemplate
)

func main() {
	testNet := "http://106.75.224.136:20336"
	testNet = ddxf_sdk.TestNet
	//testNet = "http://172.168.3.152:20336"
	//testNet = "http://127.0.0.1:20336"

	//testNet = "http://172.168.3.47:20336"
	//testNet = "http://113.31.112.154:20336"
	testNet = ddxf_sdk.MainNet
	sdk := ddxf_sdk.NewDdxfSdk(testNet)
	//106.75.224.136

	if false {
		code, err := sdk.GetOntologySdk().GetSmartContract("2e532b8ff89c58d340add366437a4572574114ff")
		fmt.Println(code, err)
		evt, err := sdk.GetSmartCodeEvent("850a6c35ef54c7dc276feee16e01c8a6e20a17cfb0df436179e868b65ce583d3")
		if err != nil {
			fmt.Println(err)
			return
		}
		data, _ := json.Marshal(evt)
		fmt.Println(string(data))
		return
	}

	pwd := []byte("123456")
	ontSdk := sdk.GetOntologySdk()
	wallet, err := ontSdk.OpenWallet("./wallet.dat")
	if err != nil {
		fmt.Printf("error in ReadFile:%s\n", err)
		return
	}

	admin, _ = wallet.GetAccountByAddress("AYnhakv7kC9R5ppw65JoE2rt6xDzCjCTvD", pwd)
	seller, _ = wallet.GetAccountByAddress("Aejfo7ZX5PVpenRj23yChnyH64nf8T1zbu", pwd)
	buyer, _ = wallet.GetAccountByAddress("AHhXa11suUgVLX1ZDFErqBd3gskKqLfa5N", pwd)
	agent, _ = wallet.GetAccountByAddress("ANb3bf1b67WP2ZPh5HQt4rkrmphMJmMCMK", pwd)
	payer, _ = wallet.GetAccountByAddress("AQCQ3Krh6qxeWKKRACNehA8kAATHxoQNWJ", pwd)

	if false {
		jiaoyisuo,_ := common.AddressFromBase58("AbbBBzXZE1EyYgYCusnZCzseyajWXXCsTQ")
		txhash, err := sdk.GetOntologySdk().Native.Ong.Transfer(2500,20000,seller,seller,jiaoyisuo,50 * 1000000000)
		fmt.Println(err)
		fmt.Println(txhash.ToHexString())
		return
	}
	fmt.Println(agent.Address.ToBase58())

	pri, _ := keypair.WIF2Key([]byte("KySMiNrDDzFyUxfpK2hV9wFivq6hEmgB81D1UynhwjXjgd7xUZ88"))
	pub := pri.Public()
	add := types.AddressFromPubKey(pub)

	admin = &ontology_go_sdk.Account{
		PrivateKey: pri,
		PublicKey:  pub,
		Address:    add,
	}

	fmt.Println(admin.Address.ToBase58())

	if false {
		utils.DeployGlobalParamContract(sdk, seller, 2500)
		return
	}

	if false {
		//e11bcb67d37f1797de02dc890344db2bd4440d4a
		//acbabd9ab5281808e1bcda70fcac05668ef00bc8
		//5a85897b8ea91da01bb965383f5cc2525bdaa663
		//faf99fb13460897b90a970b7921d873e72ae2823
		utils.DeployZeroPoolContract(sdk, admin, 2500)
		return
	}

	//270d19620533881997920364e2cbfc3ee334c8f8
	//9aa4deaf4ddb215af123f59b17b26b131375e5a4 old
	zeroContractAddr, _ := common.AddressFromHexString("20e2322eb48f0aeda8f9208280aaa6d15c46d769")
	global_param, _ := common.AddressFromHexString("a7e169ef9ab7930d573850c1152780d8d7530458")

	//df7249c2699d94ff74cc8743fd265f0576c3f973
	wingCon, _ := common.AddressFromHexString("ef79ae5d85ec819cfd51eb23d4c242fa70f713e7")

	//222db9efdd987484b246b7318d01a1655e7eceb3
	//09956946a55f61e4280321f7372c49188a911e96
	gov_addr, _ := common.AddressFromHexString("38f957755afd12b2149a3490af00a13385f2073e")


	//zeroPool升级前 wing余额 为145238856226
	if false {
		//wing := sdk.DefContract(wingCon)
		res, err := sdk.GetOntologySdk().NeoVM.PreExecInvokeNeoVMContract(wingCon,[]interface{}{"balanceOf", []interface{}{zeroContractAddr}})
		if err != nil {
			fmt.Println(err)
			return
		}
		bs, _ := res.Result.ToInteger()
		fmt.Println(bs.Uint64())
		fmt.Println(sdk.GetOntologySdk().Native.Ont.BalanceOf(gov_addr))
		return
	}

	if false {
		global := sdk.DefContract(global_param)
		txhash,err := global.Invoke("set_value", admin, []interface{}{[]byte("wing-dao-address"),gov_addr[:], admin.Address})
		if err != nil {
			fmt.Println(err)
			return
		}
		showNotify(sdk, "set_value",txhash.ToHexString())
		res, err := global.PreInvoke("get_value", []interface{}{[]byte("wing-dao-address")})
		if err != nil {
			fmt.Println(err)
			return
		}
		bs, _ := res.ToByteArray()
		fmt.Println(hex.EncodeToString(bs))
		return
	}

	if true {

		zeroPool := sdk.DefContract(zeroContractAddr)

		if false {
			wasmFile := "/Users/sss/dev/rust_project/zero_pool/output/zero_pool.wasm"
			code, _ := ioutil.ReadFile(wasmFile)
			txhash, err := zeroPool.Invoke("migrate",admin, []interface{}{code, 3,"zeroPool","version","author","email","desc"})
			if err != nil {
				fmt.Println(err)
				return
			}
			showNotify(sdk, "migrate", txhash.ToHexString())
			return
		}

		//sdk.SetGasPrice(0)
		if false {
			txhash, err := zeroPool.Invoke("init", seller, []interface{}{global_param, wingCon})
			if err != nil {
				fmt.Println(err)
				return
			}
			showNotify(sdk, "init", txhash.ToHexString())
			return
		}
		if false {
			txhash, err := zeroPool.Invoke("set_wing_address", admin, []interface{}{wingCon})
			if err != nil {
				fmt.Println(err)
				return
			}
			showNotify(sdk, "set_wing_address", txhash.ToHexString())
			return
		}
		if true {
			//AeFg9DyybfS8dYqUD79g9bAE3DC4h6nQCC   23750000, 154721428
			addr ,_:=common.AddressFromBase58("AakW1sjaVdX15xYNeMqgpZFuiQCBgQYpBi")
			method := "get_global_param_addr"
			method = "balanceOf"
			method = "get_unbound_wing"
			res, err := zeroPool.PreInvoke(method, []interface{}{addr})
			if err != nil {
				fmt.Println(err)
				return
			}
			bs, _ := res.ToByteArray()
			source :=common.NewZeroCopySource(bs)
			data, _ := source.NextI128()
			fmt.Println(data.ToNumString())
			return
		}
		if false {
			txhash, err := zeroPool.Invoke("transfer", admin, []interface{}{admin.Address, seller.Address,100})
			if err != nil {
				fmt.Println(err)
				return
			}
			showNotify(sdk, "transfer", txhash.ToHexString())
			return
		}
		if false {
			//升级前用户 质押的ont 数量
			//addr ,_:=common.AddressFromBase58("AeFg9DyybfS8dYqUD79g9bAE3DC4h6nQCC")
			res, err := zeroPool.PreInvoke("get_staking_balance", []interface{}{seller.Address})
			if err != nil {
				fmt.Println(err)
				return
			}
			bs, _ := res.ToByteArray()
			fmt.Println(hex.EncodeToString(bs))
			source := common.NewZeroCopySource(bs)
			i128, _ := source.NextI128()
			fmt.Println(i128.ToNumString())
			return
		}

		if false {
			txhash, err := zeroPool.Invoke("staking", seller, []interface{}{seller.Address, 1})
			if err != nil {
				fmt.Println(err)
				return
			}
			showNotify(sdk, "staking", txhash.ToHexString())
			return
		}
		if true {
			txhash, err := zeroPool.Invoke("unstaking", seller, []interface{}{seller.Address, 3288})
			if err != nil {
				fmt.Println(err)
				return
			}
			showNotify(sdk, "unstaking", txhash.ToHexString())
			return
		}
		if true {
			txhash, err := zeroPool.Invoke("withdraw_wing", admin, []interface{}{admin.Address})
			if err != nil {
				fmt.Println(err)
				return
			}
			showNotify(sdk, "withdraw_wing", txhash.ToHexString())
			return
		}
		return
	}

	//治理合约初始化
	if true {

		//global_param, _ := common.AddressFromHexString("a7e169ef9ab7930d573850c1152780d8d7530458")
		//wingCon, _ := common.AddressFromHexString("6d940806345b271b01efc65201ed65d940ae7db1")

		gov := sdk.DefContract(gov_addr)
		if false {
			txhash, err := gov.Invoke("register_pool", admin, []interface{}{zeroContractAddr, 100})
			if err != nil {
				fmt.Println(err)
				return
			}
			showNotify(sdk, "register_pool", txhash.ToHexString())
			return
		}
		if false {
			txhash, err := gov.Invoke("unbound_token", admin, []interface{}{})
			if err != nil {
				fmt.Println(err)
				return
			}
			showNotify(sdk, "unbound_token", txhash.ToHexString())
			//return
		}
		if false {
			txhash, err := gov.Invoke("unbound_to_pool", admin, []interface{}{})
			if err != nil {
				fmt.Println(err)
				return
			}
			showNotify(sdk, "unbound_to_pool", txhash.ToHexString())
			return
		}

		if false {
			ont, _ := common.AddressFromBase58("AFmseVrdL9f9oyCzZefL9tG6UbvhUMqNMV")
			token := utils.NewToken("ONT", 1, ont)
			sink := common.NewZeroCopySink(nil)
			sink.WriteString("set_exchange_rate")
			token.Serialize(sink)
			sink.WriteI128(common.I128FromUint64(1))

			contract := &states.WasmContractParam{}
			contract.Address = gov_addr
			//bf := bytes.NewBuffer(nil)
			argbytes := sink.Bytes()
			contract.Args = argbytes

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
			//bs, _ := res.Result.ToByteArray()
			//source := common.NewZeroCopySource(bs)
			//d, _ := source.NextI128()
			//fmt.Println(d.ToNumString())
			showNotify(sdk, "set_exchange_rate", txhash.ToHexString())
			return
		}
		if true {
			res, err := gov.PreInvoke("query_unbound_to_pool", []interface{}{zeroContractAddr, 3})
			if err != nil {
				fmt.Println(err)
				return
			}
			bs, _ := res.ToByteArray()

			fmt.Println(hex.EncodeToString(bs))
			return
		}
		if true {
			res, err := gov.PreInvoke("query_unbound_to_pool_count", []interface{}{zeroContractAddr})
			if err != nil {
				fmt.Println(err)
				return
			}
			bs, _ := res.ToByteArray()
			source := common.NewZeroCopySource(bs)
			fmt.Println(source.NextI128())
			fmt.Println(hex.EncodeToString(bs))
			return
		}
		return
	}

	if false {
		//code, err := sdk.GetOntologySdk().GetSmartContract("c485c5f4671e7acc024e63294a7573bd083b7e6d")
		//if err != nil {
		//	fmt.Println(err)
		//	return
		//}
		//fmt.Println(*code)
		addr, _ := common.AddressFromHexString("c485c5f4671e7acc024e63294a7573bd083b7e6d")
		con := sdk.DefContract(addr)
		txhash, err := con.Invoke("init", seller, []interface{}{})
		if err != nil {
			fmt.Println(err)
			return
		}
		evt, err := sdk.GetSmartCodeEvent(txhash.ToHexString())
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(evt)
		return
	}

	if false {
		utils.DeployOep4Contract(sdk, seller, 2500)
		return
	}

	if false {
		addr, _ := common.AddressFromHexString("e01d500ed0c1719b7750367ae59b4b2d308d1ceb")
		con := sdk.DefMpKit()
		con.SetContractAddress(addr)
		txhash, err := con.Update(common.ADDRESS_EMPTY, seller, []byte("01010100"),
			market_place_contract.ResourceDDO{
				Manager: seller.Address,
			}, market_place_contract.DTokenItem{}, split_policy_contract.SplitPolicyRegisterParam{})
		if err != nil {
			fmt.Println(err)
			return
		}
		evt, err := sdk.GetSmartCodeEvent(txhash.ToHexString())
		fmt.Println(evt, err)
		return
	}

	errCount := 0
	if false {
		Decimal := 100000000
		addr, _ := common.AddressFromHexString("02ec93fc885b00c2a67153d8d884959be8a97817")

		mp := sdk.DefMpKit()

		mp.SetContractAddress(addr)

		data, err := ioutil.ReadFile("./resourceId.txt")
		if err != nil {
			fmt.Println(err)
			return
		}
		res := make(map[string]interface{})
		err = json.Unmarshal(data, &res)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(len(res))
		con := sdk.DefContract(addr)

		for k, v := range res {
			id, err := hex.DecodeString(k)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("id:", string(id))
			if errCount > 10 {
				return
			}
			txhash, err := con.Invoke("buyRewardAndUseToken", seller, []interface{}{id, int(v.(float64)), seller.Address, seller.Address, Decimal, []byte{}})
			if err != nil {
				fmt.Println(err)
				errCount += 1
				continue
			}
			fmt.Println(txhash.ToHexString())
			return
		}
		return
	}

	//destroy contract
	if false {
		contractAddr, _ := common.AddressFromHexString("711ea85080759a8ebfc16088fa8bc1d2caba8a7e")
		con := sdk.DefContract(contractAddr)
		txhash, err := con.Invoke("destroy", seller, []interface{}{})
		if err != nil {
			fmt.Println(err)
			return
		}
		evt, err := sdk.GetSmartCodeEvent(txhash.ToHexString())
		fmt.Println(evt, err)
		return
	}

	if false {
		if true {
			dtoken, _ := common.AddressFromHexString("466b94488bf2ad1b1eec0ae7e49e40708e71a35d")
			split, _ := common.AddressFromHexString("f024034fe7e5ea69c53cede4774bd1dad566234f")
			addr, _ := common.AddressFromHexString("e01d500ed0c1719b7750367ae59b4b2d308d1ceb")
			mp := sdk.DefMpKit()
			mp.SetContractAddress(addr)
			bs, err := hex.DecodeString("7265736f757263655f69645f62386635613265352d666262382d346162322d393937642d366532326433323436623662")
			res, err := mp.GetPublishProductInfo(common.ADDRESS_EMPTY, bs)
			fmt.Println(string(bs))
			fmt.Println(res)
			sdk.GetOntologySdk().NewRpcClient().SetAddress("http://172.168.3.47:20336")
			txhash, err := mp.Init(common.ADDRESS_EMPTY, seller, dtoken, split)
			if err != nil {
				fmt.Println(err)
				return
			}
			evt, err := sdk.GetSmartCodeEvent(txhash.ToHexString())
			fmt.Println(err)
			fmt.Println(evt)
			return
			//dtokenContract,err := mp.PreInvoke("init",[]interface{}{})
			//if err != nil {
			//	fmt.Println(err)
			//	return
			//}
			//bs, err := dtokenContract.ToByteArray()
			//if err != nil {
			//	fmt.Println(err)
			//	return
			//}
			//addr, err = common.AddressParseFromBytes(bs)
			//fmt.Println(addr.ToHexString(), err)
			//dtoken 466b94488bf2ad1b1eec0ae7e49e40708e71a35d
			//split f024034fe7e5ea69c53cede4774bd1dad566234f
			return
		}

		code, _ := sdk.GetOntologySdk().GetSmartContract("f024034fe7e5ea69c53cede4774bd1dad566234f")
		sdk.GetOntologySdk().NewRpcClient().SetAddress("http://172.168.3.47:20336")
		//desc := "Sub-contract of DDXF series, provides a place for token/DToken owner to manage their published items, and enable token exchange between token owner and token acquiers."
		txHash, err := sdk.GetOntologySdk().WasmVM.DeployWasmVMSmartContract(
			0, 28800000, seller, hex.EncodeToString(code.GetRawCode()),
			"DDXF - dtoken", "v1.0", "", "", "")
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(txHash.ToHexString())
		evt, err := sdk.GetSmartCodeEvent(txHash.ToHexString())
		fmt.Println(err)
		fmt.Println(evt)
		return
	}

	if false {
		txHex := "00d2c75d365fc409000000000000c03ae50100000000fbe02b027e61a6d7602f26cfa9487fa58ef9ee72fd6e01bfdef23b4cfcf8e988670ff4bd762f1cd04bb10cfd57011a7265675f69645f6164645f6174747269627574655f617272617901fd38012a6469643a6f6e743a54397a715277475341487965697537784d364a3461573242737679687763317a446b012a6469643a6f6e743a4150484e504c7a3275314a5558794438726872794c616f51725734364a335036793201000000012a6469643a6f6e743a4150484e504c7a3275314a5558794438726872794c616f51725734364a33503679320100000002086461746148617368103634363137343631343836313733363801730c646174614d657461486173688036343338363533313633333436313336363233353330363536343633333733383634333236363337333736363334333933303336363433383337333136353330333136363633333433303332333933343330333133313333333336343336333933363633333433363331333733383333363533353636333833363632333033300173000141409d12dbde1896728b0adfab32e3b8c9ae4d004ce4cd6ff47f4fcdc8374f3f2a65a753e35f3a1990f22ec0bf2477ca9d9860a3c128673ad4e7d6a652d15e35466e23210240cf95b7738a102a554f83f6202fb00ed69f4354a69bb832d8df0938512adde9ac"
		tx, err := utils2.TransactionFromHexString(txHex)
		fmt.Println(err)
		mutTx, err := tx.IntoMutable()
		fmt.Println(err)
		txhash, err := ontSdk.SendTransaction(mutTx)
		fmt.Println(err)
		time.Sleep(6 * time.Second)
		fmt.Println(txhash.ToHexString())
		evt, err := ontSdk.GetSmartContractEvent(txhash.ToHexString())
		fmt.Println(err)
		fmt.Println(evt)
		return
	}

	if false {
		txhash, err := ontSdk.Native.OntId.RegIDWithPublicKey(
			2500, 20000, seller, "did:ont:"+buyer.Address.ToBase58(), buyer)
		fmt.Println(err)
		time.Sleep(10 * time.Second)
		evt, err := ontSdk.GetSmartContractEvent(txhash.ToHexString())
		fmt.Println(err)
		fmt.Println(evt)
		return
	}

	if false {
		//aa5f1094a999bc70ccc23a1f41edf3473960492c
		utils.DeployMpContract(sdk, seller, gasPrice)
		return
	}
	if false {
		//4f1e90da08d410be2e49a1ac294e9c6f6d3e967e
		//3896c83139e24ac9f5eb6d992521e6d4691404bb
		utils.DeployDTokenContract(sdk, seller, gasPrice)
		return
	}
	if false {
		//f4d41829ff666d86dd97945679e19fa7fbb324a7
		//dc7560f3e456677d3d5b4663ef250545aaf91749
		utils.DeploySplitPolicyContract(sdk, seller, gasPrice)
		return
	}
	if false {
		//5cfd487b6546b5f312ed1f341d6f5a40f7f17a60
		//0cb14bd01c2f76bdf40f6788e9f8fc4c3bf2debf
		utils.DeployDataIdContract(sdk, seller, gasPrice)
		return
	}
	if false {
		//e8497665367c5931e0319ae0aaf846dd82012e1a
		utils.DeployAccountantContract(sdk, seller, gasPrice)
		return
	}
	if false {
		//ae28201fd1ae7399a46a821bb6e78664127eb772
		//b57aae109f66dea7168c6f3443edca542c763167
		//f506d822135c492d625d81a2c052850b7a07d1ec
		//aa37941950f76bf7d15a06bd1f07dcdd0cc6e34c
		utils.DeployOpenkgContract(sdk, seller, gasPrice)
		return
	}

	mpContractAddr, _ := common.AddressFromHexString("e01d500ed0c1719b7750367ae59b4b2d308d1ceb")
	dtokenContractAddr, _ := common.AddressFromHexString("466b94488bf2ad1b1eec0ae7e49e40708e71a35d")
	splitContractAddr, _ := common.AddressFromHexString("dc7560f3e456677d3d5b4663ef250545aaf91749")
	dataIdContractAddr, _ := common.AddressFromHexString("5cfd487b6546b5f312ed1f341d6f5a40f7f17a60")
	accountantContractAddr, _ := common.AddressFromHexString("e8497665367c5931e0319ae0aaf846dd82012e1a")
	openkgContractAddr, _ := common.AddressFromHexString("aa37941950f76bf7d15a06bd1f07dcdd0cc6e34c")

	if false {
		contractAddr, _ := common.AddressFromHexString("")
		utils.DataIdTest(sdk, pwd, seller, contractAddr)
		return
	}

	fmt.Printf("mpContractAddr: %s,dtokenContractAddr:%s, \n", mpContractAddr.ToHexString(),
		dtokenContractAddr.ToHexString())
	fmt.Printf("splitContractAddr:%s,dataIdContractAddr:%s \n", splitContractAddr.ToHexString(),
		dataIdContractAddr.ToHexString())
	fmt.Printf("accountantContractAddr:%s,openkgContractAddr:%s \n", accountantContractAddr.ToHexString(),
		openkgContractAddr.ToHexString())
	fmt.Println("====================")

	if false {
		sdk.SetMpContractAddress(mpContractAddr)
		sdk.SetGasPrice(2500)
		txHash, err := sdk.DefMpKit().Init(common.ADDRESS_EMPTY, seller, dtokenContractAddr, splitContractAddr)
		if err != nil {
			fmt.Println("Init failed: ", err)
			return
		}
		showNotify(sdk, "init", txHash.ToHexString())
		return
	}
	if false {
		contractAddr, _ := common.AddressFromHexString("")
		con := sdk.DefContract(contractAddr)
		resourceId := []byte("5945133018807634517")
		res, err := con.PreInvoke("getRegisterParam", []interface{}{resourceId})
		if err != nil {
			fmt.Println(err)
			return
		}
		data, err := res.ToByteArray()
		if err != nil {
			fmt.Println(err)
			return
		}
		param := &split_policy_contract.SplitPolicyRegisterParam{}
		err = param.FromBytes(data)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(param)
		return
	}
	//openkg
	if false {
		sdk.SetGasPrice(gasPrice)
		con := sdk.DefContract(openkgContractAddr)
		if true {
			utils.Init(sdk, con, seller, mpContractAddr, dtokenContractAddr)
			return
		}
		if false {
			utils.GetMp(con)
			return
		}
		if false {
			utils.GetDToken(con)
			return
		}
		if false {
			resourceId := []byte("5945133018807634517")
			templateId, _ := hex.DecodeString("30")
			utils.BuyRewardAndUseToken(sdk, con, resourceId, 1, 10000000000, buyer, payer, templateId)
			return
		}
		if false {
			dtoken, _ := common.AddressFromHexString("3343753265152550e5a1741cea946436744ab442")
			utils.SetDtokenContractAddr(sdk, con, seller, dtoken)
			return
		}
		if false {
			mp, _ := common.AddressFromHexString("5fbcadf08b14aa737de8af429483dc4fb1ae13d3")
			utils.SetMpContractAddr(sdk, con, seller, mp)
			return
		}
		if false {
			resource_id := []byte("637088084609811033")
			template_id, _ := hex.DecodeString("30")
			utils.BuyAndUseToken(sdk, con, resource_id, 1, buyer, payer, template_id)
			return
		}
		//old
		if false {
			resourceId := []byte("5989611139890472506")
			templateBytes, _ := hex.DecodeString("012a6469643a6f6e743a544e735351344374366f474c684c61637465586a7a454d3957736646376a47614c4c0101310461616161")
			txHash, err := con.Invoke("buyRewardAndUseToken", buyer, []interface{}{resourceId, 2, buyer.Address, buyer.Address, templateBytes, 10})
			if err != nil {
				fmt.Println(err)
				return
			}
			showNotify(sdk, "buyAndUseToken", txHash.ToHexString())
			return
		}
		if true {
			resourceId := []byte("5989611139890472506")
			templateBytes, _ := hex.DecodeString("012a6469643a6f6e743a544e735351344374366f474c684c61637465586a7a454d3957736646376a47614c4c0101310461616161")
			txHash, err := con.Invoke("buyAndUseToken", buyer, []interface{}{resourceId, 2, buyer.Address, buyer.Address, templateBytes})
			if err != nil {
				fmt.Println(err)
				return
			}
			showNotify(sdk, "buyAndUseToken", txHash.ToHexString())
			return
		}
	}
	//DToken
	if true {
		dtokenContractAddr, _ = common.AddressFromHexString("3896c83139e24ac9f5eb6d992521e6d4691404bb")
		sdk.DefDTokenKit().SetContractAddr(dtokenContractAddr)
		if false {
			utils.CreateTokenTemplate(sdk, seller)
			return
		}
		if false {
			agent, _ := common.AddressFromBase58("ALBoDa6bkA3WPPSTSaGogRuxNYgHPLveee")
			tokenId, _ := hex.DecodeString("3632")
			con := sdk.DefContract(dtokenContractAddr)
			agentBa, err := con.PreInvoke("getAgentBalance", []interface{}{agent, tokenId})
			fmt.Println(err)
			bs, err := agentBa.ToByteArray()
			fmt.Println(err)
			bbb := common.BigIntFromNeoBytes(bs)
			fmt.Println("bbb:", bbb.Uint64())

			owner, _ := common.AddressFromBase58("AXP12gzwoem3Rko82SBpYdG5uFpukpveEf")
			ba, err := sdk.DefDTokenKit().BalanceOf(common.ADDRESS_EMPTY, owner, tokenId)
			fmt.Println(err, ba)
			return
		}

		if false {
			con := sdk.DefContract(dtokenContractAddr)
			res, err := con.Invoke("setDdxfContract", seller, []interface{}{mpContractAddr})
			if err != nil {
				fmt.Println(err)
				return
			}
			showNotify(sdk, "setDdxfContract", res.ToHexString())
			return
		}
		if false {
			con := sdk.DefContract(dtokenContractAddr)
			res, err := con.PreInvoke("getDdxfContract", []interface{}{})
			if err != nil {
				fmt.Println(err)
				return
			}
			bs, _ := res.ToByteArray()
			addr, _ := common.AddressParseFromBytes(bs)
			fmt.Println(addr.ToHexString())
			return
		}
		if false {
			templateId, _ := hex.DecodeString("30")
			utils.AuthorizeTokenTemplate(sdk, templateId, seller, []common.Address{buyer.Address, admin.Address})
			return
		}
		if false {
			utils.GenerateDtoken(sdk, seller)
			return
		}
		tokenId, _ := hex.DecodeString("3632")
		if false {
			utils.BalanceOf(sdk, buyer.Address, tokenId)
			return
		}
		if false {
			if err = addAgents(sdk, tokenId); err != nil {
				fmt.Println("addAgents error: ", err)
				return
			}
		}

		if true {
			pri, _ := hex.DecodeString("24cc3b4cf89b669b84e4ca0d6fde48200e8c0e8d72fa9bb342ebc8ab4b2d92a0")
			agent, _ = ontology_go_sdk.NewAccountFromPrivateKey(pri, signature.SHA256withECDSA)
			owner, _ := common.AddressFromBase58("AXP12gzwoem3Rko82SBpYdG5uFpukpveEf")
			if err = useTokenByAgent(sdk, tokenId, owner); err != nil {
				fmt.Println("useTokenByAgent error: ", err)
				return
			}
			return
		}

		if false {
			if err = removeAgents(sdk, tokenId); err != nil {
				fmt.Println("removeAgents error: ", err)
				return
			}
		}

		if false {
			if err = setTokenAgents(sdk, tokenId); err != nil {
				fmt.Println("setTokenAgents error: ", err)
				return
			}
			return
		}

		if err = addTokenAgents(sdk, tokenId); err != nil {
			fmt.Println("addTokenAgents error: ", err)
			return
		}
		if err = removeTokenAgents(sdk); err != nil {
			fmt.Println("removeTokenAgents error: ", err)
			return
		}

		err = useToken(sdk, tokenId)
		if err != nil {
			fmt.Printf("useToken: %s\n", err)
			return
		}
		return
	}

	if true {
		sdk.DefMpKit().SetContractAddress(mpContractAddr)
		resourceId := strconv.Itoa(rand.Int())
		fmt.Println("resourceId:", resourceId)
		resourceIdBytes := []byte(resourceId)
		dataId := ""
		tokenTemplate = &market_place_contract.TokenTemplate{
			DataID:     dataId,
			TokenHashs: []string{string(common.UINT256_EMPTY[:])},
		}
		//old
		if true {
			mp := sdk.DefContract(mpContractAddr)
			ddoBytes, _ := hex.DecodeString("fbe02b027e61a6d7602f26cfa9487fa58ef9ee72804669a4f71e49a598f34bf7b949676a803b6e9bc0bc483c160a6c1d80cbcd1c000000")
			itemBytes, _ := hex.DecodeString("1ffc48df93cde46d3f80523a24e82567da725d190200000000000000008cfb486669350900102700000000000001012a6469643a6f6e743a544e735351344374366f474c684c61637465586a7a454d3957736646376a47614c4c0101310461616161")
			sp := split_policy_contract.SplitPolicyRegisterParam{
				AddrAmts: []*split_policy_contract.AddrAmt{
					&split_policy_contract.AddrAmt{
						To:          seller.Address,
						Percent:     5000,
						HasWithdraw: false,
					},
					&split_policy_contract.AddrAmt{
						To:          admin.Address,
						Percent:     5000,
						HasWithdraw: false,
					},
				},
				TokenTy: split_policy_contract.ONG,
			}

			txHash, err := mp.Invoke("dtokenSellerPublish", seller,
				[]interface{}{resourceIdBytes, ddoBytes, itemBytes, sp.ToBytes()})
			if err != nil {
				fmt.Println(err)
				return
			}
			showNotify(sdk, "dtokenSellerPublish", txHash.ToHexString())
			return
		}

		if err = publish(sdk, resourceIdBytes); err != nil {
			fmt.Println("publish error: ", err)
			return
		}

		//if err := delete(sdk, resourceIdBytes); err != nil {
		//	fmt.Println("delete error: ", err)
		//	return
		//}
		//return

		//if err := update(sdk, resourceIdBytes); err != nil {
		//	fmt.Println("update error: ", err)
		//	return
		//}

		//if err = buyAndUseToken(sdk, resourceIdBytes); err != nil {
		//	fmt.Println("buyAndUseToken error: ", err)
		//	return
		//}

		//if err = buyDToken(sdk, resourceIdBytes); err != nil {
		//	fmt.Println("buyDToken error: ", err)
		//	return
		//}
	}
}

func addTokenAgents(sdk *ddxf_sdk.DdxfSdk, tokenId []byte) error {
	txHash, err := sdk.DefDTokenKit().AddTokenAgents(common.ADDRESS_EMPTY, buyer,
		[]common.Address{agent.Address}, tokenId, []int{1})
	if err != nil {
		return err
	}
	return showNotify(sdk, "addTokenAgents", txHash.ToHexString())
}
func setTokenAgents(sdk *ddxf_sdk.DdxfSdk, tokenId []byte) error {
	txHash, err := sdk.DefDTokenKit().SetTokenAgents(common.ADDRESS_EMPTY, seller,
		[]common.Address{agent.Address}, tokenId, []int{1})
	if err != nil {
		return err
	}
	return showNotify(sdk, "setTokenAgents", txHash.ToHexString())
}

func removeTokenAgents(sdk *ddxf_sdk.DdxfSdk) error {
	txHash, err := sdk.DefDTokenKit().RemoveTokenAgents(common.ADDRESS_EMPTY, []byte(""), buyer,
		[]common.Address{agent.Address})
	if err != nil {
		return err
	}
	return showNotify(sdk, "removeTokenAgents", txHash.ToHexString())
}

func removeAgents(sdk *ddxf_sdk.DdxfSdk, tokenId []byte) error {
	txHash, err := sdk.DefDTokenKit().RemoveAgents(common.ADDRESS_EMPTY, buyer, []common.Address{agent.Address}, [][]byte{tokenId})
	if err != nil {
		return err
	}
	return showNotify(sdk, "removeAgents", txHash.ToHexString())
}

func useTokenByAgent(sdk *ddxf_sdk.DdxfSdk, tokenId []byte, owner common.Address) error {

	txHash, err := sdk.DefDTokenKit().UseTokenByAgents(common.ADDRESS_EMPTY, owner, agent, tokenId, 1)
	if err != nil {
		return err
	}
	return showNotify(sdk, "useTokenByAgent", txHash.ToHexString())
}

func addAgents(sdk *ddxf_sdk.DdxfSdk, tokenId []byte) error {
	txHash, err := sdk.DefDTokenKit().AddAgents(common.ADDRESS_EMPTY, buyer,
		[]common.Address{agent.Address}, []int{1}, [][]byte{tokenId})
	if err != nil {
		fmt.Printf("AddAgents: %s\n", err)
		return err
	}
	return showNotify(sdk, "addAgents", txHash.ToHexString())
}

func useToken(sdk *ddxf_sdk.DdxfSdk, tokenId []byte) error {
	txHash, err := sdk.DefDTokenKit().UseToken(common.ADDRESS_EMPTY, seller, tokenId, 1)
	if err != nil {
		return err
	}
	return showNotify(sdk, "useToken", txHash.ToHexString())
}

//func buyAndUseToken(sdk *ddxf_sdk.DdxfSdk, resourceIdBytes []byte) error {
//	txHash, err := sdk.DefMpKit().BuyAndUseToken(buyer, payer, resourceIdBytes, 2, *tokenTemplate)
//	if err != nil {
//		return err
//	}
//	return showNotify(sdk, "buyAndUseToken", txHash.ToHexString())
//}

func buyDToken(sdk *ddxf_sdk.DdxfSdk, resourceIdBytes []byte) error {
	txHash, err := sdk.DefMpKit().BuyDToken(common.ADDRESS_EMPTY, buyer, payer, resourceIdBytes, 10)
	if err != nil {
		return err
	}
	return showNotify(sdk, "buyDToken", txHash.ToHexString())
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

func delete(sdk *ddxf_sdk.DdxfSdk, resourceIdBytes []byte) error {
	txHash, err := sdk.DefMpKit().Delete(common.ADDRESS_EMPTY, seller, resourceIdBytes)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return showNotify(sdk, "delete", txHash.ToHexString())
}
func update(sdk *ddxf_sdk.DdxfSdk, resourceIdBytes []byte) error {
	itemMeta := map[string]interface{}{
		"key": "value",
	}
	bs, err := ddxf.HashObject(itemMeta)
	if err != nil {
		return err
	}
	itemMetaHash, err := common.Uint256ParseFromBytes(bs[:])

	ddo := market_place_contract.ResourceDDO{
		Manager:      seller.Address,       // data owner id
		ItemMetaHash: itemMetaHash,         // required if len(Templates) > 1
		DTC:          []common.Address{},   // can be empty
		Accountant:   common.ADDRESS_EMPTY, // can be empty
		Split:        common.ADDRESS_EMPTY,
	}

	tokenTemplateId, _ := hex.DecodeString("30")

	item := market_place_contract.DTokenItem{
		Fee: market_place_contract.Fee{
			ContractType: 0,
			Count:        1,
		},
		ExpiredDate:      uint64(time.Now().Unix()) + 10000,
		Stocks:           10000,
		TokenTemplateIds: []string{string(tokenTemplateId)},
	}

	sp := split_policy_contract.SplitPolicyRegisterParam{
		AddrAmts: []*split_policy_contract.AddrAmt{
			&split_policy_contract.AddrAmt{
				To:          seller.Address,
				Percent:     5000,
				HasWithdraw: false,
			},
			&split_policy_contract.AddrAmt{
				To:          admin.Address,
				Percent:     5000,
				HasWithdraw: false,
			},
		},
		TokenTy: split_policy_contract.ONG,
	}

	txHash, err := sdk.DefMpKit().Update(common.ADDRESS_EMPTY, seller, resourceIdBytes, ddo, item, sp)
	if err != nil {
		fmt.Printf("update error:%s\n", err)
		return err
	}
	fmt.Println("update txHash: ", txHash.ToHexString())

	return showNotify(sdk, "update", txHash.ToHexString())
}

func publish(sdk *ddxf_sdk.DdxfSdk, resourceIdBytes []byte) error {
	itemMeta := map[string]interface{}{
		"key": "value",
	}
	bs, err := ddxf.HashObject(itemMeta)
	if err != nil {
		return err
	}
	itemMetaHash, err := common.Uint256ParseFromBytes(bs[:])

	ddo := market_place_contract.ResourceDDO{
		Manager:      seller.Address,       // data owner id
		ItemMetaHash: itemMetaHash,         // required if len(Templates) > 1
		DTC:          []common.Address{},   // can be empty
		Accountant:   common.ADDRESS_EMPTY, // can be empty
		Split:        common.ADDRESS_EMPTY,
	}

	tokenTemplateId, _ := hex.DecodeString("30")
	item := market_place_contract.DTokenItem{
		Fee: market_place_contract.Fee{
			ContractType: 1,
			Count:        0,
		},
		ExpiredDate:      uint64(time.Now().Unix()) + 10000,
		Stocks:           10000,
		TokenTemplateIds: []string{string(tokenTemplateId)},
	}

	sp := split_policy_contract.SplitPolicyRegisterParam{
		AddrAmts: []*split_policy_contract.AddrAmt{
			&split_policy_contract.AddrAmt{
				To:          seller.Address,
				Percent:     5000,
				HasWithdraw: false,
			},
			&split_policy_contract.AddrAmt{
				To:          admin.Address,
				Percent:     5000,
				HasWithdraw: false,
			},
		},
		TokenTy: split_policy_contract.ONG,
	}

	txHash, err := sdk.DefMpKit().Publish(common.ADDRESS_EMPTY, seller, resourceIdBytes, ddo, item, sp)
	if err != nil {
		fmt.Printf("Publish error:%s\n", err)
		return err
	}
	fmt.Println("publish txHash: ", txHash.ToHexString())
	return showNotify(sdk, "publish", txHash.ToHexString())
}
