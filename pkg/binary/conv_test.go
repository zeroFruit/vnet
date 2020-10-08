package binary_test

import (
	"testing"

	"github.com/zeroFruit/vnet/pkg/binary"
)

func TestUintToByteSlice(t *testing.T) {
	data := []uint16{
		802,
		1002,
	}
	for _, num := range data {
		r := binary.FromUint16(num)
		result := binary.ToUint16(r)
		if result != num {
			t.Fatalf("expected num value is %d, but got %d", num, result)
		}
	}

}
