package bytex

import (
	"testing"
)

var testData = []struct {
	In  []byte
	Out int32
}{
	{[]byte{0x60, 0xa5, 0xff, 0xff}, -23200},
	{[]byte{0xe8, 0x03, 0x00, 0x00}, 1000},
}

func TestBytesToInt(t *testing.T) {
	for _, v := range testData {
		// bytes to int
		i, err := BytesToInt(v.In)
		if err != nil {
			panic(err)
		}

		if v.Out != i {
			t.Fatal("out err: ", v.Out, i)
		}

		t.Log(i)
	}
}

func TestIntToBytes(t *testing.T) {
	for _, v := range testData {
		// bytes to int
		bytes, err := IntToBytes(v.Out)
		if err != nil {
			panic(err)
		}

		if v.In[0] != bytes[0] {
			t.Fatal("out err: ", v.In, bytes)
		}

		t.Logf("%X", bytes)
	}
}
