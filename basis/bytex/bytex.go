package bytex

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"strings"
)

func BytesToInt(in []byte) (int32, error) {
	// i32 := binary.LittleEndian.Uint32(in)

	var i int32
	reader := bytes.NewReader(in)

	if err := binary.Read(reader, binary.LittleEndian, &i); err != nil {
		return -1, err
	}

	return i, nil
}

func IntToBytes(i int32) ([]byte, error) {
	// b := make([]byte, 4)

	buf := bytes.NewBuffer([]byte{})

	if err := binary.Write(buf, binary.LittleEndian, i); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func BytesToString(in []byte) string {

	data := hex.EncodeToString(in)
	data = strings.ToUpper(data)

	return data
}

/*
// 通信地址转[]byte
func AddrToByte(addr int) ([]byte, byte) {
	len1 := addr % 256
	len2 := addr / 256 % 256
	len3 := addr / 256 / 256 % 256
	len4 := addr / 256 / 256 / 256 % 256

	len1Byte := byte(len1)
	len2Byte := byte(len2)
	len3Byte := byte(len3)
	len4Byte := byte(len4)

	return []byte{len1Byte, len2Byte, len3Byte, len4Byte}, len1Byte ^ len2Byte ^ len3Byte ^ len4Byte
}

// 通信地址转int
func AddrToInt(bytes []byte) int {
	if len(bytes) == 4 {
		len1 := int(bytes[0])
		len2 := int(bytes[1])
		len3 := int(bytes[2])
		len4 := int(bytes[3])

		return len1%256 + len2*256 + len3*256*256 + len4*256*256*256
	}

	log.Error("Serial Addr Error: %X", bytes)
	return 0
}
*/
