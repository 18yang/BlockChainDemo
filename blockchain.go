package main

import (
	"BlockChain/bolt"
	"log"
)

//4. 引入区块链
type BlockChain struct {
	//定义一个区块链数组
	//blocks []*Block
	db *bolt.DB
	tail []byte //存储最后一个区块的值
}

const blockChainDb = "blockChain.db"
const blockBucket  = "blockBucket"

//返回一个链
func NewBlockChain() *BlockChain {
	//作为第一个创世块并加入到区块链中
	genesisBlock := GenesisBlock()
	//return &BlockChain{
	//	blocks: []*Block{genesisBlock},
	//}
	//该变量保存最后一个区块的哈希
	var lastHash []byte
	//1. 打开数据库
	db, err := bolt.Open(blockChainDb, 0600, nil)
	if err != nil {
		log.Panic("打开数据库失败")
	}
	//2. 找到抽屉bucket（没有就创建）
	db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			//没有抽屉，创建
			bucket, err = tx.CreateBucket([]byte(blockBucket))
			if err != nil {
				log.Panic("创建bucket("+blockBucket+")失败")
			}
			//3. 写数据  hash 最为key block大的字节流最为value
			bucket.Put(genesisBlock.Hash, genesisBlock.Serialize())
			bucket.Put([]byte("lastHash"), genesisBlock.Hash)
			lastHash = genesisBlock.Hash
		}else{
			lastHash = bucket.Get([]byte("lastHash"))
		}
		return nil
	})
	return &BlockChain{
		db:   db,
		tail: lastHash,
	}
}

//定义一个创世块
func GenesisBlock() *Block {
	return NewBlock("创世块", []byte{})
}

//5. 添加区块
func (chain *BlockChain) AddBlock(data string) {
	//获取前区块的哈希
	db := chain.db  // 获取区块链数据库
	lastHash := chain.tail
	//a. 创建新的区块
	db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			log.Panic("bucket不应该为空，请检查！")
		}
		//b. 添加到区块链数据库中
		block := NewBlock(data, lastHash)
		//写数据  hash 最为key block大的字节流最为value
		bucket.Put(block.Hash, block.Serialize())
		bucket.Put([]byte("lastHash"), block.Hash)
		lastHash = block.Hash
		//更新一下内存中的区块链
		chain.tail = block.Hash

		return nil
	})


}

func (bc *BlockChain)NewIterator() *BlockChainIterator {
	return &BlockChainIterator{
		db:                 bc.db,
		//最初指向区块链的最后一个区块，随着Next方法，不断变化
		currentHashPointer: bc.tail,
	}
}
