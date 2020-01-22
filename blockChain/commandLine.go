package main

import (
	"flag"
	"fmt"
	"os"
)

type CLI struct {
}

const usage = `
	createChain --address ADDRESS "create a block Chain"
	printChain			  "print all blocks"
	getBalance --address ADDRESS   "get balance of the address"
	send --from FROM --to TO --amount AMOUNT   "send coin from FROM to TO"
	`

const PrintChainCmdString = "printChain"
const CreateChainCmdString = "createChain"
const GetBalanceCmdString = "getBalance"
const SendCmdString = "send"

// 提示信息打印
func (cli *CLI) printUsage() {
	fmt.Println("invalid input")
	fmt.Println(usage)
	os.Exit(1)
}

// 参数检查
func (cli *CLI) parameterCheck() {
	if len(os.Args) < 2 {
		cli.printUsage()
	}
}

// 命令行参数接收
func (cli *CLI) Run() {
	cli.parameterCheck()
	// 添加新命令（Flag ：命令标记）
	createChainCmd := flag.NewFlagSet(CreateChainCmdString, flag.ExitOnError)
	getBalanceCmd := flag.NewFlagSet(GetBalanceCmdString, flag.ExitOnError)
	printChainCmd := flag.NewFlagSet(PrintChainCmdString, flag.ExitOnError)
	sendCmd := flag.NewFlagSet(SendCmdString, flag.ExitOnError)

	createChainCmdPara := createChainCmd.String("address", "", "address info! ")
	getBalanceCmdPara := getBalanceCmd.String("address", "", "address info! ")
	fromCmdPara := sendCmd.String("from", "", "from address info! ")
	toCmdPara := sendCmd.String("to", "", "to address info! ")
	amountCmdPara := sendCmd.Float64("amount", 0, "amount info! ")

	switch os.Args[1] {

	// 命令行 产生交易
	case SendCmdString:
		err := sendCmd.Parse(os.Args[2:])
		CheckErr("cli *CLI Run()2", err)
		if sendCmd.Parsed() {
			if *fromCmdPara == "" || *toCmdPara == "" || *amountCmdPara == 0 {
				fmt.Println("err: from address should not be empty")
				cli.printUsage()
			}
			cli.Send(*fromCmdPara, *toCmdPara, *amountCmdPara)
		}

	// 创建区块链（初始化）
	case CreateChainCmdString:
		err := createChainCmd.Parse(os.Args[2:])
		CheckErr("cli *CLI Run()2", err)
		if createChainCmd.Parsed() {
			if *createChainCmdPara == "" {
				fmt.Println("err: send cmd parameters invalid")
				cli.printUsage()
			}
			cli.CreateChain(*createChainCmdPara)
		}

	// 命令行遍历区块，打印数据
	case PrintChainCmdString:
		err := printChainCmd.Parse(os.Args[2:])
		CheckErr("cli *CLI Run()3", err)
		if printChainCmd.Parsed() {
			cli.PrintChain()
		}

	// 查询用户余额
	case GetBalanceCmdString:
		err := getBalanceCmd.Parse(os.Args[2:])
		CheckErr("cli *CLI Run()4", err)
		if getBalanceCmd.Parsed() {
			if *getBalanceCmdPara == "" {
				fmt.Println("err: getBalance data should not be empty")
				cli.printUsage()
			}
			cli.GetBalance(*getBalanceCmdPara)
		}
	default:
		cli.printUsage()
	}
}
