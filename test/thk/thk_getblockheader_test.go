package test

import (
	"github.com/JFJun/web3.go-thk/web3"
	"github.com/JFJun/web3.go-thk/web3/providers"
	"testing"
)

func TestThkGetBlockHeader(t *testing.T) {
	var connection = web3.NewWeb3(providers.NewHTTPProvider("rpctest.thinkey.xyz", 10, false))
	res, err := connection.Thk.GetBlockHeader("1", "30")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log("BlockHeader:", res)
}

func TestThkGetBlockTxs(t *testing.T) {
	// var connection = web3.NewWeb3(providers.NewHTTPProvider("rpctest.thinkey.xyz", 10, false))
	// res, err := connection.Thk.GetBlockTxs("1", "30","1","10")
	// if err != nil {
	// 	t.Error(err)
	// 	t.FailNow()
	// }
	// t.Log("BlockHeader:", res)
}
