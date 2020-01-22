package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
)

// int to byte
func IntToByte(num int64) []byte {
	var buffer bytes.Buffer //Buffer是一个实现了读写方法的可变大小的字节缓冲。本类型的零值是一个空的可用于读写的缓冲。
	//binary.Write()参数:将data的binary编码格式写入w，data必须是定长值、定长值的切片、定长值的指针。order指定写入数据的字节序（大端还是小端模式），写入结构体时，名字中有'_'的字段会置为0。
	err:=binary.Write(&buffer,binary.BigEndian,num)
	CheckErr("IntToByte",err)
	return buffer.Bytes()
}

// check err
func CheckErr(pos string, err error)  {
	if err!=nil{
		// Exit让当前程序以给出的状态码code退出。一般来说，状态码0表示成功，非0表示出错。程序会立刻终止，defer的函数不会被执行。
		fmt.Println("error, pos: ",pos,err)
		os.Exit(1)

	}
}

