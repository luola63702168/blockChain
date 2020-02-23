package src

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"time"
)

type Block struct {
	Version      int64
	PreBlockHash []byte // 前区块哈希值
	MerKelRoot   []byte // 梅克尔根
	TimeStamp    int64  // 时间戳
	Bits         int64  // 难度值
	Nonce        int64  //随机值
	//Data         []byte //交易信息
	Transactions []*Transaction //交易信息

	// 区块中是不存hash值的，是节点接收区块后独立计算并存储本地的
	Hash []byte // 当前区块的hash值
}

// block序列化
func (block *Block) Serialize() []byte {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(block)
	CheckErr("Serialize: ", err)
	return buffer.Bytes()
}

// 自由函数，反序列化
func Deserialize(data []byte) *Block {
	if len(data) == 0 {
		return nil
	}
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)
	CheckErr("Deserialize", err)
	return &block
}

// Create block
func NewBlock(txs []*Transaction, preBlockHash []byte, ) *Block {
	var block Block
	block = Block{
		Version:      1,
		PreBlockHash: preBlockHash,
		MerKelRoot:   []byte{},
		TimeStamp:    time.Now().Unix(),
		Bits:         targetBits,
		Nonce:        0,
		Transactions: txs,
	}
	pow := NewProofOfWork(&block)
	nonce, hash := pow.Run()
	block.Nonce = nonce
	block.Hash = hash
	return &block
}

// Create a first block
func NewGenesisBlock(coinBaseTx *Transaction) *Block {
	return NewBlock([]*Transaction{coinBaseTx}, []byte{})
}


// 生成交易信息哈希值，生成粗糙的默克尔树
func (block *Block)HashTransactions() []byte {

	var txHashes [][]byte
	txs:=block.Transactions
	for _,tx:=range txs{
		//[]bytes
		txHashes=append(txHashes,tx.TXID)
	}
	data :=bytes.Join(txHashes,[]byte{})
	hash:=sha256.Sum256(data)
	return hash[:]
}