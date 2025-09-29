package counter

import (
	"encoding/hex"
	"strings"
	"testing"
)

func TestLoadFile(t *testing.T) {
	data := loadContractHexBytes()
	bin := CounterMetaData.Bin
	hexStr, found := strings.CutPrefix(bin, "0x")
	if !found {
		t.Fatal("Contract hex prefix not found")
	}
	decode, err := hex.DecodeString(hexStr)
	if err != nil {
		t.Fatal(err)
	}
	if len(decode) != len(data) {
		t.Fatal("Decoded length does not match")
	}
	for i := 0; i < len(decode); i++ {
		if decode[i] != data[i] {
			t.Fatal("Decoded value does not match")
		}
	}
}

func TestCounter(t *testing.T) {
	callCounter()
}
