package link

import (
	"github.com/zeroFruit/vnet/pkg/types"
)

type Addr string

func (a Addr) Equal(o types.HwAddr) bool {
	return a == o
}

func (a Addr) String() string {
	return string(a)
}

func AddrFromStr(s string) types.HwAddr {
	return Addr(s)
}

const BroadcastAddr Addr = "FF:FF:FF:FF:FF:FF"
