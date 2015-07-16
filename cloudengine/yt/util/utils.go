package util

import (
	"encoding/binary"
)

func IntToBytes(i int) []byte {
	var buf = make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, uint32(i))
	return buf
}

func BytesToInt(buf []byte) int {
	return int(binary.LittleEndian.Uint32(buf))
}

func ShortToBytes(i int) []byte {
	var buf = make([]byte, 2)
	binary.LittleEndian.PutUint16(buf, uint16(i))
	return buf
}

func BytesToShort(buf []byte) int {
	return int(binary.LittleEndian.Uint16(buf))
}
