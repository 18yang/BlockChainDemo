package main

import (
	"BlockChainProject/bolt"
	"bytes"
	"fmt"
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
func NewBlockChain(address string) *BlockChain {

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
			//作为第一个创世块并加入到区块链中
			genesisBlock := GenesisBlock(address)
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
func GenesisBlock(address string) *Block {
	coinbase := NewCoinBaseTX(address,"创世块")
	return NewBlock([]*Transaction{coinbase}, []byte{})
}

//5. 添加区块
func (chain *BlockChain) AddBlock(txs []*Transaction) {
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
		block := NewBlock(txs, lastHash)
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
//找到指定地址的所有UTXO
func (bc *BlockChain) FindUTXOs(address string) []TXOutput {
	var UTXO []TXOutput
	//map[string][]uint64
	//定义一个map来保存笑给过的output，key是这个output的交易id，value是这个交易中索引的数组
	spentOutputs := make(map[string][]int64)
	//创建迭代器
	it := bc.NewIterator()
	for   {
		//遍历区块
		block := it.Next()
		//遍历交易
		for _,tx := range block.Transactions {
			fmt.Printf("current txid: %x\n",tx.TXID)

		OUTPUT:
			//遍历output， 找到和自己相关的utxo（再添加output之前检查一下自己是否消耗过）
			for i , output := range tx.TXoutputs {
				fmt.Printf("current index: %x\n",i)
				//在这里做一个过滤，将所有消耗过的output和当前的所即将添加output对比一下
				//如果相同，即跳过，否则添加
				//如果当前的交易id存在于我们已经表示的map，那么说明这个交易是有消耗过的
				if spentOutputs[string(tx.TXID)] != nil{
					for _,j := range spentOutputs[string(tx.TXID)]{
						if int64(i) == j {
							//当前准备添加output已经消耗过了
							continue OUTPUT
						}
					}
				}

				//这个output和我们目标的地址相同，满足条件，加到返回utxo数组中
				if output.PukKeyHash == address {
					UTXO = append(UTXO, output)
				}
			}
			//如果当前交易时挖矿交易的话，那么不做遍历
			if !tx.IsCoinbase() {


				//遍历input ， 找到自己花费过的utxo集合（把自己消耗过的标识出来）
				for _,input := range tx.TXInputs {
					//判断一下当前这个input和目标是否一致，如果相同，表示是消耗过的output
					if input.Sig == address {
						//indexArray := spentOutputs[string(input.TXid)]
						spentOutputs[string(input.TXid)] = append(spentOutputs[string(input.TXid)], input.Index)
					}
				}
			}else{
				fmt.Printf("这是coinbase, 不做遍历\n")
			}
		}
		if len(block.PrevHash) == 0 {
			fmt.Println("区块链遍历完成退出")
			break
		}
	}
	return UTXO
}

func (bc *BlockChain) FindNeedUTXOs(from string, amount float64) (map[string][]uint64,float64) {
	//找到的合理的utxo集合
	utxos := make(map[string][]uint64)
	//找到utxos里面包含钱的总数
	calculate := 0.0
	//map[string][]uint64
	//定义一个map来保存笑给过的output，key是这个output的交易id，value是这个交易中索引的数组
	spentOutputs := make(map[string][]int64)
	//创建迭代器
	it := bc.NewIterator()
	for   {
		//遍历区块
		block := it.Next()
		//遍历交易
		for _,tx := range block.Transactions {
			//fmt.Printf("current txid: %x\n",tx.TXID)

		OUTPUT:
			//遍历output， 找到和自己相关的utxo（再添加output之前检查一下自己是否消耗过）
			for i , output := range tx.TXoutputs {
				//fmt.Printf("current index: %x\n",i)
				//在这里做一个过滤，将所有消耗过的output和当前的所即将添加output对比一下
				//如果相同，即跳过，否则添加
				//如果当前的交易id存在于我们已经表示的map，那么说明这个交易是有消耗过的
				if spentOutputs[string(tx.TXID)] != nil{
					for _,j := range spentOutputs[string(tx.TXID)]{
						if int64(i) == j {
							//当前准备添加output已经消耗过了
							continue OUTPUT
						}
					}
				}

				//这个output和我们目标的地址相同，满足条件，加到返回utxo数组中
				if output.PukKeyHash == from {
					if calculate < amount {
						//1. 把utxo加进来
						utxos[string(tx.TXID)] = append(utxos[string(tx.TXID)], uint64(i))
						//2. 统计一下当前utxo的总额
						calculate += output.Value
						//3. 比较一下是否满足转账需求
						//	a. 满足的话，直接返回，
						if calculate >= amount {
							return utxos,calculate
						}
						//	b. 不满足继续统计

					}
				}
			}
			//如果当前交易时挖矿交易的话，那么不做遍历
			if !tx.IsCoinbase() {
				//遍历input ， 找到自己花费过的utxo集合（把自己消耗过的标识出来）
				for _,input := range tx.TXInputs {
					//判断一下当前这个input和目标是否一致，如果相同，表示是消耗过的output
					if input.Sig == from {
						//indexArray := spentOutputs[string(input.TXid)]
						spentOutputs[string(input.TXid)] = append(spentOutputs[string(input.TXid)], input.Index)
					}
				}
			}//else{
			//	fmt.Printf("这是coinbase, 不做遍历\n")
			//}
		}
		if len(block.PrevHash) == 0 {
			//fmt.Println("区块链遍历完成退出")
			break
		}
	}

	return utxos,calculate
}


func (bc *BlockChain) Printchain() {

	blockHeight := 0
	bc.db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte("blockBucket"))

		//从第一个key-> value 进行遍历，到最后一个固定的key时直接返回
		b.ForEach(func(k, v []byte) error {
			if bytes.Equal(k, []byte("LastHashKey")) {
				return nil
			}

			block := Deserialize(v)
			//fmt.Printf("key=%x, value=%s\n", k, v)
			fmt.Printf("=============== 区块高度: %d ==============\n", blockHeight)
			blockHeight++
			fmt.Printf("版本号: %d\n", block.Version)
			fmt.Printf("前区块哈希值: %x\n", block.PrevHash)
			fmt.Printf("梅克尔根: %x\n", block.MerkelRoot)
			fmt.Printf("时间戳: %d\n", block.TimeStamp)
			fmt.Printf("难度值(随便写的）: %d\n", block.Difficulty)
			fmt.Printf("随机数 : %d\n", block.Nonce)
			fmt.Printf("当前区块哈希值: %x\n", block.Hash)
			fmt.Printf("区块数据 :%s\n", block.Transactions[0].TXInputs[0].Sig)
			return nil
		})
		return nil
	})
}