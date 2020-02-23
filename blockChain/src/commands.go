package src

import (
	"fmt"
)


// 命令行遍历区块，打印数据
func (cli *CLI) PrintChain() {
	bc := GetBlockChainHandler()
	it := bc.NewIterator()
	defer bc.db.Close()
	for {
		block := it.Next()
		fmt.Printf("Version%d\n", block.Version)
		fmt.Printf("PreBlockHash:%x\n", block.PreBlockHash)
		fmt.Printf("Hash:%x\n", block.Hash)
		//fmt.Println(block.Hash) //测试无%x 为什么样
		fmt.Printf("TimeStamp:%d\n", block.TimeStamp)
		fmt.Printf("Bits:%d\n", block.Bits)
		fmt.Printf("Nonce:%d\n", block.Nonce)
		fmt.Printf("IsValid:%v\n", NewProofOfWork(block).IsValid())
		if len(block.PreBlockHash) == 0 {
			fmt.Println("print over")
			break
		}
	}
}

// 命令行，创建区块链（初始化）
func (cli *CLI) CreateChain(address string) {
	bc := InitBlockChain(address)
	defer bc.db.Close()
	fmt.Println("Create block successfully")
}

// 命令行，获取余额
func (cli *CLI) GetBalance(address string)  {
	bc := GetBlockChainHandler()
	defer bc.db.Close()
	utxos := bc.FindUTXO(address)

	var total float64 = 0
	for _, utxo := range utxos {
		total += utxo.Value
	}
	fmt.Printf("The balance of %s is : %f\n",address,total)
}

// 命令行： 发生交易
func (cli *CLI)Send(from,to string, amount float64)  {
	bc := GetBlockChainHandler()
	defer bc.db.Close()

	tx := NewTransaction(from,to,amount,bc)
	bc.AddBlock([]*Transaction{tx})
	fmt.Println("send successfully!")

}