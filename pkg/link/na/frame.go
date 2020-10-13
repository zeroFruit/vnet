package na

import (
	"github.com/zeroFruit/vnet/pkg/types"
	"time"
)

type Frame struct {
	Src     types.HwAddr
	Dest    types.HwAddr
	Payload []byte
}

type FrameData struct {
	// Buf is serialized Frame structure data
	Buf []byte

	// Incoming represents interface address which receives this FrameData
	Incoming  types.HwAddr
	Timestamp time.Time
}
