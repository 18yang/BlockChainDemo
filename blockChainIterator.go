package main

import (
	"BlockChain/bolt"
	"log"
)

type BlockChainIterator struct {
	db *bolt.DB
	currentHashPointer []byte
}

//用于将指针不断向前迭代，取出数据库中的值
func (it *BlockChainIterator)Next() *Block {
	var block Block
	it.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			log.Panic("迭代器遍历时，bucket不应该为空，请检查！")
		}
		//未解码的数据
		blockTmp := bucket.Get(it.currentHashPointer)
		block = Deserialize(blockTmp)
		//指针左移
		it.currentHashPointer = block.PrevHash
		return nil
	})
	return &block
}
