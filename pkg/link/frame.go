package link

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/zeroFruit/vnet/pkg/link/na"
)

type FrameEncoder struct{}

func NewFrameEncoder() *FrameEncoder {
	gob.Register(Addr(""))
	return &FrameEncoder{}
}

func (e *FrameEncoder) Encode(frame na.Frame) ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0))
	if err := gob.NewEncoder(buf).Encode(frame); err != nil {
		return nil, fmt.Errorf("failed to encode frame: %v", err)
	}
	b := buf.Bytes()
	return b, nil
}

type FrameDecoder struct{}

func NewFrameDecoder() *FrameDecoder {
	return &FrameDecoder{}
}

func (e *FrameDecoder) Decode(b []byte) (na.Frame, error) {
	var frame na.Frame
	decoder := gob.NewDecoder(bytes.NewBuffer(b))
	if err := decoder.Decode(&frame); err != nil {
		return na.Frame{}, fmt.Errorf("failed to decode frame: %v", err)
	}
	return frame, nil
}
