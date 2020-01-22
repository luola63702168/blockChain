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
	// 目标值，由于是hash（256位），数字很大，所以使用go自带的big.Int类型进行存储
	target *big.Int
}

// 这是难度值
//const targetBits int = 24
const targetBits = 24

// 计算工作量Pow的函数
func NewProofOfWork(block *Block) *ProofOfWork {
	target := big.NewInt(1)                  // target 是一个int指针，值为1
	// 目标值前面有五个0 因为末尾的1本身占一位
	target.Lsh(target, uint(256-targetBits)) // 二进制左移,(为什么是uint，因为这是无符号int，二进制的第一位不表示符号，仍要参加运算)
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
	// 拼装数据
	var hash[32]byte
	var nonce int64 = 0
	var hashInt big.Int

	fmt.Println("Begin Mining...")
	fmt.Printf("target hash : %x\n",pow.target.Bytes())
	for nonce<math.MaxInt64{
		data:=pow.PrepareDat(nonce)
		hash = sha256.Sum256(data)
		//哈希值转成big.Int类型
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