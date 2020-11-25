package main

import (
	"BlockChainProject/base58"
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"io/ioutil"
	"log"
	"os"
)

const walletFile = "wallet.dat"

//定一个 Wallets结构，它保存所有的wallet以及他的地址
type Wallets struct {
	//map[地址]钱包
	WalletsMap map[string] *Wallet
}

//创建方法
func NewWallets() *Wallets {
	var ws Wallets
	//ws.WalletsMap = make(map[string]*Wallet)
	ws.loadFile()
	return &ws
}
func (ws *Wallets)CreateWallet() string{
	wallet := NewWallet()
	address := wallet.NewAddress()
	//var wallets Wallets
	//wallets.WalletsMap = make(map[string] *Wallet)
	ws.WalletsMap[address] = wallet
	ws.saveToFile()
	return address
}

//保存文件
func (ws *Wallets)saveToFile()  {
	var buffer bytes.Buffer
	//应为椭圆曲线是接口类型，不能直接进行序列化，需要先注册
	gob.Register(elliptic.P256())
	encode := gob.NewEncoder(&buffer)
	err := encode.Encode(ws)
	if err != nil {
		log.Panic(err)
	}
	//对钱包序列化，存储到本地
	ioutil.WriteFile(walletFile,buffer.Bytes(),0600)
}

func (ws *Wallets)loadFile()  {
	//在读取之间，需要确定文件是否存在，不存在即退出
	_, err2 := os.Stat(walletFile)
	if os.IsNotExist(err2) {
		ws.WalletsMap = make(map[string] *Wallet)
		return
	}
	//读取内容
	content, err := ioutil.ReadFile(walletFile)
	if err != nil {
		log.Panic(err)
	}
	//解码
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(content))
	var wsLocal Wallets
	err = decoder.Decode(&wsLocal)
	if err != nil {
		log.Panic(err)
	}
	//注意不能直接用ws和wsLocal直接赋值
	ws.WalletsMap = wsLocal.WalletsMap
}
func (ws *Wallets)GetAllAddresses() []string {
	var addresses []string
	//遍历钱包，将所有的key取出来返回
	for address := range ws.WalletsMap{
		addresses = append(addresses, address)
	}
	return addresses
}
//通过地址返回公钥的哈希值
func GetPubKeyFromAddress(address string) []byte {
	//1. 解码
	//2. 截取出公钥哈希：去除version（1字节），去除校验码（4字节）
	addressByte := base58.Decode(address) //25字节
	lenth := len(addressByte)

	pubKeyHash := addressByte[1:lenth-4]

	return pubKeyHash
}