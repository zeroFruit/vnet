package link_test

import (
	"github.com/zeroFruit/vnet/pkg/link"
	"github.com/zeroFruit/vnet/pkg/link/na"
	"testing"
)

func TestNewFrameEncodeDecode(t *testing.T) {
	encoder := link.NewFrameEncoder()
	decoder := link.NewFrameDecoder()

	frame, err := encoder.Encode(na.Frame{
		Src: link.AddrFromStr("11-11-11-11-11-11"),
		Dest: link.AddrFromStr("22-22-22-22-22-22"),
		Payload: []byte("data"),
	})
	if err != nil {
		t.Fatalf("failed to encode frame: %v", err)
	}

	frame2, err := decoder.Decode(frame)
	if err != nil {
		t.Fatalf("failed to decode bytes to frame: %v", err)
	}
	if !frame2.Src.Equal(link.AddrFromStr("11-11-11-11-11-11")) {
		t.Fatalf("expected src address is '11-11-11-11-11-11', but got '%s'", frame2.Src)
	}
	if !frame2.Dest.Equal(link.AddrFromStr("22-22-22-22-22-22")) {
		t.Fatalf("expected src address is '22-22-22-22-22-22', but got '%s'", frame2.Src)
	}
	if string(frame2.Payload) != "data" {
		t.Fatalf("expected payload value is 'data', but got '%s'", string(frame2.Payload))
	}
}
