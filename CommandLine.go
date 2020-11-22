package main

import "fmt"

func (cli *CLI)PrintBlockChain()  {
	chain := cli.bc
	iterator := chain.NewIterator()

	for  {
		block := iterator.Next()
		fmt.Println("=============================================")
		fmt.Printf("版本号:  %d\n", block.Version)
		fmt.Printf("前区块的哈希值:  %x\n", block.PrevHash)
		fmt.Printf("Merkel根:  %x\n", block.MerkelRoot)
		fmt.Printf("时间戳: %d\n", block.TimeStamp)
		fmt.Printf("难度值:  %d\n", block.Diffculty)
		fmt.Printf("随机数:  %d\n", block.Nonce)
		fmt.Printf("当前区块哈希:  %x\n", block.Hash)
		fmt.Printf("区块数据:  %s\n", block.Data)
		fmt.Println("=============================================")
		fmt.Println()
		if len(block.PrevHash) == 0  {
			fmt.Println("遍历结束")
			return
		}
	}
}
func (cli *CLI)AddBlock(data string)  {
	cli.bc.AddBlock(data)
	fmt.Printf("添加区块成功！\n")
}

