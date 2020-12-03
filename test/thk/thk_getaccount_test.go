package test

import (
	"github.com/JFJun/web3.go-thk/web3"
	"github.com/JFJun/web3.go-thk/web3/providers"
	"testing"
)

func TestThkGetBalance(t *testing.T) {
	var connection = web3.NewWeb3(providers.NewHTTPProvider("rpctest.thinkey.xyz", 10, false))
	connection.DefaultAddress = "0x6422ce49bd62dba229568eb6ad868d15c49d7fe8"
	bal, err := connection.Thk.GetBalance(connection.DefaultAddress, "1")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log("Balance:", bal)
}

func TestThkGetNonce(t *testing.T) {
	var connection = web3.NewWeb3(providers.NewHTTPProvider("rpctest.thinkey.xyz", 10, false))
	connection.DefaultAddress = "0x6422ce49bd62dba229568eb6ad868d15c49d7fe8"
	nonce, err := connection.Thk.GetNonce(connection.DefaultAddress, "1")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log("nonce:", nonce)
}
