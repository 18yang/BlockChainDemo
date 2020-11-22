package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/gob"
	"log"
	"time"
)

//1. 定义结构
type Block struct {
	//1. 版本号
	Version uint64
	//2. 前区块哈希
	PrevHash []byte
	//3. Merkel根
	//(所有的交易每一个进行单独哈希，再两两结合算出哈希，直到最后得到的一个哈希为Merkel root)
	MerkelRoot []byte
	//4. 时间戳
	TimeStamp uint64
	//5. 难度值
	Diffculty uint64
	//6. 随机数，挖矿要找的数
	Nonce uint64
	//7. 当前区块哈希 正常比特币区块中没有当前区块的哈希，为了方便简化//TODO
	Hash []byte
	//8. 数据
	Transactions []*Transaction
}
//实现一个辅助函数，将uint64转成[]byte
func Uint64ToByte(num uint64) []byte {
	var buffer bytes.Buffer

	//通过重新二进制编码转换
	err := binary.Write(&buffer, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}
	return buffer.Bytes()
}
//2. 创建区块
func NewBlock(txs []*Transaction, preBlockHash []byte) *Block {
	block := Block{
		Version:    00,
		PrevHash:   preBlockHash,
		MerkelRoot: []byte{},
		TimeStamp:  uint64(time.Now().Unix()),
		Diffculty:  0,
		Nonce:      0,
		Hash:       []byte{},
		Transactions: txs,
	}
	block.MerkelRoot = block.MakeMerkelRoot()
	//设置哈希值
	//block.SetHash()
	//创建一个pow对象 不停的进行hash运算 查找随机数
	pow := NewProofOfWork(&block)
	hash, nonce := pow.Run()
	//根据挖矿结果对区块数据进行更新
	block.Hash = hash
	block.Nonce = nonce

	return &block
}

//3. 生成哈希
func (block *Block) SetHash() {
	// 1. 拼装数据
	/*
		blockInfo = append(blockInfo, Uint64ToByte(block.Version)...)
		blockInfo = append(blockInfo, block.PrevHash...)
		blockInfo = append(blockInfo, block.MerkelRoot...)
		blockInfo = append(blockInfo, block.Data...)
		blockInfo = append(blockInfo, Uint64ToByte(block.TimeStamp)...)
		blockInfo = append(blockInfo, Uint64ToByte(block.Diffculty)...)
		blockInfo = append(blockInfo, Uint64ToByte(block.Nonce)...)
	*/
	tmp:= [][]byte{
		//只对区块头拼接做hash运算
		Uint64ToByte(block.Version),
		block.PrevHash,
		block.MerkelRoot,
		Uint64ToByte(block.TimeStamp),
		Uint64ToByte(block.Diffculty),
		Uint64ToByte(block.Nonce),
	}
	//将二维切片数组链接起来，返回一个一维切片数组
	blockInfo := bytes.Join(tmp, []byte{})
	// 2. sha256
	hash := sha256.Sum256(blockInfo)
	block.Hash = hash[:]
}
//序列化
func (block *Block)Serialize() []byte {
	var buffer bytes.Buffer
	//使用gob进行序列化得到字节流
	// 1. 定义一个编码器
	encoder := gob.NewEncoder(&buffer)
	//2. 使用编码器编码
	err := encoder.Encode(&block)
	if err != nil {
		panic(err)
	}

	//fmt.Printf("%v\n",buffer)
	return buffer.Bytes()
}
//反序列化
func Deserialize(data []byte) Block {
	decoder := gob.NewDecoder(bytes.NewReader(data))
	//2. 使用解码器解码
	var block Block
	err := decoder.Decode(&block)
	if err != nil {
		panic(err)
	}
	return block
}
func (block *Block) MakeMerkelRoot () []byte{
	//TODO
	return []byte{}
}