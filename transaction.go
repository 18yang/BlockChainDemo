package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
)
const reward = 25
//定义交易结构
type Transaction struct {
	TXID []byte  //交易id
	TXInputs []TXInput //交易输入数组
	TXoutputs []TXOutput //交易输出数组
}
//定义交易输入
type TXInput struct {
	//引用的交易ID
	TXid []byte
	//引用的output的索引值
	Index int64
	//解锁脚本，我们用地址来模拟
	Sig string
}
//定义交易输出
type TXOutput struct {
	//转账金额
	value float64
	//锁定脚本，用地址模拟
	PukKeyHash string
}
//设置交易ID
func (tx *Transaction) SetHash(){
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)

	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	data := buffer.Bytes()
	hash := sha256.Sum256(data)
	tx.TXID = hash[:]
}
//提供创建交易方法(挖矿交易)
func NewCoinBaseTX(address string,data string) *Transaction{
	//挖矿交易的特点：
	//1. 只有一个input
	//2. 无需引用交易id
	//3. 无需引用index
	//矿工由于挖矿时无需指定签名，所以无需这个Sig字段，可以由矿工自由填写
	input := TXInput{
		TXid:  []byte{},
		Index: -1,
		Sig:   data,
	}
	output := TXOutput{
		value:      reward,
		PukKeyHash: address,
	}
	//对于挖矿交易来说，只有一个input和output
	tx := Transaction{
		TXID:      []byte{},
		TXInputs:  []TXInput{input},
		TXoutputs: []TXOutput{output},
	}
	tx.SetHash()
	return &tx
}
//创建挖矿交易
//根据交易调整程序

