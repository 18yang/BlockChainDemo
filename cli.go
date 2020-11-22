package main

import (
	"fmt"
	"os"
)
//这是一个用来接受命令行参数并且控制区块链操作的文件
type CLI struct {
	bc *BlockChain
}
const Usage = `
	addBlock --data DATA       "add data to blockchain"
	printChain				   "print all blockchain data"
`
//接收参数的动作，放在一个函数中
func (cli *CLI)Run()  {
	//得到所有的命令
	args := os.Args
	// ./block printChain
	// ./block addBlock
	if len(args) < 2 {
		fmt.Printf(Usage)
		return
	}
	//2. 分析命令
	cmd := args[1]
	switch cmd {
		//3. 执行相应的动作
	case "addBlock":
		//添加区块
		fmt.Println("添加区块")
		if len(args) ==4 && args[2] == "--data"  {
			//获取命令的数据
			data := args[3]
			//添加区块
			cli.AddBlock(data)
		}else{
			fmt.Println("参数使用不当，请检查！")
			fmt.Print(Usage)
		}
	case "printChain":
		//打印区块
		fmt.Println("打印区块")
		cli.PrintBlockChain()
	default:
		fmt.Printf(Usage)
	}
}
