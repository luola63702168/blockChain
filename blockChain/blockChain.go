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
	// 检查文件信息，返回文件对象，和err
	_, err := os.Stat(dbFile)
	// 校验这个err，如果满足，则说明文件不存在
	return !os.IsNotExist(err)
	//if os.IsNotExist(err){
	//	return false
	//}
	//return true
}

// 创建一个新的区块链（也就是：Bucket）
func InitBlockChain(address string) *BlockChain {
	if IsDBExist() {
		fmt.Println("blockChain is all already, no need to create")
		os.Exit(1)
	}

	db, err := bolt.Open(dbFile, 0600, nil)
	CheckErr("InitBlockChain1", err)

	var lastHash []byte // 初始化在匿名函数外面，这样可以直接该将数据返回
	//db.View()
	err4 := db.Update(func(tx *bolt.Tx) error {
		// 没有bucket, 创建创世块并将数据填写到数据库里
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

	var lastHash []byte // 初始化在匿名函数外面，这样可以直接该将数据返回
	//db.View()
	err4 := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket)) // 取bucket
		if bucket != nil {
			// 取出最后区块的hash值
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

// 迭代器，就是一个对象，它里面包含了一个游标，一直向前(后)移动，完成整个容器的遍历
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
		// 获取区块
		byteData := bucket.Get(it.currHash)
		block = Deserialize(byteData)
		// 游标移动
		it.currHash = block.PreBlockHash
		return nil
	})
	CheckErr("Next err: ", err)
	return
}

// 在整个区块链中找寻指定地址能够支配的utxo的交易集合（未消耗的）
func (bc *BlockChain) FindUTXOTransactions(address string) []Transaction {
	// 包含目标utxo的交易集合（未消耗的）
	var UTXOTransaction []Transaction
	//  存储已经使用过 的utxo集合（已经消耗的）
	//  因为可能存在 0x111 :0,1 也就是 一个交易id的两个输出指向了同一个输入，
	//  所以这个存储结构应该是：map[交易id] []int64
	spentUTXO := make(map[string][]int64)

	it := bc.NewIterator()
	for {
		block := it.Next()
		for _, tx := range block.Transactions {
			// 如果是挖矿交易，没有输入，则不处理。
			if !tx.IsCoinBase() {
				// 遍历input
				// 处理该交易中所有的input 确定其还未使用这个utxo（如果已经使用了，那么这个交易的私钥就会由输入方的地址更改为自己的地址）
				// 所以需要两个字段来标识使用过的utxo： 交易ID 和output的索引值
				for _, input := range tx.TXInputs {
					// 私钥是address的地址，说明这个input.TXID和input.Vout共同标记的utxo被address消耗了
					if input.CanUnlockUTXOWith(address) {
						spentUTXO[string(input.TXID)] = append(spentUTXO[string(input.TXID)], input.Vout)
					}
				}
			}

		OUTPUTS:
			// 处理该交易中所有的output 确定是其所有者
			for currindex, output := range tx.TXOutputs { // currindex其实对应的就是:input.Vout
				// 检查当前的output是否已经被消耗，如果被消耗了就不添加对应的tx进去了
				if spentUTXO[string(tx.TXID)] != nil {
					// 非空，代表交易里有被adress消耗的utxo
					indexs := spentUTXO[string(tx.TXID)]
					for _, index := range indexs {
						// 此时输出索引和已消耗的输出索引相等，说明这个utxo被address消耗了，跳过，进行下一个output的判断
						if int64(currindex) == index {
							continue OUTPUTS // 跳转标签
						}
					}
				}
				// 再判断输出公钥是不是和该地址对应，确定该utxo是不是自己可以消耗的utxo。
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
	// 遍历交易
	for _, tx := range txs {
		// 遍历output
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
	// 在整个区块链中找寻指定地址能够支配的utxo的交易集合（未消耗的）
	txs := bc.FindUTXOTransactions(address)
	vaidUTXOs := make(map[string][]int64)
	var total float64
FIND:
	// 遍历交易(utxo)
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
