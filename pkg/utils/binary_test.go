package utils_test

import (
	"testing"

	"github.com/zeroFruit/vnet/pkg/utils"
)

func TestUintToByteSlice(t *testing.T) {
	data := []uint16{
		802,
		1002,
	}
	for _, num := range data {
		r := utils.Uint16ToByteSlice(num)
		result := utils.ByteSliceToUint16(r)
		if result != num {
			t.Fatalf("expected num value is %d, but got %d", num, result)
		}
	}

}
