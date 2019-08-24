package test

import (
	"encoding/json"
	"io/ioutil"
	"math/big"
	"strconv"
	"testing"
	"time"
	"wallet/models"
	"web3.go/common/cryp/crypto"
	"web3.go/common/hexutil"
	"web3.go/web3"
	"web3.go/web3/providers"
	"web3.go/web3/thk/util"
)

func TestRelease(t *testing.T) {
	from := "0x2c7536e3605d9c16a7a3d7b1898e529396a65c23"
	contentTokenVesting, err := ioutil.ReadFile(TokenVestingJsonName)
	if err != nil {
		println(err)

	}
	var TokenVestingResponse TruffleContract
	_ = json.Unmarshal(contentTokenVesting, &TokenVestingResponse)

	var connection = web3.NewWeb3(providers.NewHTTPProvider("192.168.1.13:8089", 10, false))

	//bytecodeTokenVesting := TokenVestingResponse.Bytecode
	contractTokenVesting, err := connection.Thk.NewContract(TokenVestingResponse.Abi)
	privatekey, err := crypto.HexToECDSA(key)
	hashTokenVesting := "0xe0c77573eb6e3447dcd7c361a9c6ef91d51e5376520664d2f6aad3aed04a2ae6"
	receiptTokenVesting, err := connection.Thk.GetTransactionByHash(Chain, hashTokenVesting)
	toTokenVesting := receiptTokenVesting.ContractAddress

	//Vesting
	nonce, err := connection.Thk.GetNonce(from, Chain)
	if err != nil {
		println("get nonce error")
		return
	}
	transactionTokenVesting := util.Transaction{
		ChainId: Chain, FromChainId: Chain, ToChainId: Chain, From: from,
		To: toTokenVesting, Value: "0", Input: "", Nonce: strconv.Itoa(int(nonce)),
	}
	tmcliff, errc := time.Parse("2006-01-02 15:04:05", "2019-07-19 17:47:00")
	tmstart, errc:= time.Parse("2006-01-02 15:04:05", "2019-07-19 17:48:00")
	tmend, errc := time.Parse("2006-01-02 15:04:05", "2019-07-19 17:49:00")
	println(errc)
	toAddress, err := hexutil.Decode("0x14723a09acff6d2a60dcdf7aa4aff308fddc160c")
	vestAddress := models.BytesToAddress(toAddress)
	//vestAddressP,boolerr := new(big.Int).SetString(strings.TrimPrefix(string("0x14723a09acff6d2a60dcdf7aa4aff308fddc160c"), "0x"),16)
	//println(boolerr)
	cliff:=tmcliff.Unix()
	start:=tmstart.Unix()
	end:=tmend.Unix()
	tmcliffP:= new(big.Int).SetInt64(cliff)
	tmstartP:= new(big.Int).SetInt64(start)
	tmendP:= new(big.Int).SetInt64(end)
	timesP:= new(big.Int).SetInt64(2)
	total := new(big.Int).SetUint64(uint64(100))
	resultTokenVesting, err := contractTokenVesting.Send(transactionTokenVesting, "addPlan",privatekey,
		vestAddress, tmcliffP, tmstartP, timesP, tmendP, total, false, "上交易所后私募锁仓10分钟，之后每10分钟释放50%")
	if err != nil {
		t.Error(err)
	}
	t.Log("result:", resultTokenVesting)




	nonce, err = connection.Thk.GetNonce(from, Chain)
	if err != nil {
		println(err)
		println("get nonce error")
		return
	}
	transactionReleasable := util.Transaction{
		ChainId: Chain, FromChainId: Chain, ToChainId: Chain, From: from,
		To: toTokenVesting, Value: "0", Input: "", Nonce: strconv.Itoa(int(nonce)),
	}
	value := new(big.Int).SetUint64(uint64(0))
	result, err := contractTokenVesting.Call(transactionReleasable, "releasableAmount",vestAddress)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log("result:", result)
	err = contractTokenVesting.Parse(result, "releasableAmount", &value)
	if err != nil {
		println("failed")
		return
	}


	transactionRelease := util.Transaction{
		ChainId: Chain, FromChainId: Chain, ToChainId: Chain, From: from,
		To: toTokenVesting, Value: "0", Input: "", Nonce: strconv.Itoa(int(nonce)),
	}

	hashTokenVesting, err = contractTokenVesting.Send(transactionRelease, "release", privatekey, vestAddress)
	if err != nil {
		t.Error(err)
	}
	t.Log("result:", hashTokenVesting)

}
