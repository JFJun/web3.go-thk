package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Alex-Chris/log/log"
	"io"
	"io/ioutil"
	"math/big"
	"mime/multipart"
	"net/http"
	"path"
	"strconv"
	"strings"
	"testing"
	"time"
	"web3.go/common/cryp/crypto"
	"web3.go/web3"
	"web3.go/web3/providers"
	"web3.go/web3/thk/util"
)

type ContractJson struct {
	ContractName string `json:"contractName"`
	ABI          string `json:"abi"`
	ByteCode     string `json:"bytecode"`
}

var jsonName = "../resources/ERC20.json"

var Symbol = "GBDJ"
var Name = "GOLDDJ"
var Des = "金币"
var Chain = "2"
var Url = "test.thinkey.xyz"

type TruffleContract struct {
	Abi      string `json:"abi"`
	Bytecode string `json:"bytecode"`
}
type Token struct {
	Name            string  `json:"name"`
	Symbol          string  `json:"symbol"`
	Total           float64 `json:"total"`
	ContractAddress string  `json:"contractaddress"`
	ABI             string  `json:"abi"`
	Icon            string  `json:"icon"`
	Website         string  `json:"website"`
	Introduction    string  `json:"introduction"`
	State           string  `json:"state"`
	Date            string  `json:"date"`
	ChainId         string  `json:"chainid"`
	Decimal         int64   `json:"decimal"`
}

func TestDeploy(t *testing.T) {
	ctct, err := CompileContract("../resources/contract/ERC20.sol",
		"../resources/contract/IERC20.sol", "../resources/contract/Pausable.sol",
		"../resources/contract/SafeMath.sol", "../resources/contract/Ownable.sol",
		 )
	var ERC20Json ContractJson
	amount := new(big.Int).SetUint64(uint64(100000))
	decimal := uint8(8)
	//for keyname, value := range ctct {
	//	contractJson.ContractName = keyname
	//	mapvalue = value.(map[string]interface{})
	//	contractJson.ByteCode = mapvalue["code"].(string)
	//	info = mapvalue["info"].(map[string]interface{})
	//	abidef := info["abiDefinition"]
	//	abibytes, err := json.Marshal(abidef)
	//	contractJson.ABI = string(abibytes)
	//	if err != nil {
	//		return
	//	}
	//}
	ERC20Json, _, err = GetContractJson(ctct)
	fmt.Println("contractJson的值\n", ctct)
	data, err := json.MarshalIndent(ERC20Json, "", "  ")
	if ioutil.WriteFile(jsonName, data, 0644) == nil {
		fmt.Println("写入文件成功")
	}
	content, err := ioutil.ReadFile(jsonName)
	if err != nil {
		log.Error(err)
		panic("get file error")
	}

	var unmarshalResponse TruffleContract
	err = json.Unmarshal(content, &unmarshalResponse)
	if err != nil {
		log.Error(err)
	}

	var connection = web3.NewWeb3(providers.NewHTTPProvider(Url, 10, false))

	bytecode := unmarshalResponse.Bytecode
	contract, err := connection.Thk.NewContract(unmarshalResponse.Abi)
	from := "0x2c7536e3605d9c16a7a3d7b1898e529396a65c23"
	nonce, err := connection.Thk.GetNonce(from, Chain)
	if err != nil {
		log.Error(err)
		println("get nonce error")
		return
	}
	transaction := util.Transaction{
		ChainId: Chain, FromChainId: Chain, ToChainId: Chain, From: from,
		To: "", Value: "0", Input: "", Nonce: strconv.Itoa(int(nonce)),
	}

	privateKey, err := crypto.HexToECDSA(key)
	hash, err := contract.Deploy(transaction, bytecode, privateKey, Symbol, Name, decimal, amount)
	if err != nil {
		log.Error(err)
		fmt.Println("get hash error")
		return
	}
	log.Info("contract hash:", hash)
	time.Sleep(time.Second * 10)
	receipt, err := connection.Thk.GetTransactionByHash(Chain, hash)
	if err != nil {
		log.Error(err)
		fmt.Println("get hash error")
		return
	}
	log.Info("contract addr:", receipt.ContractAddress)
	to := receipt.ContractAddress
	newtransaction := util.Transaction{
		ChainId: Chain, FromChainId: Chain, ToChainId: Chain, From: from,
		To: to, Value: "0", Input: "", Nonce: strconv.Itoa(int(nonce)),
	}
	result, err := contract.Call(newtransaction, "symbol")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log("result:", result)
	var str string
	err = contract.Parse(result, "symbol", &str)
	if err != nil {
		log.Error(err)
		println("failed")
		return
	}

	var token Token
	token.Name = Name
	token.Symbol = Symbol
	token.Total = 100000
	token.ContractAddress = receipt.ContractAddress
	token.ABI = unmarshalResponse.Abi
	token.Icon = "icon"
	token.Website = "www.thinkey.com"
	token.Introduction = Des
	token.Date = time.Now().Format("2006-01-02")
	token.ChainId = "2"
	token.Decimal = 8
	PostInfo(token)
	//PostFile(token, "../resources/abc.sol")
}

func PostInfo(token Token) {
	stringUrl := "http://ext.thinkey.xyz/v1/wallet/token/tokeninfo/"
	bodybuf := bytes.NewBufferString("")
	bodywriter := multipart.NewWriter(bodybuf)
	bodywriter.SetBoundary("Pp7Ye2EeWaFDdAY")
	err := bodywriter.WriteField("Name", token.Name)
	err = bodywriter.WriteField("Symbol", token.Symbol)
	err = bodywriter.WriteField("Total", fmt.Sprintf("%f", token.Total))
	err = bodywriter.WriteField("ContractAddress", token.ContractAddress)
	err = bodywriter.WriteField("ABI", token.ABI)
	err = bodywriter.WriteField("Icon", token.Icon)
	err = bodywriter.WriteField("Website", token.Website)
	err = bodywriter.WriteField("Introduction", token.Introduction)
	err = bodywriter.WriteField("Date", token.Date)
	err = bodywriter.WriteField("ChainId", token.ChainId)
	err = bodywriter.WriteField("Decimal", fmt.Sprintf("%d", token.Decimal))
	bodywriter.Close()
	reqreader := io.MultiReader(bodybuf)
	resp, err := http.Post(stringUrl,
		bodywriter.FormDataContentType(),
		reqreader)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("读取回应消息异常")
	}
	str := string(body)
	fmt.Println("发送回应数据:", str)
}

func TestRun(t *testing.T) {
	var token Token
	token.Name = Name
	token.Symbol = Symbol
	token.Total = 100000
	token.ContractAddress = ""
	token.ABI = ""
	token.Icon = "icon"
	token.Website = ""
	PostFile(token, "../resources/abc.sol")
}
func PostFile(token Token, name string) {
	filebytes, err := ioutil.ReadFile(name)
	stringUrl := "http://192.168.1.164:8201/v1/thkadmin/wallet/token/"
	bodybuf := bytes.NewBufferString("")
	bodywriter := multipart.NewWriter(bodybuf)
	bodywriter.SetBoundary("Pp7Ye2EeWaFDdAY")
	err = bodywriter.WriteField("Name", token.Name)
	err = bodywriter.WriteField("Symbol", token.Symbol)
	err = bodywriter.WriteField("Total", fmt.Sprintf("%f", token.Total))
	err = bodywriter.WriteField("ContractAddress", token.ContractAddress)
	err = bodywriter.WriteField("ABI", token.ABI)
	err = bodywriter.WriteField("Icon", token.Icon)
	err = bodywriter.WriteField("Website", token.Website)
	filename := path.Base(name)
	_, err = bodywriter.CreateFormFile(token.Icon, filename)
	if err != nil {
		fmt.Printf("创建FormFile1文件信息异常！")
	}
	bodybuf.Write(filebytes)
	bodywriter.Close()
	//application/json
	//multipart/form-data

	reqreader := io.MultiReader(bodybuf)

	//req, err := http.NewRequest("POST", stringUrl, reqreader)
	//if err != nil {
	//	fmt.Printf("站点相机上传图片，创建上次请求异常！异常信息")
	//
	//}
	//req.Header.Set("Connection", "close")
	//req.Header.Set("Pragma", "no-cache")
	//req.Header.Set("Content-Type", bodywriter.FormDataContentType())
	//req.ContentLength = int64(bodybuf.Len())
	//fmt.Printf("发送消息长度:")
	//client := &http.Client{}
	//resp, err := client.Do(req)
	//defer resp.Body.Close()

	resp, err := http.Post(stringUrl,
		bodywriter.FormDataContentType(),
		reqreader)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("读取回应消息异常")
	}
	fmt.Println("发送回应数据:", string(body))

}

func GetContractJson(ctct map[string]interface{}) (ContractJson, ContractJson, error) {
	var contractJson, ERC20Json, TokenVestingJson ContractJson
	for keyname, value := range ctct {
		contractJson.ContractName = keyname
		arr := strings.Split(contractJson.ContractName, ":")
		length := len(arr) - 1
		if (arr[length] == "ERC20") {
			mapvalue := value.(map[string]interface{})
			ERC20Json.ByteCode = mapvalue["code"].(string)
			info := mapvalue["info"].(map[string]interface{})
			abidef := info["abiDefinition"]
			abibytes, _ := json.Marshal(abidef)
			ERC20Json.ABI = string(abibytes)
		}
		if (arr[length] == "TokenVesting") {
			mapvalue := value.(map[string]interface{})
			TokenVestingJson.ByteCode = mapvalue["code"].(string)
			info := mapvalue["info"].(map[string]interface{})
			abidef := info["abiDefinition"]
			abibytes, _ := json.Marshal(abidef)
			TokenVestingJson.ABI = string(abibytes)
		}
	}
	return ERC20Json, TokenVestingJson, nil
}
