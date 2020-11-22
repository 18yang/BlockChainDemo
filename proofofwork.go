package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
)

//定义一个工作量证明的结构ProofOfWork
type ProofOfWork struct {
	//a. block
	block *Block
	//b. 目标值 big.Int 是一个非常大的数，有其特殊的方法
	target *big.Int
}
//提供创建POW的函数
func NewProofOfWork(block *Block) *ProofOfWork {
	pow := ProofOfWork{
		block:  block,
	}
	//指定的目标值   --->   工作量证明就是为了找到比目标值小的hash数
	targetStr := "0000f0000000000000000000000000000000000000000000000000000000000"
	tmpInt := big.Int{}//转成 big.int 类型
	//将目标值赋值给tmpInt 并指定为16进制
	tmpInt.SetString(targetStr,16)
	pow.target = &tmpInt
	return &pow
}
//提供计算不断计算hash的函数
func (pow *ProofOfWork) Run() ([]byte, uint64) {
	fmt.Println("开始挖矿...")
	var nonce uint64 //随机数
	block := pow.block  //区块
	var hash [32]byte   //哈希值
	for   {
		//1. 拼装数据（区块数据，还有不断变化的随机数 Nonce）
		tmp:= [][]byte{
			Uint64ToByte(block.Version),
			block.PrevHash,
			block.MerkelRoot,
			block.Data,
			Uint64ToByte(block.TimeStamp),
			Uint64ToByte(block.Diffculty),
			Uint64ToByte(nonce),
		}
		//将二维切片数组链接起来，返回一个一维切片数组
		blockInfo := bytes.Join(tmp, []byte{})
		//2. 做哈希运算
		hash = sha256.Sum256(blockInfo)
		//3. 与目标值进行比较
		tmpInt := big.Int{}
		//将我们得到的hash数组转换成一个 big.int 应为目标值是big.int
		tmpInt.SetBytes(hash[:])
		//和目标值进行比较，小就退出返回，大就让 nonce++  继续循环
		if tmpInt.Cmp(pow.target) == -1{
			// -1 为小于  0 为等于    1 为大于
			fmt.Printf("挖矿成功！ hash：%x,%d\n",hash,nonce)
			break
		}else {
			nonce++
		}
	}
	return hash[:], nonce
}
//提供一个校验函数