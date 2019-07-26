package main

import (
	"context"
	"encoding/hex"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
	"time"
	"web3.go/web3/thk/abi"
)

var jsonName = "./test/resources/dahantest.json"
var jsonName1 = "./test/resources/dahan.json"

type TruffleContract struct {
	Abi      string `json:"abi"`
	Bytecode string `json:"bytecode"`
}
type Token struct {
	Name            string  `json:"name"`
	Symbol          string  `json:"symbol"`
	Total           float64 `json:"total"`
	ContractAddress string  `json:"contractAddress"`
	ABI             string  `json:"abi"`
	Icon            string  `json:"icon"`
	Website         string  `json:"website"`
	Introduction    string  `json:"introduction"`
	State           string  `json:"state"`
	Date            string  `json:"date"`
	ChainId         string  `json:"chainid"`
}

func main() {
	contractAddress := "0x74b00d00f2c0d0cd932a894054454c4798c11c3e"
	input := "0x23b872dd0000000000000000000000002c7536e3605d9c16a7a3d7b1898e529396a65c230000000000000000000000002c7536e3605d9c16a7a3d7b1898e529396a65c240000000000000000000000000000000000000000000000000000000000000014"
	mgo, err := InitMongod()
	if err != nil {
		println(err)
	}
	str := GetInputStr(mgo, contractAddress, input)
	println(str)
}
func GetInputStr(db *mongo.Database, contractAddress, input string) interface{} {
	findfilter := bson.D{{"contractaddress", contractAddress}}
	var token Token
	res := db.Collection("wallet").FindOne(context.Background(), findfilter)
	if res.Err() != nil {
		return nil
	}
	err:=res.Decode(&token)
	if err != nil {
		return nil
	}
	abistr := token.ABI
	//name := "transferFrom"

	// 解析Abi格式成为Json格式
	abiDecoder, err := abi.JSON(strings.NewReader(abistr))
	if err != nil {
		println(err)
		return nil
	}

	// 剔除最前面的0x标记
	var inputString string = ""
	hexFlag := strings.Index(input, "0x")
	if hexFlag == -1 {
		inputString = input
	} else {
		inputString = input[2:]
	}
	decodeBytes, err := hex.DecodeString(inputString)
	if err != nil {
		println(err)
		return nil
	}
	inputMap := map[string]interface{}{}
	method, err := abiDecoder.MethodById(decodeBytes)
	if err != nil {
		println(err)
		return nil
	}
	if err1 := abiDecoder.InputUnpack(inputMap, method.Name, decodeBytes[4:]); err1 != nil {
		println("DecodeAbi ", "err", err)
		println(err)
	}
	return inputMap
}

func InitMongod() (collection *mongo.Database, err error) {
	opts := options.ClientOptions{Hosts: []string{"192.168.1.108:27017"}}
	//credential := options.Credential{
	//	Username: config.Key("username").String(),
	//	Password: config.Key("password").String(),
	//	AuthSource: config.Key("db").String(),
	//}
	//opts.Auth = &credential
	client, err := mongo.NewClient(&opts)

	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}
	collection = client.Database("publicChain")
	return collection, nil
}
