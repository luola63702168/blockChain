package src

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"os"
)

const reward  = 12.5

type Transaction struct {
	TXID []byte // 交易id
	TXInputs []TXInput //交易输入
	TXOutputs []TXOutput //交易输入
}


type TXInput struct {
	TXID []byte // 引用uxto所在交易的ID(上一个交易id)
	Vout int64// 所消费utxo在所在的的output中的索引（来自上个交易的第几个输出）
	ScriptSig string// 解锁脚本,指明可以使用某个output的条件（所有的输出都是可见的，但是使用权只能是持有私钥签名的数据的人）

}

// 代表 TXID []byte input.Vout共同标记的之前交易的那个UTXO自己不能使用了
func (input *TXInput)CanUnlockUTXOWith(unlockData string) bool {
	return input.ScriptSig == unlockData
}


type TXOutput struct {
	Value float64// 输出金额（支付）
	ScriptPubKey string// 锁定脚本，指定收款方的地址
}

// 检查自己能不能消耗这个UTXO
func (output *TXOutput)CanBeUnlockedWith(unlockData string) bool {
	return output.ScriptPubKey == unlockData
}


// 生成交易ID
func (tx *Transaction)SetTXID()  {
	var buffer bytes.Buffer
	encoder :=gob.NewEncoder(&buffer)
	err:=encoder.Encode(tx)
	CheckErr("(tx *Transaction)SetTXID ",err)
	hash := sha256.Sum256(buffer.Bytes())
	tx.TXID=hash[:]
}

// 创建普通交易，send的辅助函数
func NewTransaction(from,to string, amount float64,bc *BlockChain) *Transaction {

	vaidUTXOs, total := bc.FindSuitableUTXO(from,amount)
	if total<amount{
		fmt.Println("Not enough money")
		os.Exit(1)
	}
	var inputs  []TXInput
	var outputs  []TXOutput

	//output-->input
	for txId,outputIndex := range vaidUTXOs{
		for _,index :=range outputIndex{
			input:=TXInput{
				TXID:      []byte(txId),
				Vout:      index,
				ScriptSig: from,
			}
			inputs = append(inputs,input)
		}
	}

	// 创建outputs
	// 普通output
	output := TXOutput{
		Value:        amount,
		ScriptPubKey: to,
	}
	// 找零output
	if total>amount{
		output := TXOutput{
			Value:        total-amount,
			ScriptPubKey: from,
		}
		outputs = append(outputs,output)
	}
	outputs = append(outputs,output)


	tx:= Transaction{
		TXID:      []byte{},
		TXInputs:  inputs,
		TXOutputs: outputs,
	}

	tx.SetTXID()
	return &tx
}

// 创建NewCoinBaseTx(只有收款人，没有付款人，是矿工的奖励交易)
func NewCoinBaseTx(address string,data string) *Transaction {
	if data == ""{
		data = fmt.Sprintf("reward to %s %f btc",address,reward)
	}

	input:=TXInput{
		TXID:      []byte{},
		Vout:      -1,
		ScriptSig: data,
	}
	output:=TXOutput{
		Value:        reward,
		ScriptPubKey: address,
	}
	tx:= Transaction{
		TXID:      []byte{},
		TXInputs:  []TXInput{input},
		TXOutputs: []TXOutput{output},
	}
	tx.SetTXID()
	return &tx
}

// 判断交易是否为第一个交易（挖矿交易）
func (tx *Transaction)IsCoinBase() bool  {
	if len(tx.TXInputs)==1{
		if len(tx.TXInputs[0].TXID)==0 && tx.TXInputs[0].Vout==-1{
			return true
		}
	}
	return false
}