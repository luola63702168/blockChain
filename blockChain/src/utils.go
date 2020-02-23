package src

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
)

// int to byte
func IntToByte(num int64) []byte {
	var buffer bytes.Buffer
	err := binary.Write(&buffer, binary.BigEndian, num)
	CheckErr("IntToByte", err)
	return buffer.Bytes()
}

// check err
func CheckErr(pos string, err error) {
	if err != nil {
		fmt.Println("error, pos: ", pos, err)
		os.Exit(1)

	}
}

