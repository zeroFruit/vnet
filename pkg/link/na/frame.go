package na

import (
	"time"

	"github.com/zeroFruit/vnet/pkg/types"
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
	Incoming types.HwAddr

	// Timestamp represents times when these frame data created
	Timestamp time.Time
}
