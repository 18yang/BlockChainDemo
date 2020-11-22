package main

func main() {
	chain := NewBlockChain("班长")
	cli := CLI{chain}
	cli.Run()
	//chain.AddBlock("11111111111")
	//chain.AddBlock("22222222222")
	//
	////调用迭代器，返回每一个区块
	//iterator := chain.NewIterator()
	//
	//for  {
	//	block := iterator.Next()
	//	fmt.Println("=============================================")
	//	fmt.Printf("前区块的哈希值:  %x\n", block.PrevHash)
	//	fmt.Printf("当前区块的哈希值: %x\n", block.Hash)
	//	fmt.Printf("当前区块的数据:  %s\n", block.Data)
	//	fmt.Println("=============================================")
	//	fmt.Println()
	//	if len(block.PrevHash) == 0  {
	//		fmt.Println("遍历结束")
	//		return
	//	}
	//}
	chain.db.Close()
}
