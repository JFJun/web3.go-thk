package test

import (
	"fmt"
	"time"

	"math/big"
	"strconv"
	"testing"
	"web3.go/common/cryp/crypto"
	"web3.go/common/encoding"
	"web3.go/common/hexutil"
	"web3.go/web3"
	"web3.go/web3/providers"
	"web3.go/web3/thk/util"
)

var (
	key1 = "0811c0ad9ce2effe2d8e09deef3c472ec68348ee79dd154d514e4aeda26df5d9"
)

func TestThkGetStats(t *testing.T) {
	var connection = web3.NewWeb3(providers.NewHTTPProvider("rpctest.thinkey.xyz", 10, false))
	stats, err := connection.Thk.GetStats("2")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	println(stats.Currentheight)
}
func TestVccCash(t *testing.T) {
	var err error
	var connection = web3.NewWeb3(providers.NewHTTPProvider("rpctest.thinkey.xyz", 10, false))
	from := "0x151b3a46fb5c1ecffd8feccb975acab63bf52652"
	to := "0x0000000000000000000000000000000000020000"
	toAddress := "0x1111111111111111111111111111111111111112"
	value := "3" + "000000000000000000"

	stats, err := connection.Thk.GetStats("2")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	println(stats.Currentheight)

	expireHeight := Height(stats.Currentheight) + 200

	nonce, err := connection.Thk.GetNonce(from, "1")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	from_str, err := hexutil.Decode(from)
	to_str, err := hexutil.Decode(toAddress)
	bigValue, _ := big.NewInt(0).SetString(value, 10)
	vcc := &CashCheck{
		FromChain:    1,
		FromAddress:  BytesToAddress(from_str),
		Nonce:        uint64(nonce),
		ToChain:      2,
		ToAddress:    BytesToAddress(to_str),
		ExpireHeight: expireHeight,
		Amount:       bigValue,
	}
	println(vcc.Nonce)
	intput, err := encoding.Marshal(vcc)
	println(intput)

	str := hexutil.Encode(intput)
	fmt.Println("------------------")
	fmt.Println(str)
	transaction := util.Transaction{
		ChainId: "1", FromChainId: "1", ToChainId: "2", From: from,
		To: to, Value: "0", Input: str, Nonce: strconv.Itoa(int(nonce)),
	}

	privatekey, err := crypto.HexToECDSA(key1)
	err = connection.Thk.SignTransaction(&transaction, privatekey)

	txhash, err := connection.Thk.SendTx(&transaction)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	fmt.Println("txhash:", txhash)

	time.Sleep(10 * time.Second)

	res, err := connection.Thk.GetTransactionByHash("1", txhash)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	fmt.Println("res:", res)


	// 0x472a80cd5a8aa4664fcca5f3a4fd72c3ff25681c2511325f4613f04c128966e9

	// nonce, err = connection.Thk.GetNonce(from, "1")
	// if err != nil {
	// 	t.Error(err)
	// 	t.FailNow()
	// }

	// stats, _ := connection.Thk.GetStats(2)

	// fmt.Println(stats.Currentheight)

	fmt.Println(expireHeight)

	transaction = util.Transaction{
		ChainId: "1", FromChainId: "1", ToChainId: "2", From: from,
		To: toAddress, Nonce: strconv.Itoa(int(nonce)), Value: value, ExpireHeight: int64(expireHeight),
	}
	input, err := connection.Thk.RpcMakeVccProof(&transaction)
	t.Log("input:", input)


	// save
	to = "0x0000000000000000000000000000000000030000"
	nonce, err = connection.Thk.GetNonce(from, "2")
	fmt.Println(nonce)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	//
	// from_str,err:=hexutil.Decode("0x2c7536e3605d9c16a7a3d7b1898e529396a65c23")
	// to_str,err:=hexutil.Decode("0x4fa1c4e6182b6b7f3bca273390cf587b50b47311")
	// vcc := &CashCheck{
	//	FromChain:    2,
	//	FromAddress:  BytesToAddress(from_str),
	//	Nonce:        uint64(nonce),
	//	ToChain:      3,
	//	ToAddress:    BytesToAddress(to_str),
	//	ExpireHeight: 33772,
	//	Amount:       big.NewInt(1),
	// }
	// println(vcc.Nonce)
	// intput,err:=encoding.Marshal(vcc)
	// println(intput)
	//
	// str:=hexutil.Encode(intput)
	transaction = util.Transaction{
		ChainId: "2", FromChainId: "2", ToChainId: "2", From: from,
		To: to, Value: "0", Input: input, Nonce: strconv.Itoa(int(nonce)),
	}

	err = connection.Thk.SignTransaction(&transaction, privatekey)

	txhash, err = connection.Thk.SendTx(&transaction)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log("txhash:", txhash)

	time.Sleep(30 * time.Second)

	res, err = connection.Thk.GetTransactionByHash("2", txhash)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log("res:", res)


}
func TestVCCThkSaveCashCheck(t *testing.T) {
	var err error
	var connection = web3.NewWeb3(providers.NewHTTPProvider("test.thinkey.xyz", 10, false))
	from := "0x2c7536e3605d9c16a7a3d7b1898e529396a65c23"
	to := "0x0000000000000000000000000000000000030000"

	nonce, err := connection.Thk.GetNonce(from, "2")
	fmt.Println(nonce)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	//
	// from_str,err:=hexutil.Decode("0x2c7536e3605d9c16a7a3d7b1898e529396a65c23")
	// to_str,err:=hexutil.Decode("0x4fa1c4e6182b6b7f3bca273390cf587b50b47311")
	// vcc := &CashCheck{
	//	FromChain:    2,
	//	FromAddress:  BytesToAddress(from_str),
	//	Nonce:        uint64(nonce),
	//	ToChain:      3,
	//	ToAddress:    BytesToAddress(to_str),
	//	ExpireHeight: 33772,
	//	Amount:       big.NewInt(1),
	// }
	// println(vcc.Nonce)
	// intput,err:=encoding.Marshal(vcc)
	// println(intput)
	//
	// str:=hexutil.Encode(intput)
	transaction := util.Transaction{
		ChainId: "2", FromChainId: "2", ToChainId: "2", From: from,
		To: to, Value: "0", Input: "0x95000000022c7536e3605d9c16a7a3d7b1898e529396a65c230000000000010009000000034fa1c4e6182b6b7f3bca273390cf587b50b473110000000000045644010102a301ba1fc09cd5f9d10c23f8e2db49d4d4e529a32b5b951e3685f314eda7f6d13289dc6aa894941093a1a0df6cfaa2c89bf9deeed6a9c03667d40ca358b2adc9e091d2598a2b7e7220a000c200008080940e934080c2084080810001d0a2ea876f373a05d990e1c46041af438ff0e25e7d5d6953dcc9c43e2845026f0001019403934080c27aa2808100038187aa9f339cf1ba6ffe6986f68c639a835fac453ac37d0df6e72091b1cd1cd3d42acb443bbd30466cf2f099f5fc277f9beb032a09f8b074201404d94cb21947ade490581abc936a49b4754aaac0816195d4af0d77a6fd454210762d8da590180001019424930080c20000c0b514b73aa5d9299ebaa524822220c50a1c884bcd6e1193c279b4b2023e4fc5c181000509f47f9feafa18ad06f468d253c4d9aa5bebe0438fe01a00a830f0546d5d60b8625dc71f6529f508c2f6411029909f5207b556920cf45d64951b1781a9e8b17431f3959a8327f5d093bc5fae377a4a831f70d74bccf65eab93cfde3d2a8fab34eca078605c1b0ad6ff4323f7c23307585d3dddd504f96e7a7f722f9802d2a1b787f28d0a0b5499f8c6dc7afdcb43e1feddb8e21beb4750c81c947f0aed109090000110", Nonce: strconv.Itoa(int(nonce)),
	}
	privatekey, err := crypto.HexToECDSA(key)
	err = connection.Thk.SignTransaction(&transaction, privatekey)

	txhash, err := connection.Thk.SendTx(&transaction)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log("txhash:", txhash)
	// 0x920a95dc3af9d6ed801258fc8eeb1455b7e6b35a72d4142c995a27f4f0e78c8d
}

// Chainid:
// From:
// To:
// Nonce:   nonce,
// ExpireHeight: expireheight,
// Amount: value.(string),
func TestVCCThkRpcMakeVccProof(t *testing.T) {
	var err error
	var connection = web3.NewWeb3(providers.NewHTTPProvider("rpctest.thinkey.xyz", 10, false))
	from := "0x2c7536e3605d9c16a7a3d7b1898e529396a65c23"
	to := "0x0000000000000000000000000000000000020000"

	nonce, err := connection.Thk.GetNonce(from, "2")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// stats, _ := connection.Thk.GetStats(2)
	expireHeight := 279228 + 5000

	// fmt.Println(stats.Currentheight)

	fmt.Println(expireHeight)

	transaction := util.Transaction{
		ChainId: "2", FromChainId: "2", ToChainId: "3", From: from,
		To: to, Nonce: strconv.Itoa(int(nonce)), Value: "2333", ExpireHeight: int64(expireHeight),
	}
	input, err := connection.Thk.RpcMakeVccProof(&transaction)
	t.Log("input:", input)
}

func TestVCCThkMakeCCCExistenceProof(t *testing.T) {
	var err error
	var connection = web3.NewWeb3(providers.NewHTTPProvider("rpctest.thinkey.xyz", 10, false))
	from := "0x2c7536e3605d9c16a7a3d7b1898e529396a65c23"
	to := "0x0000000000000000000000000000000000020000"

	nonce, err := connection.Thk.GetNonce(from, "3")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// stats, _ := connection.Thk.GetStats(3)
	expireHeight := 279228 + 5000

	transaction := util.Transaction{
		ChainId: "2", FromChainId: "2", ToChainId: "2", From: from,
		To: to, Nonce: strconv.Itoa(int(nonce)), Value: "2333", ExpireHeight: int64(expireHeight),
	}
	input, err := connection.Thk.MakeCCCExistenceProof(&transaction)
	t.Log("input:", input)
}
