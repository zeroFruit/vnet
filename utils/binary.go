package utils

import (
	"bytes"
	"encoding/binary"
)

func Uint16ToByteSlice(num uint16) []byte {
	b := new(bytes.Buffer)
	if err := binary.Write(b, binary.LittleEndian, num); err != nil {
		panic(err)
	}
	return b.Bytes()
}

func ByteSliceToUint16(b []byte) uint16 {
	var num uint16
	if err := binary.Read(bytes.NewBuffer(b), binary.LittleEndian, &num); err != nil {
		panic(err)
	}
	return num
}
