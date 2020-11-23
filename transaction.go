package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
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
	Value float64
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
//实现一个函数，判断当前的交易是否为挖矿交易
func (tx *Transaction) IsCoinbase() bool {
	//1. 交易input只有一个
	//2. 交易id为空
	//3. 交易的index 为 -1
	if len(tx.TXInputs) == 1 && len(tx.TXInputs[0].TXid) == 0 && tx.TXInputs[0].Index == -1 {
		return true
	}
	return false
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
		Value:      reward,
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
//创建普通交易
func NewTransaction(from, to string,amount float64,bc *BlockChain) *Transaction{
	//TODO
	//1. 找到最合理的utxo集合
	utxos, resValue := bc.FindNeedUTXOs(from,amount)

	if resValue < amount {
		fmt.Printf("余额不足，交易失败")
	}

	var inputs []TXInput
	var outputs []TXOutput
	//2. 将utxo逐一装成inputs   创建交易输入
	for id, indexArray := range utxos {
		for _,i := range indexArray {
			input := TXInput{
				TXid:  []byte(id),
				Index: int64(i),
				Sig:   from,
			}
			inputs = append(inputs,input)
		}
	}
	//3. 创建outputs 创建交易输出
	output:= TXOutput{
		Value:      amount,
		PukKeyHash: to,
	}
	outputs = append(outputs,output)
	//4. 如果有零钱，找零
	if resValue > amount {
		outputs = append(outputs,TXOutput{resValue-amount,from})
	}
	//将交易输入与输出添加到一个交易中
	tx := Transaction{[]byte{},inputs,outputs}
	tx.SetHash()
	return &tx
}
