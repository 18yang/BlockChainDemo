package main

import "fmt"

func (cli *CLI) PrintBlockChain() {
	//cli.bc.Printchain()
	//fmt.Printf("打印区块链完成\n")
	chain := cli.bc
	iterator := chain.NewIterator()

	for {
		block := iterator.Next()
		fmt.Println("=============================================")
		fmt.Printf("版本号:  %d\n", block.Version)
		fmt.Printf("前区块的哈希值:  %x\n", block.PrevHash)
		fmt.Printf("Merkel根:  %x\n", block.MerkelRoot)
		fmt.Printf("时间戳: %d\n", block.TimeStamp)
		fmt.Printf("难度值:  %d\n", block.Difficulty)
		fmt.Printf("随机数:  %d\n", block.Nonce)
		fmt.Printf("当前区块哈希:  %x\n", block.Hash)
		fmt.Printf("区块数据:  %s\n", block.Transactions[0].TXInputs[0].PubKey)
		fmt.Println("=============================================")
		fmt.Println()
		if len(block.PrevHash) == 0 {
			fmt.Printf("遍历结束\n")
			return
		}
	}

}

func (cli *CLI) PrintBlockChainReverse() {
	bc := cli.bc
	//创建迭代器
	it := bc.NewIterator()

	//调用迭代器，返回我们的每一个区块数据
	for {
		//返回区块，左移
		block := it.Next()

		for _, tx := range block.Transactions {
			fmt.Println(tx)
		}
		if len(block.PrevHash) == 0 {
			fmt.Printf("区块链遍历结束！")
			break
		}
	}
}


func (cli *CLI) GetBalance(address string) {
	//校验地址是否有效
	if !IsValidAddress(address){
		fmt.Printf("地址无效： %s\n",address)
		return
	}

	pubKeyHash := GetPubKeyFromAddress(address)
	utxos := cli.bc.FindUTXOs(pubKeyHash)

	total := 0.0
	for _, utxo := range utxos {
		total += utxo.Value
	}
	fmt.Printf("%s的余额为： %f\n", address, total)
}

func (cli *CLI) Send(from, to string, amount float64, miner, data string) {
	//校验地址是否有效
	if !IsValidAddress(from){
		fmt.Printf("地址无效： %s\n",from)
		return
	}
	if !IsValidAddress(to){
		fmt.Printf("地址无效： %s\n",to)
		return
	}
	if !IsValidAddress(miner){
		fmt.Printf("地址无效： %s\n",miner)
		return
	}

	//1. 创建挖矿交易
	coinbase := NewCoinBaseTX(miner, data)
	//2. 创建一个普通交易
	tx := NewTransaction(from, to, amount, cli.bc)
	if tx == nil {
		//fmt.Printf("无效的交易")
		return
	}
	//3. 添加到区块
	cli.bc.AddBlock([]*Transaction{coinbase, tx})
	fmt.Printf("转账成功！\n")
}

func (cli *CLI)NewWallet()  {
	//wallet := NewWallet()
	//address := wallet.NewAddress()
	//fmt.Printf("私钥：%v\n",wallet.Private)
	//fmt.Printf("公钥：%v\n",wallet.PubKey)
	//fmt.Printf("地址：%s\n",address)
	ws := NewWallets()
	address := ws.CreateWallet()
	fmt.Printf("地址： %s\n",address)
}
func (cli *CLI) ListAddresses()  {
	ws := NewWallets()
	addresses := ws.GetAllAddresses()
	for _,address :=range addresses {
		fmt.Printf("地址： %s\n",address)
	}
}


