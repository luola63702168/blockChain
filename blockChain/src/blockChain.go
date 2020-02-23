package main

import (
	"fmt"
	"github.com/boltdb/bolt"
	"os"
)

const dbFile = "./obj_file/blockChain.db" //db
const blockBucket = "bucket"              //bucket
const lastHashKey = "lastHashKey"         // block lastHashKey
const genesisInfo = "this genesisInfo"    // 创世块 ScriptSig

type BlockChain struct {
	db *bolt.DB
	// 尾巴，表示最后一个区块的hash值
	tail []byte
}

// 检查有没有数据库文件（如果有这个文件就一定有那个：Bucket）
func IsDBExist() bool {
	_, err := os.Stat(dbFile)
	return !os.IsNotExist(err)
}

// 创建一个新的区块链（也就是：Bucket）
func InitBlockChain(address string) *BlockChain {
	if IsDBExist() {
		fmt.Println("blockChain is all already, no need to create")
		os.Exit(1)
	}

	db, err := bolt.Open(dbFile, 0600, nil)
	CheckErr("InitBlockChain1", err)

	var lastHash []byte
	//db.View()
	err4 := db.Update(func(tx *bolt.Tx) error {
		coinBaseTx := NewCoinBaseTx(address, genesisInfo)
		genesisBlock := NewGenesisBlock(coinBaseTx)
		bucket, err := tx.CreateBucket([]byte(blockBucket))
		CheckErr("InitBlockChain2", err)
		err1 := bucket.Put(genesisBlock.Hash, genesisBlock.Serialize())
		CheckErr("InitBlockChain3", err1)
		err2 := bucket.Put([]byte(lastHashKey), genesisBlock.Hash)
		CheckErr("InitBlockChain4", err2)
		lastHash = genesisBlock.Hash
		return nil
	})

	CheckErr("NewBlockChain3", err4)
	return &BlockChain{
		db:   db,
		tail: lastHash,
	}

}

// 获取区块链（bucket）
func GetBlockChainHandler() *BlockChain {
	if !IsDBExist() {
		fmt.Println("please create blockChain first")
		os.Exit(1)
	}
	db, err := bolt.Open(dbFile, 0600, nil)
	CheckErr("GetBlockChainHandler1", err)

	var lastHash []byte
	//db.View()
	err4 := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket)) // 取bucket
		if bucket != nil {
			lastHash = bucket.Get([]byte(lastHashKey))
		}
		return nil
	})

	CheckErr("GetBlockChainHandler2", err4)
	return &BlockChain{
		db:   db,
		tail: lastHash,
	}
}

// 添加区块
func (bc *BlockChain) AddBlock(txs []*Transaction) {
	var prevBlockHash []byte
	err := bc.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			os.Exit(1)
		}
		prevBlockHash = bucket.Get([]byte(lastHashKey))
		return nil
	})
	CheckErr("bc.db.View", err)
	block := NewBlock(txs, prevBlockHash)
	err1 := bc.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			os.Exit(1)
		}
		err1 := bucket.Put(block.Hash, block.Serialize())
		CheckErr("AddBlock1", err1)
		err2 := bucket.Put([]byte(lastHashKey), block.Hash)
		CheckErr("AddBlock2", err2)
		bc.tail = block.Hash
		return nil
	})
	CheckErr("bc.db.Update", err1)
}

// 迭代器
type BlockChainIterator struct {
	currHash []byte // 游标，总指向当前的hash
	db       *bolt.DB
}

// 创建一个迭代器对象，同时初始化指向最后一个区块
func (bc *BlockChain) NewIterator() *BlockChainIterator {
	return &BlockChainIterator{currHash: bc.tail, db: bc.db}
}

// 迭代器本身
func (it *BlockChainIterator) Next() (block *Block) {
	err := it.db.View(func(tx *bolt.Tx) error {
		// 获取bucket
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			return nil
		}
		byteData := bucket.Get(it.currHash)
		block = Deserialize(byteData)
		it.currHash = block.PreBlockHash
		return nil
	})
	CheckErr("Next err: ", err)
	return
}

// 在整个区块链中找寻指定地址能够支配的utxo的交易集合（未消耗的）
func (bc *BlockChain) FindUTXOTransactions(address string) []Transaction {
	var UTXOTransaction []Transaction
	spentUTXO := make(map[string][]int64)

	it := bc.NewIterator()
	for {
		block := it.Next()
		for _, tx := range block.Transactions {
			if !tx.IsCoinBase() {
				for _, input := range tx.TXInputs {
					if input.CanUnlockUTXOWith(address) {
						spentUTXO[string(input.TXID)] = append(spentUTXO[string(input.TXID)], input.Vout)
					}
				}
			}

		OUTPUTS:
			for currindex, output := range tx.TXOutputs {
				if spentUTXO[string(tx.TXID)] != nil {
					indexs := spentUTXO[string(tx.TXID)]
					for _, index := range indexs {
						if int64(currindex) == index {
							continue OUTPUTS
						}
					}
				}
				if output.CanBeUnlockedWith(address) {
					UTXOTransaction = append(UTXOTransaction, *tx)
				}
			}

		}
		if len(block.PreBlockHash) == 0 {
			break
		}
	}
	return UTXOTransaction
}

// 寻找当前地址能够使用的utxo
func (bc *BlockChain) FindUTXO(address string) []TXOutput {

	var UTXOs []TXOutput
	txs := bc.FindUTXOTransactions(address)
	for _, tx := range txs {
		for _, utxo := range tx.TXOutputs {
			if utxo.CanBeUnlockedWith(address) {
				UTXOs = append(UTXOs, utxo)
			}
		}
	}
	return UTXOs

}

// 找到合理的utxo

func (bc *BlockChain) FindSuitableUTXO(address string, amount float64) (map[string][]int64, float64) {
	txs := bc.FindUTXOTransactions(address)
	vaidUTXOs := make(map[string][]int64)
	var total float64
FIND:
	for _, tx := range txs {
		outputs := tx.TXOutputs
		for index, output := range outputs {
			if output.CanBeUnlockedWith(address) {
				if total < amount {
					total += output.Value
					vaidUTXOs[string(tx.TXID)] = append(vaidUTXOs[string(tx.TXID)], int64(index))
				} else {
					break FIND
				}

			}
		}
	}
	return vaidUTXOs, total
}
