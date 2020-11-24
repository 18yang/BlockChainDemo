package main

import "fmt"

func (cli *CLI) PrintBlockChain(){
	cli.bc.Printchain()
	fmt.Println("打印区块链完成\n")
}

func (cli *CLI)PrintBlockChainReverse()  {
	chain := cli.bc
	iterator := chain.NewIterator()

	for  {
		block := iterator.Next()
		fmt.Println("=============================================")
		fmt.Printf("版本号:  %d\n", block.Version)
		fmt.Printf("前区块的哈希值:  %x\n", block.PrevHash)
		fmt.Printf("Merkel根:  %x\n", block.MerkelRoot)
		fmt.Printf("时间戳: %d\n", block.TimeStamp)
		fmt.Printf("难度值:  %d\n", block.Difficulty)
		fmt.Printf("随机数:  %d\n", block.Nonce)
		fmt.Printf("当前区块哈希:  %x\n", block.Hash)
		fmt.Printf("区块数据:  %s\n", block.Transactions[0].TXInputs[0].Sig)
		fmt.Println("=============================================")
		fmt.Println()
		if len(block.PrevHash) == 0  {
			fmt.Println("遍历结束")
			return
		}
	}
}
func (cli *CLI)AddBlock(data string)  {
	//cli.bc.AddBlock(data)//TODO
	fmt.Printf("添加区块成功！\n")
}

func (cli *CLI) GetBalance(address string)  {
	utxos := cli.bc.FindUTXOs(address)

	total := 0.0
	for _, utxo := range utxos  {
		total += utxo.Value
	}
	fmt.Printf("%s的余额为： %f\n",address,total)
}

func (cli *CLI)Send(from ,to string, amount float64, miner ,data string)  {
	//1. 创建挖矿交易
	coinbase := NewCoinBaseTX(miner,data)
	//2. 创建一个普通交易
	tx := NewTransaction(from, to, amount, cli.bc)
	if tx == nil {
		//fmt.Printf("无效的交易")
		return
	}
	//3. 添加到区块
	cli.bc.AddBlock([]*Transaction{coinbase,tx})
	fmt.Printf("转账成功！\n")
	////1. 创建挖矿交易
	//coinbase := NewCoinBaseTX(miner, data)
	////2. 创建一个普通交易
	//tx := NewTransaction(from, to, amount, cli.bc)
	//fmt.Printf("22222222222\n")
	//if tx == nil {
	//	fmt.Printf("3333333333\n")
	//	//fmt.Printf("无效的交易")
	//	return
	//}
	////3. 添加到区块
	//fmt.Printf("444444444444\n")
	//cli.bc.AddBlock([]*Transaction{coinbase, tx})
	//fmt.Printf("555555555555\n")
	//fmt.Printf("转账成功！\n")
}