package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
)

// 工具量证明的结构体ProofOWork
type ProofOfWork struct {
	block *Block
	// 目标值
	target *big.Int
}

// 这是难度值
const targetBits = 24

// 计算工作量Pow的函数
func NewProofOfWork(block *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))
	pow := ProofOfWork{block: block, target: target}
	return &pow
}

// 工作量证明准备方法
func (pow *ProofOfWork)PrepareDat(nonce int64) []byte {
	block:=pow.block
	block.MerKelRoot=block.HashTransactions()
	tmp:=[][]byte{
		IntToByte(block.Version),
		block.PreBlockHash,
		block.MerKelRoot,
		IntToByte(block.TimeStamp),
		IntToByte(targetBits),
		IntToByte(nonce),
	}
	data:=bytes.Join(tmp,[]byte{})
	return data
}
// 工作量证明方法：计算
func (pow *ProofOfWork)Run() (int64, []byte) {
	var hash[32]byte
	var nonce int64 = 0
	var hashInt big.Int

	fmt.Println("Begin Mining...")
	fmt.Printf("target hash : %x\n",pow.target.Bytes())
	for nonce<math.MaxInt64{
		data:=pow.PrepareDat(nonce)
		hash = sha256.Sum256(data)
		hashInt.SetBytes(hash[:])
		if hashInt.Cmp(pow.target)==-1{
			fmt.Printf("found nonce, hash:%x,nonce:%d\n",hash,nonce)
			break
		}else {
			nonce++
		}
	}
	return nonce,hash[:]
}

// 校验工作量是否可信方法
func (pow *ProofOfWork)IsValid() bool {
	var hashInt big.Int
	data:=pow.PrepareDat(pow.block.Nonce)
	hash:=sha256.Sum256(data)
	hashInt.SetBytes(hash[:])
	return hashInt.Cmp(pow.target) == -1
}