package test

import (
	"fmt"
	"github.com/JFJun/web3.go-thk/web3/providers"
	"github.com/JFJun/web3.go-thk/web3/thk"
	"testing"
)

func Test_Rpc(t *testing.T) {
	c := thk.NewThk(providers.NewHTTPProvider("https://rpcproxy.thinkium.vip", 60, true))
	resp, err := c.GetStats("1")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(resp)
	fmt.Println(resp.Currentheight)
}
