package main

import (
	"BlockChainProject/base58"
	"BlockChainProject/ripemd160"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"log"
)

//每一个钱包都保存了公钥，私钥对
type Wallet struct {
	//私钥
	Private *ecdsa.PrivateKey
	//PubKey *ecdsa.PublicKey
	//约定，这里的PubKey不存储原始的公钥，而是存储X和Y拼接的字符串，再校验端重新拆分
	PubKey []byte
}
//创建钱包
func NewWallet() *Wallet {
	//创建曲线
	curve := elliptic.P256()
	//生成私钥
	privateKey, err := ecdsa.GenerateKey(curve,rand.Reader)
	if err != nil {
		log.Panic(err)
	}
	//生成公钥
	pubKeyOrig := privateKey.PublicKey

	//拼接X，Y
	pubKey := append(pubKeyOrig.X.Bytes(),pubKeyOrig.Y.Bytes()...)

	return &Wallet{
		Private: privateKey,
		PubKey:  pubKey,
	}
}

//生成地址
func (w *Wallet)NewAddress() string {
	pubKey := w.PubKey
	rip160HashValue := HashPubKey(pubKey)
	//拼接version
	version := byte(00)
	payload := append([]byte{version},rip160HashValue...)

	//25字节数据
	payload = CheckSum(payload)
	//对地址进行编码
	address := base58.Encode(payload)

	return address
}

func HashPubKey(data []byte) []byte {
	hash := sha256.Sum256(data)
	//理解为编码器
	rip160hasher := ripemd160.New()
	_, err := rip160hasher.Write(hash[:])
	if err != nil {
		log.Panic(err)
	}
	//返回rip160的哈希结果
	rip160HashValue := rip160hasher.Sum(nil)
	return rip160HashValue
}

func CheckSum(data []byte) []byte {
	hash1 := sha256.Sum256(data)
	hash2 := sha256.Sum256(hash1[:])
	//前四字节校验码
	checkCode := hash2[:4]
	//25字节数据
	data = append(data,checkCode...)
	return data
}