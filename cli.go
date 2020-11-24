package main

import (
	"fmt"
	"os"
	"strconv"
)
//这是一个用来接受命令行参数并且控制区块链操作的文件
type CLI struct {
	bc *BlockChain
}
const Usage = `
	addBlock --data DATA       "添加区块"
	printChain                 "正向打印区块链"
	printChainR				   "反向打印区块链"
	getBalance --address ADDRESS "获取指定地址的余额"
	send FROM TO AMOUNT MINER DATA "由FROM转AMOUNT给TO,由MINER挖矿,同时写入DATA"
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
	case "printChain" :
		fmt.Println("正向打印区块链\n")
		cli.PrintBlockChain()
	case "printChainR":
		//打印区块
		fmt.Println("反向打印区块链")
		cli.PrintBlockChainReverse()
	case "getBalance":
		fmt.Printf("获取余额...\n")
		if len(args) ==4 && args[2] == "--address"  {
			address := args[3]
			cli.GetBalance(address)
		}
	case "send":
		if len(args) != 7{
			fmt.Printf("参数个数错误，请检查！")
			fmt.Printf(Usage)
			return
		}
		fmt.Printf("开始转账...\n")
		from := args[2]
		to := args[3]
		amount,_ := strconv.ParseFloat( args[4] , 64 )
		miner := args[5]
		data := args[6]
		cli.Send(from,to,amount,miner,data)


	default:
		fmt.Printf(Usage)
	}
}
